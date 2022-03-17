package album

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/fir1/story-title/pkg/utilerror"
	"googlemaps.github.io/maps"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GetAlbumParams struct {
	AlbumID string
	Sort    *GetAlbumSort
	Filter  *GetAlbumFilter
}

type GetAlbumSort struct {
	PhotoDate    *Sort
	PhotoCountry *Sort
}

type GetAlbumFilter struct {
	PhotoAddressKeyword *string
	PhotoDateFrom       *time.Time
	PhotoDateTo         *time.Time
	WeatherStatus       *WeatherStatus
}

func (s Service) GetAlbum(ctx context.Context, params GetAlbumParams) (Album, error) {
	if params.AlbumID == "" {
		return Album{}, utilerror.ErrArgument{Wrapped: fmt.Errorf("album_id must not be empty")}
	}

	// we are searching existing albums from the memory cache before making further API calls
	// to third party this is to prevent from hitting free tier of API provider.
	album, ok := s.mapIDToAlbum[params.AlbumID]
	if !ok {
		f, err := os.Open(fmt.Sprintf("%s/%s.csv", s.config.DataDir, params.AlbumID))
		if err != nil {
			return Album{}, err
		}
		defer f.Close()

		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		data, err := csvReader.ReadAll()
		if err != nil {
			return Album{}, err
		}

		// convert records to array of structs
		photos, err := createPhotosFromCsvData(data)
		if err != nil {
			return Album{}, err
		}

		// generate story titles for a photo
		var wg sync.WaitGroup
		errChan := make(chan error)
		done := make(chan struct{})

		for i, p := range photos {
			wg.Add(1)

			go func(p Photo, i int) {
				defer wg.Done()

				resp, err := s.getPhotoStoryTitles(ctx, p)
				if err != nil {
					errChan <- err
					return
				}

				photos[i].Titles = resp.Titles
				photos[i].Address = resp.Address
			}(p, i)
		}

		go func() {
			wg.Wait()
			done <- struct{}{}
			close(errChan)
		}()

		select {
		case <-ctx.Done():
			return Album{}, ctx.Err()
		case err = <-errChan:
			return Album{}, err
		case <-done:
			break
		}

		album = Album{Photos: photos}
		s.mapIDToAlbum[params.AlbumID] = album
	}

	photos := album.Photos

	if params.Filter != nil {
		photos = filterPhotos(photos, *params.Filter)
	}

	// user can ask to sort photos taken date oldest to newest ascending
	if params.Sort != nil {
		photos = sortPhotos(photos, *params.Sort)
	}

	return Album{
		Photos: photos,
	}, nil
}

func createPhotosFromCsvData(data [][]string) ([]Photo, error) {
	var photos []Photo

	for _, line := range data {
		var photo Photo
		for column, field := range line {
			switch column {
			case 0:
				// the layout of date in CSV files are different so if we want to use native
				// GoLang time package we have to know which kind of date layout they have before
				// parsing it to time.Time object.
				// So decided to use third party package for doing this parsing.
				date, err := dateparse.ParseAny(field)
				if err != nil {
					return nil, err
				}

				photo.Date = date
			case 1:
				latitude, err := strconv.ParseFloat(field, 64)
				if err != nil {
					return nil, err
				}

				photo.Latitude = latitude
			case 2:
				longitude, err := strconv.ParseFloat(field, 64)
				if err != nil {
					return nil, err
				}

				photo.Longitude = longitude
			default:
				return nil, fmt.Errorf("unexpected field found in row %d", column)
			}
		}
		photos = append(photos, photo)
	}
	return photos, nil
}

type getPhotoStoryTitlesResponse struct {
	Titles  []string
	Address Address
}

func (s Service) getPhotoStoryTitles(ctx context.Context, p Photo) (getPhotoStoryTitlesResponse, error) {
	c, err := maps.NewClient(maps.WithAPIKey(s.config.GoogleMapsAPIKey))
	if err != nil {
		return getPhotoStoryTitlesResponse{}, err
	}

	response, err := c.ReverseGeocode(ctx, &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: p.Latitude,
			Lng: p.Longitude,
		},
	})
	if err != nil {
		return getPhotoStoryTitlesResponse{}, err
	}

	var geocodingResult maps.GeocodingResult
	for i, r := range response {
		if r.Geometry.Location.Lat == p.Latitude && r.Geometry.Location.Lng == p.Longitude {
			geocodingResult = r
			break
		}

		if i == len(response)-1 {
			geocodingResult = r
		}
	}

	var titles []string
	var address Address
	for _, addressComponent := range geocodingResult.AddressComponents {
		for _, t := range addressComponent.Types {
			switch t {
			case "street_number": // 532
				address.Line1 = addressComponent.LongName
			case "route": // LaGuardia Place
				address.Line1 = strings.Trim(fmt.Sprintf(
					"%s %s", address.Line1, addressComponent.LongName), " ")
			case "sublocality": // Manhattan
				if isWeekend(p.Date) {
					titles = append(titles,
						fmt.Sprintf("A weekend in %s", addressComponent.LongName))
				}
				address.Line2 = addressComponent.LongName
			case "locality": // New York
				if isWeekend(p.Date) {
					titles = append(titles,
						fmt.Sprintf("A weekend in %s", addressComponent.LongName))
				}
				titles = append(titles,
					fmt.Sprintf("A trip to %s", addressComponent.LongName),
					fmt.Sprintf("%s in %s", addressComponent.LongName, p.Date.Month()))

				address.City = addressComponent.LongName
			case "country": // United States
				titles = append(titles,
					fmt.Sprintf("A trip to the %s", addressComponent.LongName),
					fmt.Sprintf("%s in %s", addressComponent.LongName, p.Date.Month()))

				address.Country = addressComponent.LongName
				address.CountryCode = addressComponent.ShortName
			case "postal_code": // 10012
				address.PostCode = addressComponent.LongName
			case "administrative_area_level_1":
				address.State = addressComponent.LongName
			}
		}
	}

	return getPhotoStoryTitlesResponse{
		Titles:  titles,
		Address: address,
	}, nil
}

func isWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

func sortPhotos(pp []Photo, s GetAlbumSort) []Photo {
	if s.PhotoDate != nil {
		switch *s.PhotoDate {
		case SortDESC:
			sort.Slice(pp, func(i, j int) bool {
				return pp[i].Date.After(pp[j].Date)
			})
		case SortASC:
			sort.Slice(pp, func(i, j int) bool {
				return pp[i].Date.Before(pp[j].Date)
			})
		}
	}

	if s.PhotoCountry != nil {
		switch *s.PhotoCountry {
		case SortDESC:
			sort.Slice(pp, func(i, j int) bool {
				return pp[i].Address.Country > pp[j].Address.Country
			})
		case SortASC:
			sort.Slice(pp, func(i, j int) bool {
				return pp[i].Address.Country < pp[j].Address.Country
			})
		}
	}

	return pp
}

func filterPhotos(pp []Photo, f GetAlbumFilter) []Photo {
	if f.PhotoAddressKeyword != nil {
		pp = meetCriteriaPhotoAddress(pp, *f.PhotoAddressKeyword)
	}

	if f.PhotoDateFrom != nil || f.PhotoDateTo != nil {
		pp = meetCriteriaPhotoDate(pp, f.PhotoDateFrom, f.PhotoDateTo)
	}

	if f.WeatherStatus != nil {
		pp = meetCriteriaPhotoWeather(pp, *f.WeatherStatus)
	}
	return pp
}
