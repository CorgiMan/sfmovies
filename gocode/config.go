package sfmovies

import (
	"fmt"
	"strings"
)

const (
	APIVersion = "v1"
	HostName   = "http://corgiman.infty.nl"
	TableURL   = "https://data.sfgov.org/api/views/yitu-d5am/rows.csv?accessType=DOWNLOAD"
	OmdbURL    = "http://www.omdbapi.com/"
	MongoURL   = "192.168.59.103"
)

var (
	GeocodingKey       = "AIzaSyA8Px3Nesn6PsDhA0DIppHX16OEDT85WfA"
	GeocodingURL       = "https://maps.googleapis.com/maps/api/geocode/json?address="
	bounds             = fmt.Sprintf("bounds=%f,%f|%f,%f", MinLat, MinLng, MaxLat, MaxLng)
	GeocodingURLSuffix = ",+San+Francisco,+CA&" + bounds + "&key=" + GeocodingKey
)

// The size of the response of the queries handled by the API server.
const (
	NearQuerySize         = 20
	AutoCompleteQuerySize = 10
)

// San Francisco Bounds.
const (
	MinLat = 37.571026
	MinLng = -122.674607
	MaxLat = 37.928325
	MaxLng = -122.000637
)

// Usage string is used when root ("/") is requested.
var Usage = strings.Replace(fmt.Sprintf(
	`{
  "api_description": "San Francisco Movies API %s. Location and movie info of films recorded in San Francisco",
  "api_examples": {
    "{{.}}/status":                     "the status of the api server that handled the request",
    "{{.}}/movies/tt0028216":           "movie info of the specified IMDB ID",
    "{{.}}/complete?term=franc":        "auto complete results for the specified term parameter",
    "{{.}}/search?q=francisco":         "searches for movie title, film location, release year, director, production company, distributer, writer and actors",
    "{{.}}/near?lat=37.76&lng=-122.39": "searches for film locations near the presented gps coordinates"
    "{{.}}/?callback=XXX":              "use the callback parameter on any request to return JSONP in stead of just JSON"
  }
}`, APIVersion), "{{.}}", HostName, -1)
