package album

import "go.uber.org/fx"

var Module = fx.Provide(
	NewAlbumService,
)
