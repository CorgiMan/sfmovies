package sfmovies

const (
	//
	TableURL           = "https://data.sfgov.org/api/views/yitu-d5am/rows.csv?accessType=DOWNLOAD"
	OmdbURL            = "http://www.omdbapi.com/"
	GeocodingKey       = "AIzaSyA8Px3Nesn6PsDhA0DIppHX16OEDT85WfA"
	GeocodingURL       = "https://maps.googleapis.com/maps/api/geocode/json?address="
	GeocodingURLSuffix = ",+San+Francisco,+CA&key=" + GeocodingKey
	MongoURL           = "192.168.59.103"
	// MongoUser          = "crazy_tango:firefly_bargain@"
	MongoUser = ""
	HostName  = "http://infty.nl:12000"

	// APIserver constants
	NearQuerySize         = 20
	AutoCompleteQuerySize = 10
)
