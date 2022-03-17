package main

import (
	"context"
	"github.com/fir1/story-title/http/grpc"
	"github.com/fir1/story-title/pkg/dependency"
	"go.uber.org/fx"
	"log"
)

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	app := fx.New(
		dependency.GrpcServerFxOptions,
		fx.Invoke(GrpcServerRun),
	)
	app.Run()

	return app.Err()
}

func GrpcServerRun(s *grpc.Server, lifecycle fx.Lifecycle) error {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return s.RunServer()
			},
			OnStop: func(ctx context.Context) error {
				return s.StopServer()
			},
		},
	)

	return nil
}
