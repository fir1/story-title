package service

import (
	"context"
	"github.com/araddon/dateparse"
	pb "github.com/fir1/story-title/http/grpc/proto-gen"
	"github.com/fir1/story-title/pkg/album"
	"google.golang.org/grpc"
	// pb "github.com/fortuwealth/protorepo-front-office-go/session"
)

const DefaultOriginIP = "0.0.0.0"

type AlbumHandler struct {
	pb.UnimplementedAlbumServiceServer
	albumSvc album.Service
}

func NewAlbumHandler(albumSvc album.Service) *AlbumHandler {
	return &AlbumHandler{
		albumSvc: albumSvc,
	}
}

func (ah *AlbumHandler) GetAlbum(ctx context.Context, req *pb.GetAlbumRequest) (*pb.GetAlbumResponse, error) {
	var sort *album.GetAlbumSort
	if req.Sort != nil {
		sort = &album.GetAlbumSort{}
		if req.Sort.PhotoDate != nil && *req.Sort.PhotoDate != pb.Order_UNDEFINED {
			st := album.Sort(req.Sort.GetPhotoDate().String())
			sort.PhotoDate = &st
		}

		if req.Sort.PhotoCountry != nil && *req.Sort.PhotoCountry != pb.Order_UNDEFINED {
			st := album.Sort(req.Sort.GetPhotoCountry().String())
			sort.PhotoCountry = &st
		}
	}

	var filter *album.GetAlbumFilter
	if req.Filter != nil {
		filter = &album.GetAlbumFilter{
			PhotoAddressKeyword: req.Filter.PhotoAddressKeyword,
		}

		if req.Filter.PhotoDateTo != nil {
			date, err := dateparse.ParseAny(*req.Filter.PhotoDateTo)
			if err != nil {
				return nil, err
			}
			filter.PhotoDateTo = &date
		}

		if req.Filter.PhotoDateFrom != nil {
			date, err := dateparse.ParseAny(*req.Filter.PhotoDateFrom)
			if err != nil {
				return nil, err
			}
			filter.PhotoDateFrom = &date
		}

		if req.Filter.Weather != nil {
			ws := weatherStatusPbToModel(*req.Filter.Weather)
			filter.WeatherStatus = &ws
		}
	}

	albums, err := ah.albumSvc.GetAlbum(ctx, album.GetAlbumParams{
		AlbumID: req.AlbumId,
		Sort:    sort,
		Filter:  filter,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetAlbumResponse{
		Album: albumModelToPb(albums),
	}, nil
}

func (ah *AlbumHandler) ListAlbums(ctx context.Context, req *pb.ListAlbumsRequest) (*pb.ListAlbumsResponse, error) {
	var sort *album.ListAlbumsSort
	if req.Sort != nil && req.Sort.PhotoCounts != nil &&
		*req.Sort.PhotoCounts != pb.Order_UNDEFINED {
		st := album.Sort(req.Sort.GetPhotoCounts().String())
		sort = &album.ListAlbumsSort{PhotoCounts: &st}
	}

	var filter *album.ListAlbumsFilter
	if req.Filter != nil {
		filter = &album.ListAlbumsFilter{
			MinimumPhotoItems: req.Filter.MinimumPhotoItems,
			MaximumPhotoItems: req.Filter.MaximumPhotoItems,
		}
	}

	var pagination *album.Pagination
	if req.Pagination != nil {
		pagination = &album.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}

	listAlbumsResponse, err := ah.albumSvc.ListAlbums(ctx, album.ListAlbumsParams{
		Filter:     filter,
		Sort:       sort,
		Pagination: pagination,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ListAlbumsResponse{
		Albums: albumsModelToPb(listAlbumsResponse.Albums),
		Pagination: &pb.PaginationResponse{
			Total:    listAlbumsResponse.Pagination.TotalItems,
			PageSize: listAlbumsResponse.Pagination.PageSize,
		},
	}, nil
}

func (ah *AlbumHandler) registerAlbumService(s *grpc.Server) {
	pb.RegisterAlbumServiceServer(s, ah)
}
