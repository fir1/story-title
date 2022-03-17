package service

import (
	"google.golang.org/grpc"
)

type Handler struct {
	albumHandler *AlbumHandler
}

func NewHandler(ah *AlbumHandler) *Handler {
	return &Handler{
		albumHandler: ah,
	}
}

func (h *Handler) RegisterHandler(s *grpc.Server) {
	h.albumHandler.registerAlbumService(s)
}
