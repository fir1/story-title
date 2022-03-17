package album

import "time"

type WeatherStatus uint8

const (
	WeatherStatusUndefined WeatherStatus = iota
	WeatherStatusSunny
	WeatherStatusRainy
	WeatherStatusSnowy
)

type Sort string

const (
	SortASC  Sort = "ASC"
	SortDESC Sort = "DESC"
)

type Address struct {
	Line1       string
	Line2       string
	City        string
	State       string
	PostCode    string
	Country     string
	CountryCode string
}

type Photo struct {
	Date          time.Time
	Latitude      float64
	Longitude     float64
	Titles        []string
	WeatherStatus WeatherStatus
	Address       Address
}

type Album struct {
	Photos []Photo
}
