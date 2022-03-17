package service

import (
	pb "github.com/fir1/story-title/http/grpc/proto-gen"
	"github.com/fir1/story-title/pkg/album"
)

func albumsModelToPb(aa []album.Album) []*pb.Album {
	result := make([]*pb.Album, len(aa))
	for i, a := range aa {
		result[i] = albumModelToPb(a)
	}
	return result
}

func albumModelToPb(a album.Album) *pb.Album {
	return &pb.Album{Photos: photosModelToPb(a.Photos)}
}

func photosModelToPb(pp []album.Photo) []*pb.Photo {
	result := make([]*pb.Photo, len(pp))
	for i, p := range pp {
		result[i] = photoModelToPb(p)
	}
	return result
}

func photoModelToPb(p album.Photo) *pb.Photo {
	return &pb.Photo{
		PhotoDate: p.Date.String(),
		Latitude:  float32(p.Latitude),
		Longitude: float32(p.Longitude),
		Titles:    p.Titles,
		Weather:   weatherStatusModelToPb(p.WeatherStatus),
		Address:   addressModelToPb(p.Address),
	}
}

func weatherStatusModelToPb(w album.WeatherStatus) pb.Weather {
	switch w {
	case album.WeatherStatusRainy:
		return pb.Weather_WEATHER_RAINY
	case album.WeatherStatusSnowy:
		return pb.Weather_WEATHER_SNOWY
	case album.WeatherStatusSunny:
		return pb.Weather_WEATHER_SUNNY
	}
	return pb.Weather_WEATHER_UNDEFINED
}

func weatherStatusPbToModel(w pb.Weather) album.WeatherStatus {
	switch w {
	case pb.Weather_WEATHER_RAINY:
		return album.WeatherStatusRainy
	case pb.Weather_WEATHER_SNOWY:
		return album.WeatherStatusSnowy
	case pb.Weather_WEATHER_SUNNY:
		return album.WeatherStatusSunny
	}
	return album.WeatherStatusUndefined
}

func addressModelToPb(a album.Address) *pb.Address {
	return &pb.Address{
		Line1:       a.Line1,
		Line2:       a.Line2,
		City:        a.City,
		PostCode:    a.PostCode,
		Country:     a.Country,
		CountryCode: a.CountryCode,
		State:       a.State,
	}
}
