package dependency

import (
	"github.com/fir1/story-title/http/grpc"
	"github.com/fir1/story-title/http/grpc/service"
	"go.uber.org/fx"
)

var GrpcServerFxOptions = fx.Options(
	CommonFxOptions,
	service.Module,
	grpc.Module,
)
