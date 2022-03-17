package dependency

import (
	"github.com/fir1/story-title/pkg/album"
	"github.com/fir1/story-title/pkg/config"
	"github.com/fir1/story-title/pkg/logger"
	"go.uber.org/fx"
)

var CommonFxOptions = fx.Options(
	album.Module,
	config.Module,
	logger.Module,
)
