package album

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ListAlbumsParams struct {
	Filter     *ListAlbumsFilter
	Sort       *ListAlbumsSort
	Pagination *Pagination
}

type ListAlbumsFilter struct {
	MinimumPhotoItems *uint64
	MaximumPhotoItems *uint64
}

type ListAlbumsSort struct {
	PhotoCounts *Sort
}

type ListAlbumsResponse struct {
	Albums     []Album
	Pagination PaginationResponse
}

func (s Service) ListAlbums(ctx context.Context, params ListAlbumsParams) (ListAlbumsResponse, error) {

	// get list of CSV files
	var files []string
	err := filepath.Walk(s.config.DataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return ListAlbumsResponse{}, err
	}

	var albums []Album

	for _, file := range files {
		id := strings.TrimRight(file, ".csv")

		sortDesc := SortDESC
		album, err := s.GetAlbum(ctx, GetAlbumParams{
			AlbumID: id,
			Sort:    &GetAlbumSort{PhotoDate: &sortDesc},
		})
		if err != nil {
			return ListAlbumsResponse{}, err
		}

		if (params.Filter != nil && params.Filter.MinimumPhotoItems != nil &&
			len(album.Photos) < int(*params.Filter.MinimumPhotoItems)) ||
			(params.Filter != nil && params.Filter.MaximumPhotoItems != nil &&
				len(album.Photos) > int(*params.Filter.MaximumPhotoItems)) {
			continue
		}

		albums = append(albums, album)
	}

	if params.Sort != nil && params.Sort.PhotoCounts != nil {
		if *params.Sort.PhotoCounts == SortASC {
			sort.SliceStable(albums, func(i, j int) bool {
				return len(albums[i].Photos) < len(albums[j].Photos)
			})
		} else {
			sort.SliceStable(albums, func(i, j int) bool {
				return len(albums[i].Photos) > len(albums[j].Photos)
			})
		}
	}

	paginationResponse := PaginationResponse{
		TotalItems: uint32(len(albums)),
	}

	if params.Pagination != nil {
		if int(params.Pagination.Offset) > len(albums) {
			return ListAlbumsResponse{}, nil
		}

		if params.Pagination.Offset > 0 {
			albums = append(albums[params.Pagination.Offset:])
		}
		if params.Pagination.Limit > 0 && int(params.Pagination.Limit) < len(albums) {
			albums = append(albums[:params.Pagination.Limit])
		}
	}

	paginationResponse.PageSize = uint32(len(albums))

	return ListAlbumsResponse{
		Albums:     albums,
		Pagination: paginationResponse,
	}, nil
}
