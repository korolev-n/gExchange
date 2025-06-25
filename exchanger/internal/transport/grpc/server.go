package grpc

import (
	"log/slog"
	"net"

	"github.com/korolev-n/gExchange/shared/api"
	"google.golang.org/grpc"
)

type Server struct {
	srv *grpc.Server
}

func NewGRPCServer(handler *Handler, logger *slog.Logger) *Server {
	s := grpc.NewServer()
	api.RegisterExchangerServiceServer(s, handler)
	return &Server{srv: s}
}

func (s *Server) Start(addr string, logger *slog.Logger) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("Failed to start gRPC listener", "error", err)
		return err
	}

	//logger.Info("Starting gRPC server", "addr", addr)
	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
