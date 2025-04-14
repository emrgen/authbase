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
	"github.com/emrgen/authbase/pkg/secret"
	"github.com/emrgen/authbase/pkg/service"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/emrgen/authbase/x/mail"
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
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Server is the main server struct
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
	ready           chan struct{}
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
	return &Server{config: config, ready: make(chan struct{})}
}

// Start the server
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

func (s *Server) Ready() <-chan struct{} {
	return s.ready
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

// register the services with the grpc server
func (s *Server) registerServices() error {
	var err error
	keyProvider := x.NewStaticKeyProvider(x.JWTSecretFromEnv())
	verifier := x.NewStoreBasedTokenVerifier(s.provider, s.redis)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcvalidator.UnaryServerInterceptor(),
			x.AuthInterceptor(verifier, keyProvider, s.provider),
			UnaryGrpcRequestTimeInterceptor(),
		)),
	)
	s.grpcServer = grpcServer

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

	secrets := secret.NewMemStore()

	// Register the grpc services
	v1.RegisterAdminProjectServiceServer(grpcServer, service.NewAdminProjectService(s.provider, redis))
	v1.RegisterProjectServiceServer(grpcServer, service.NewProjectService(perm, s.provider, redis))
	v1.RegisterClientServiceServer(grpcServer, service.NewClientService(perm, s.provider, secrets))
	v1.RegisterAuthServiceServer(grpcServer, service.NewAuthService(s.provider, keyProvider, perm, s.mailer, redis))
	v1.RegisterAccountServiceServer(grpcServer, service.NewAccountService(perm, s.provider, redis))
	v1.RegisterAccessKeyServiceServer(grpcServer, service.NewAccessKeyService(perm, s.provider, redis, keyProvider, verifier))
	v1.RegisterPoolServiceServer(grpcServer, service.NewPoolService(s.provider, perm))
	v1.RegisterPoolMemberServiceServer(grpcServer, service.NewPoolMemberService(s.provider))
	v1.RegisterTokenServiceServer(grpcServer, service.NewTokenService(verifier))
	v1.RegisterGroupServiceServer(grpcServer, service.NewGroupService(s.provider))
	v1.RegisterRoleServiceServer(grpcServer, service.NewRoleService(s.provider))
	v1.RegisterApplicationServiceServer(grpcServer, service.NewApplicationService(s.provider))
	v1.RegisterProjectMemberServiceServer(grpcServer, service.NewProjectMemberService(perm, s.provider, redis))
	v1.RegisterAdminAuthServiceServer(grpcServer, service.NewAdminAuthService(s.provider, s.config.AdminOrg, keyProvider, redis))

	// Register the http gateway
	if err = v1.RegisterAdminProjectServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterProjectServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterClientServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterAuthServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterAccountServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterAccessKeyServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterPoolServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterPoolMemberServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterTokenServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}
	if err = v1.RegisterGroupServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}

	if err = v1.RegisterRoleServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}

	if err = v1.RegisterApplicationServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}

	if err = v1.RegisterProjectMemberServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}

	if err = v1.RegisterAdminAuthServiceHandlerFromEndpoint(context.TODO(), s.mux, endpoint, opts); err != nil {
		return err
	}

	return err
}

// run the server and listen on grpc and http ports
func (s *Server) run() error {
	apiMux := http.NewServeMux()
	docsPath := "/v1/docs/"
	openapiDocs := packr.NewBox("../../docs/v1")
	apiMux.Handle(docsPath, http.StripPrefix(docsPath, http.FileServer(openapiDocs)))
	apiMux.Handle("/", s.mux)

	logger := logrus.New()
	logger.SetReportCaller(true)

	grpclog.SetLoggerV2(grpclog.NewLoggerV2(io.Discard, io.Discard, io.Discard))

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // All origins are allowed
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Logger:           logger,
	})

	restServer := &http.Server{
		Addr:    s.httpPort,
		Handler: c.Handler(apiMux),
	}

	// make sure to wait for the servers to stop before exiting
	var wg sync.WaitGroup

	logrus.Info("API documentation: http://localhost", s.httpPort, "/v1/docs/")

	// Start the grpc server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.Info("starting rest gateway on: ", s.httpPort)
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

	logrus.Infof("-----------------------------------------------")
	logrus.Info("BOOTSTRAP: waiting for the server to be ready")
	logrus.Infof("-----------------------------------------------")

	// if an admin project is provided, create the org and the super admin user
	// TODO: there are some issues with the admin project creation, need to fix it.
	if s.config.AdminOrg.Valid() {
		//TODO: remove this check as this logs the client secret
		adminOrgService := service.NewAdminProjectService(s.provider, s.redis)
		_, err := adminOrgService.CreateAdminProject(context.TODO(), &v1.CreateAdminProjectRequest{
			Name:         s.config.AdminOrg.OrgName,
			VisibleName:  s.config.AdminOrg.VisibleName,
			Email:        s.config.AdminOrg.Email,
			Password:     &s.config.AdminOrg.Password,
			ClientId:     s.config.AdminOrg.ClientId,
			ClientSecret: s.config.AdminOrg.ClientSecret,
		})
		if err != nil {
			if !errors.Is(err, store.ErrProjectExists) {
				logrus.Infof("admin project already exists, admin project creation skipped")
			} else {
				logrus.Errorf("error creating admin project: %v", err)
				return err
			}
		}
	} else {
		logrus.Errorf("no admin project provided, skipping admin project creation")
		os.Exit(0)
	}

	go func() {
		// wait for 1sec
		time.Sleep(100 * time.Millisecond)
		s.ready <- struct{}{}
	}()

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
