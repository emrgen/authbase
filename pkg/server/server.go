package server

import (
	"context"
	"errors"
	"fmt"
	gatewayfile "github.com/black-06/grpc-gateway-file"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/config"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/service"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/emrgen/authbase/x/mail"
	gopackv1 "github.com/emrgen/gopack/apis/v1"
	"github.com/gobuffalo/packr"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcvalidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

type Server struct {
	config          *config.Config
	provider        store.Provider
	redis           *cache.Redis
	permission      permission.AuthBasePermission
	mailer          mail.MailerProvider
	adminOrgService v1.AdminProjectServiceServer
	gl              net.Listener
	rl              net.Listener
	grpcServer      *grpc.Server
	mux             *runtime.ServeMux
	httpPort        string
	grpcPort        string
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
	logrus.Infof("APP_MODE: %v", s.config.Mode)

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
	db := store.GetDB()
	// if multistore mode use multistore provider
	s.provider = store.NewDefaultProvider(db)
	s.redis = cache.NewRedisClient()
	//s.permission = permission.NewStoreBasedPermission(s.provider)
	s.permission = permission.NewStoreBasedPermission(s.provider)
	s.mailer = mail.NewMailerProvider("smtp.gmail.com", 587, "", "")

	// migrate the database
	err := db.Migrate()
	if err != nil {
		return err
	}

	s.grpcPort = ":" + grpcPort
	s.httpPort = ":" + httpPort

	gl, err := net.Listen("tcp", s.grpcPort)
	if err != nil {
		return err
	}
	s.gl = gl

	rl, err := net.Listen("tcp", s.httpPort)
	if err != nil {
		return err
	}
	s.rl = rl

	return nil
}

func (s *Server) registerServices() error {
	var err error
	verifier := x.NewStoreBasedUserVerifier(s.provider, s.redis)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcvalidator.UnaryServerInterceptor(),
			x.AuthInterceptor(verifier),
			UnaryGrpcRequestTimeInterceptor(),
		)),
	)

	cookieStore := NewCookieStore(s.redis)

	// connect the rest gateway to the grpc server
	s.mux = runtime.NewServeMux(
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
		runtime.WithForwardResponseOption(InjectCookie(cookieStore)),
		runtime.WithMetadata(ExtractCookie(cookieStore)),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(UnaryRequestTimeInterceptor()),
	}
	endpoint := "localhost" + s.grpcPort

	redis := s.redis
	perm := s.permission

	// Register the grpc services

	//oauthService := service.NewOAuth2Service(s.provider, redis)
	//offlineTokenService := service.NewOfflineTokenService(perm, s.provider, redis)

	v1.RegisterAdminProjectServiceServer(grpcServer, service.NewAdminProjectService(s.provider, redis))
	v1.RegisterProjectServiceServer(grpcServer, service.NewProjectService(perm, s.provider, redis))
	v1.RegisterClientServiceServer(grpcServer, service.NewClientService(perm, s.provider, redis))
	v1.RegisterAuthServiceServer(grpcServer, service.NewAuthService(s.provider, perm, s.mailer, redis))
	//v1.RegisterOAuth2ServiceServer(grpcServer, oauthService)
	//v1.RegisterAccountServiceServer()
	//gopackv1.RegisterTokenServiceServer(grpcServer, service.NewTokenService(offlineTokenService, oauthService))

	// Register the rest gateway
	if err = v1.RegisterProjectServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterProjectServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterAuthServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = gopackv1.RegisterTokenServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}

	s.grpcServer = grpcServer

	return err
}

// run the server and listen on grpc and http ports
func (s *Server) run() error {
	apiMux := http.NewServeMux()
	docsPath := "/v1/docs/"
	openapiDocs := packr.NewBox("../../docs/v1")
	apiMux.Handle(docsPath, http.StripPrefix(docsPath, http.FileServer(openapiDocs)))
	apiMux.Handle("/", s.mux)

	// CORS middleware
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

	logrus.Infof("Press Ctrl+C to stop the server")

	// if an admin project is provided, create the org and the super admin user
	if s.config.AdminOrg.Valid() {
		logrus.Infof("trying to create admin project: %v", s.config.AdminOrg)
		adminOrgService := service.NewAdminProjectService(s.provider, s.redis)
		_, err := adminOrgService.CreateAdminProject(context.TODO(), &v1.CreateAdminProjectRequest{
			Name:        s.config.AdminOrg.OrgName,
			VisibleName: s.config.AdminOrg.VisibleName,
			Email:       s.config.AdminOrg.Email,
			Password:    &s.config.AdminOrg.Password,
		})
		if err != nil {
			if !errors.Is(err, store.ErrProjectExists) {
				logrus.Infof("admin project already exists, skipping creation")
			} else {
				logrus.Errorf("error creating admin project: %v", err)
				return err
			}
		}
	}

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
