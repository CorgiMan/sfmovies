package sfmovies

import (
	"math"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type Location struct {
	Name string
	Lat  float64 `json: lat`
	Lng  float64 `json: lng`
}

func (l1 *Location) Distance(l2 *Location) float64 {
	dLat := l2.Lat - l1.Lat
	dLng := l2.Lng - l1.Lng

	latrad := (l1.Lat + l2.Lat) / 180 * math.Pi / 2
	dX := dLng * math.Cos(latrad)
	dY := dLat
	return math.Sqrt(dX*dX + dY*dY)
}

func (loc *Location) IsInBounds() bool {
	return (MinLat <= loc.Lat && loc.Lat <= MaxLat) &&
		(MinLng <= loc.Lng && loc.Lng <= MaxLng)
}

type Movie struct {
	Title    string
	Year     string
	Rated    string
	Released string
	Runtime  string
	Genre    string
	Director string
	Writer   string
	Actors   string
	Plot     string
	Poster   string
	IMDBID   string `json: imdbID`
}

type Scene struct {
	IMDBID string
	*Location
}

type APIData struct {
	Time   time.Time
	Movies map[string]*Movie
	Scenes map[string]*Scene
}

func NewAPIData() *APIData {
	r := &APIData{}
	r.Movies = make(map[string]*Movie)
	r.Scenes = make(map[string]*Scene)
	return r
}

func GetLatestAPIData() (*APIData, error) {
	session, err := mgo.Dial(MongoURL)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	c := session.DB("apiserverdb").C("apiserverdata")

	var ad APIData
	err = c.Find(nil).Sort("-time").One(&ad)
	if err != nil {
		return nil, err
	}
	return &ad, nil
}

func StoreAPIData(ad *APIData) error {
	session, err := mgo.Dial(MongoURL)
	if err != nil {
		return err
	}
	defer session.Close()

	c := session.DB("apiserverdb").C("apiserverdata")
	err = c.Insert(ad)
	if err != nil {
		return err
	}

	return nil
}
