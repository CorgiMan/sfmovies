package sfmovies

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	tableURL           = "https://data.sfgov.org/api/views/yitu-d5am/rows.csv?accessType=DOWNLOAD"
	omdbURL            = "http://www.omdbapi.com/"
	geocodingKey       = "AIzaSyA8Px3Nesn6PsDhA0DIppHX16OEDT85WfA"
	geocodingURL       = "https://maps.googleapis.com/maps/api/geocode/json?address="
	geocodingURLSuffix = ",+San+Francisco,+CA&key=" + geocodingKey
)

// geocoding key in config file

func GetAndParseTable() (*APIData, error) {
	f, err := http.Get(tableURL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ad, err := ParseRows(csv.NewReader(f.Body))
	return ad, err
}

func ParseRows(r *csv.Reader) (*APIData, error) {
	ad := NewAPIData()

	// skip line 1. It contains the field descriptions
	_, err := r.Read()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for i := 0; i < 1100; i++ {
		fmt.Println(i)
		fields, err := r.Read()
		if err != nil {
			fmt.Println(err)
			break
		}

		movie, scene, err := ParseRow(fields)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if _, ok := ad.Movies[movie.IMDBID]; !ok {
			ad.Movies[movie.IMDBID] = movie
		}

		hasher := fnv.New32()
		bts, err := json.Marshal(scene)
		if err != nil {
			fmt.Println("Hashing error")
			continue
		}
		_, err = hasher.Write(bts)
		if err != nil {
			fmt.Println("Hashing error")
			continue
		}
		hash := fmt.Sprintf("%x", hasher.Sum32())
		ad.Scenes[hash] = scene
	}
	return ad, nil
}

func ParseRow(record []string) (*Movie, *Scene, error) {
	movie := new(Movie)
	scene := new(Scene)

	if len(record) < 2 {
		return movie, scene, errors.New("Not enough record fields")
	}

	title := record[0]
	movie, err := GetOMDBMovieInfo(title)
	if err != nil {
		return movie, scene, err
	}
	loc := record[2]
	location, err := GeoEncoding(loc)
	if err != nil {
		return movie, scene, err
	}
	location.Name = loc

	scene = &Scene{movie.IMDBID, location}

	return movie, scene, nil
}

func GeoEncoding(location string) (*Location, error) {
	// mi := Movie{}
	location = strings.Replace(location, " ", "+", -1)
	r, _ := http.Get(geocodingURL + location + geocodingURLSuffix)
	jsonbts, _ := ioutil.ReadAll(r.Body)

	response := struct {
		Results []struct {
			Geometry struct {
				Location *Location `json: location`
			} `json: geometry`
		} `json: results`
	}{}

	_ = json.Unmarshal(jsonbts, &response)
	rs := response.Results
	if len(rs) > 1 {
		fmt.Println("multiple geoencodings")
	}
	if len(rs) == 0 {
		fmt.Println("No results")
		return nil, errors.New("No geo encoding available")
	}
	return rs[0].Geometry.Location, nil
}

func GetOMDBMovieInfo(title string) (*Movie, error) {
	mi := new(Movie)
	title = strings.Replace(title, " ", "+", -1)
	r, err := http.Get(omdbURL + "/?t=" + title)
	if err != nil {
		return mi, err
	}
	jsonbts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return mi, err
	}
	err = json.Unmarshal(jsonbts, &mi)

	if err != nil {
		return mi, err
	}
	return mi, nil
}
