package album

import (
	"github.com/fir1/story-title/pkg/config"
)

type Service struct {
	config       config.Config
	mapIDToAlbum map[string]Album
}

func NewAlbumService(cnf config.Config) Service {
	return Service{
		config:       cnf,
		mapIDToAlbum: make(map[string]Album),
	}
}
