package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/config"
	"github.com/emrgen/authbase/pkg/service"
	"github.com/emrgen/authbase/pkg/store"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	gatewayfile "github.com/black-06/grpc-gateway-file"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/gobuffalo/packr"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcvalidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func UnaryRequestTimeInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		reqTime := time.Since(start)
		logrus.Infof("request time: %v: %v", method, reqTime)
		return err
	}
}

type Server struct {
	config     *config.Config
	db         store.AuthBaseStore
	redis      *cache.Redis
	gl         net.Listener
	rl         net.Listener
	grpcServer *grpc.Server
	mux        *runtime.ServeMux
	httpPort   string
	grpcPort   string
}

// NewServerFromEnv creates a new server instance from the environment configuration.
func NewServerFromEnv() *Server {
	cfg, err := config.FromEnv()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	return NewServer(cfg)
}

// NewServer creates a new server instance.
func NewServer(config *config.Config) *Server {
	return &Server{config: config}
}

func (s *Server) Start(grpcPort, httpPort string) error {
	if err := s.init(grpcPort, httpPort); err != nil {
		return err
	}

	if err := s.registerServices(); err != nil {
		return err
	}

	if err := s.run(); err != nil {
		return err
	}

	return nil
}

func (s *Server) init(grpcPort, httpPort string) error {
	s.db = config.GetDB()
	s.redis = cache.NewRedisClient()

	grpcPort = ":" + grpcPort
	httpPort = ":" + httpPort

	gl, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return err
	}
	s.gl = gl

	rl, err := net.Listen("tcp", httpPort)
	if err != nil {
		return err
	}
	s.rl = rl

	return nil
}

func (s *Server) registerServices() error {
	var err error
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcvalidator.UnaryServerInterceptor(),
		)),
	)

	// connect the rest gateway to the grpc server
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
		gatewayfile.WithHTTPBodyMarshaler(),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(UnaryRequestTimeInterceptor()),
	}
	endpoint := "localhost" + s.grpcPort

	rdb := s.db
	redis := s.redis

	// Register the grpc server
	v1.RegisterOrganizationServiceServer(grpcServer, service.NewOrganizationService(rdb, redis))
	v1.RegisterMemberServiceServer(grpcServer, service.NewMemberService(rdb, redis))
	v1.RegisterUserServiceServer(grpcServer, service.NewUserService(rdb, redis))
	v1.RegisterPermissionServiceServer(grpcServer, service.NewPermissionService(rdb, redis))
	v1.RegisterAuthServiceServer(grpcServer, service.NewAuthService(rdb, redis))
	v1.RegisterOauthServiceServer(grpcServer, service.NewOauthService(rdb, redis))
	v1.RegisterTokenServiceServer(grpcServer, service.NewTokenService(rdb, redis))

	// Register the rest gateway
	if err = v1.RegisterOrganizationServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterMemberServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterUserServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterPermissionServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterAuthServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterOauthServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterTokenServiceHandlerFromEndpoint(context.TODO(), mux, endpoint, opts); err != nil {
		return err
	}

	s.grpcServer = grpcServer

	return err
}

// run the server and listen on grpc and http ports
func (s *Server) run() error {

	apiMux := http.NewServeMux()
	openapiDocs := packr.NewBox("../../docs/v1")
	docsPath := "/v1/docs/"
	apiMux.Handle(docsPath, http.StripPrefix(docsPath, http.FileServer(openapiDocs)))
	apiMux.Handle("/", s.mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // All origins are allowed
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
	})

	restServer := &http.Server{
		Addr:    s.httpPort,
		Handler: c.Handler(apiMux),
	}

	// make sure to wait for the servers to stop before exiting
	var wg sync.WaitGroup

	// Start the grpc server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.Info("starting rest gateway on: ", s.httpPort)
		logrus.Info("click on the following link to view the API documentation: http://localhost", s.httpPort, "/v1/docs/")
		if err := restServer.Serve(s.rl); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logrus.Errorf("error starting rest gateway: %v", err)
			}
		}
		logrus.Infof("rest gateway stopped")
	}()

	// Start the rest gateway
	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.Info("starting grpc server on: ", s.grpcPort)
		if err := s.grpcServer.Serve(s.gl); err != nil {
			logrus.Infof("grpc failed to start: %v", err)
		}
		logrus.Infof("grpc server stopped")
	}()

	time.Sleep(1 * time.Second)
	logrus.Infof("Press Ctrl+C to stop the server")

	// listen for interrupt signal to gracefully shut down the server
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	<-sigs
	// clean Ctrl+C output
	fmt.Println()

	s.grpcServer.Stop()
	err := restServer.Shutdown(context.Background())
	if err != nil {
		logrus.Errorf("error stopping rest gateway: %v", err)
		return err
	}

	wg.Wait()

	return nil
}
