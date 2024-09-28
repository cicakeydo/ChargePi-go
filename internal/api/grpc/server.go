package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/ChargePi/ChargePi-go/internal/auth"
	chargePoint "github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	settings "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ChargePi-go/internal/users"
	grpc2 "github.com/ChargePi/ChargePi-go/pkg/proto/v1/grpc"
	"github.com/ChargePi/ChargePi-go/pkg/tls"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// GRPC API configuration
type Configuration struct {
	// Enabled is a flag to enable or disable the API
	Enabled bool `json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`

	// Address is the address where the API will be served. It should be in the format of host:port.
	Address string `json:"address,omitempty" yaml:"address" mapstructure:"address"`

	// TLS is the configuration for the TLS
	TLS tls.TLS `json:"tls,omitempty" yaml:"tls" mapstructure:"tls"`
}

type Server struct {
	server             *grpc.Server
	logger             *log.Logger
	address            string
	service            *Service
	authService        *TagAuthService
	chargePointService *ChargePointService
	logService         *LogService
	userService        *UserService
}

func NewServer(
	settings Configuration,
	point chargePoint.ChargePoint,
	authCache auth.Service,
	manager manager.Manager,
	settingsManager settings.Manager,
	userService users.Service,
) (*Server, error) {
	var opts []grpc.ServerOption

	logger := log.StandardLogger()

	if settings.TLS.IsEnabled {
		// Add TLS if enabled
		tlsCredentials, err := credentials.NewServerTLSFromFile(settings.TLS.CACertificatePath, settings.TLS.PrivateKeyPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to TLS certificates")
		}

		opts = append(opts, grpc.Creds(tlsCredentials))
	}

	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	// Add authentication, recovery and logging middleware
	opts = append(opts, grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
		logging.UnaryServerInterceptor(interceptorLogger(logger), logOpts...),
		grpcauth.UnaryServerInterceptor(authMiddleware(userService)),
		grpcrecovery.UnaryServerInterceptor(),
	)))

	return &Server{
		logger:             logger,
		server:             grpc.NewServer(opts...),
		address:            settings.Address,
		service:            NewEvseService(manager),
		authService:        NewTagAuthService(authCache),
		chargePointService: NewChargePointService(point, settingsManager),
		logService:         NewLogService(nil),
		userService:        NewUserService(userService),
	}, nil
}

// Run starts the gRPC server
func (s *Server) Run() error {
	// Business services
	grpc2.RegisterChargePointServer(s.server, s.chargePointService)
	grpc2.RegisterEvseServer(s.server, s.service)
	grpc2.RegisterLogServer(s.server, s.logService)
	grpc2.RegisterTagServer(s.server, s.authService)
	grpc2.RegisterUsersServer(s.server, s.userService)

	// Health check
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(s.server, healthcheck)

	s.logger.Infof("Exposing API endpoints at %s", s.address)

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("unable to listen to provided address: %s", s.address)
	}

	err = s.server.Serve(listener)
	if err != nil {
		return fmt.Errorf("unable to listen to provided address: %s", s.address)
	}

	return nil
}

// Stop stops the gRPC server
func (s *Server) Stop() {
	s.server.GracefulStop()
}

// authMiddleware is a middleware function that authenticates incoming requests using basic auth.
func authMiddleware(userService users.Service) func(context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		log.Debug("Authenticating request")
		token, err := grpcauth.AuthFromMD(ctx, "basic")
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "no basic header found: %v", err)
		}

		if userService.CheckPassword(token, token) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth credentials: %v", err)
		}

		return ctx, nil
	}
}

func interceptorLogger(l log.FieldLogger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make(map[string]any, len(fields)/2)
		i := logging.Fields(fields).Iterator()
		for i.Next() {
			k, v := i.At()
			f[k] = v
		}
		l := l.WithFields(f)

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg)
		case logging.LevelInfo:
			l.Info(msg)
		case logging.LevelWarn:
			l.Warn(msg)
		case logging.LevelError:
			l.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
