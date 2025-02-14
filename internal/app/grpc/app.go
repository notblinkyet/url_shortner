package grpcapp

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/notblinkyet/url_shortner/internal/services"
	urlshortner "github.com/notblinkyet/url_shortner/internal/transport/grpc/url_shortner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log    *log.Logger
	server *grpc.Server
	host   string
	port   int
}

func New(log *log.Logger, services services.IServices, port int, host string, timeout time.Duration) *App {

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived,
			logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Printf("recovery of panic: %v", p)
			return status.Errorf(codes.Internal, "interanal error")
		}),
	}

	gRRPCserver := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
		),
	)
	urlshortner.Register(gRRPCserver, services, timeout)

	return &App{
		log:    log,
		server: gRRPCserver,
		host:   host,
		port:   port,
	}
}

func InterceptorLogger(l *log.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Println(msg)
	})
}

func (app *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", app.host, app.port))
	if err != nil {
		return err
	}
	app.log.Printf("grpc server started\nhost: %s\nport: %d\n", app.host, app.port)

	if err := app.server.Serve(l); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (app *App) Stop() {
	app.log.Println("stopping grpc server")

	app.server.GracefulStop()
}

func (app *App) MustRun() {
	if err := app.Run; err != nil {
		panic(err)
	}
}
