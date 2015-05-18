package sfmovies

import "fmt"

const (
	MinLat = 37.571026
	MinLng = -122.674607
	MaxLat = 37.928325
	MaxLng = -122.000637

	NearQuerySize         = 20
	AutoCompleteQuerySize = 10
)
const (
	APIVersion   = "v1"
	TableURL     = "https://data.sfgov.org/api/views/yitu-d5am/rows.csv?accessType=DOWNLOAD"
	OmdbURL      = "http://www.omdbapi.com/"
	GeocodingKey = "AIzaSyA8Px3Nesn6PsDhA0DIppHX16OEDT85WfA"
	MongoURL     = "localhost"
	HostName     = "http://corgiman.infty.nl"
)

var (
	GeocodingURL       = "https://maps.googleapis.com/maps/api/geocode/json?address="
	bounds             = fmt.Sprintf("bounds=%s,%s|%s,%s", MinLat, MinLng, MaxLat, MaxLng)
	GeocodingURLSuffix = ",+San+Francisco,+CA&" + bounds + "&key=" + GeocodingKey
)
