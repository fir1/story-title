package grpc

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func (s *Server) logRequest(req interface{}) {
	if s.config.Log.EnableDetailedRequest {
		s.logger.Infof("Request - Data{%+v}", req)
	} else {
		s.logger.Infof("Request - Data{%q}", "detailed request logger disabled")
	}
}

func (s *Server) logResponse(ctx context.Context, start time.Time, info *grpc.UnaryServerInfo, res interface{}, err error) {
	if s.config.Log.EnableDetailedResponse {
		s.logger.Infof(`Response - Data{%v} Method:%s Duration:%s Session-ID:%v \n`, res, info.FullMethod, time.Since(start), ctx.Value("session"))
	} else {
		s.logger.Infof(`Response - Data{%q}, Method:%s, Status: OK, Duration:%s Session-ID:%v \n`,
			"detailed response logger disabled", info.FullMethod, time.Since(start), ctx.Value("session"))
	}
}
