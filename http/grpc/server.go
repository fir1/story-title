package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/fir1/story-title/http/grpc/service"
	"github.com/fir1/story-title/pkg/config"
	"github.com/fir1/story-title/pkg/utilerror"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

const (
	MaxGrpcMessageSize = 100 * 1024 * 1024
)

type Server struct {
	grpcServer         *grpc.Server
	grpcServiceHandler *service.Handler
	config             config.Config
	logger             *logrus.Logger
}

func NewServer(cnf config.Config, lg *logrus.Logger, gsh *service.Handler) *Server {
	return &Server{
		grpcServiceHandler: gsh,
		config:             cnf,
		logger:             lg,
	}
}

func (s *Server) RunServer() error {
	s.newGrpcServer()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPCPort))
	if err != nil {
		return fmt.Errorf("listen port %d error: %w", s.config.GRPCPort, err)
	}

	s.logger.Infof("server running on port: %v", s.config.GRPCPort)

	go func() {
		err = s.grpcServer.Serve(lis)
		if err != nil {
			s.logger.Errorf("grpc-server start error: %v", err)
		}
	}()

	return err
}

func (s *Server) StopServer() error {
	if s == nil || s.grpcServer == nil {
		return errors.New("grpc server can not be nil")
	}

	// graceful shutdown
	s.grpcServer.GracefulStop()
	s.logger.Warn("grpc-server graceful stop")
	return nil
}

func (s *Server) newGrpcServer() {
	server := grpc.NewServer(s.withServerUnaryInterceptord(), grpc.MaxRecvMsgSize(MaxGrpcMessageSize))
	s.grpcServiceHandler.RegisterHandler(server)
	s.grpcServer = server
}

func (s *Server) withServerUnaryInterceptord() grpc.ServerOption {
	return grpc.UnaryInterceptor(s.interceptor)
}

func (s *Server) interceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	if info.FullMethod == "/grpc.health.v1.Health/Check" {
		h, err := handler(ctx, req)
		if err != nil {
			logrus.WithFields(logrus.Fields{"method": info.FullMethod}).Errorf("Health error: %v (duration %v)", err, time.Since(start))
		}

		return h, err
	}

	s.logRequest(req)

	h, err := handler(ctx, req)
	if err != nil {
		return nil, utilerror.ConversionToGrpcDetailedError(err)
	}

	s.logResponse(ctx, start, info, h, err)
	return h, nil
}
