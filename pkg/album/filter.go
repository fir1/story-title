package album

import (
	"strings"
	"time"
)

func meetCriteriaPhotoAddress(photos []Photo, keyWord string) []Photo {
	var filteredPhotos []Photo
	for _, p := range photos {
		if isKeywordFoundInAddress(p.Address, keyWord) {
			filteredPhotos = append(filteredPhotos, p)
		}
	}
	return filteredPhotos
}

func isKeywordFoundInAddress(a Address, keyWord string) bool {
	line1Found := strings.Contains(a.Line1, keyWord)
	line2Found := strings.Contains(a.Line2, keyWord)
	cityFound := strings.Contains(a.City, keyWord)
	stateFound := strings.Contains(a.State, keyWord)
	postCodeFound := strings.Contains(a.PostCode, keyWord)
	countryFound := strings.Contains(a.Country, keyWord)
	countryCodeFound := strings.Contains(a.CountryCode, keyWord)

	return line1Found || line2Found || cityFound || stateFound || postCodeFound ||
		countryFound || countryCodeFound
}

func meetCriteriaPhotoDate(photos []Photo, dateFrom, dateTo *time.Time) []Photo {
	var filteredPhotos []Photo
	for _, p := range photos {
		switch {
		case (dateFrom != nil && dateTo != nil &&
			p.Date.After(*dateFrom) && p.Date.Before(*dateTo)) ||
			(dateFrom != nil && dateTo != nil &&
				p.Date.Equal(*dateFrom) && p.Date.Equal(*dateTo)):
			filteredPhotos = append(filteredPhotos, p)
		case (dateFrom != nil && p.Date.Equal(*dateFrom)) ||
			(dateFrom != nil && p.Date.After(*dateFrom)):
			filteredPhotos = append(filteredPhotos, p)
		case (dateTo != nil && p.Date.Equal(*dateTo)) ||
			(dateTo != nil && p.Date.Before(*dateTo)):
			filteredPhotos = append(filteredPhotos, p)
		}
	}
	return filteredPhotos
}

func meetCriteriaPhotoWeather(photos []Photo, w WeatherStatus) []Photo {
	var filteredPhotos []Photo
	for _, p := range photos {
		if p.WeatherStatus == w {
			filteredPhotos = append(filteredPhotos, p)
		}
	}
	return filteredPhotos
}
