package internalgrpc

import (
	"context"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/seth2810/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
)

type Server struct {
	address string
	logger  Logger
	server  *grpc.Server
	service pb.CalendarServiceServer
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, id, title string) error
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func NewServer(address string, logger Logger, app Application) *Server {
	return &Server{address, logger, nil, &calendarServiceServer{app: app}}
}

func (s *Server) Start(ctx context.Context) error {
	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			loggingInterceptor(s.logger),
			grpc_validator.UnaryServerInterceptor(),
		)),
	)

	pb.RegisterCalendarServiceServer(s.server, s.service)

	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	if err := s.server.Serve(lis); err != nil {
		return err
	}

	<-ctx.Done()

	return s.Stop()
}

func (s *Server) Stop() error {
	s.server.GracefulStop()

	return nil
}
