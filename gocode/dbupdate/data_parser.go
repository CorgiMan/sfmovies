package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/CorgiMan/sfmovies/gocode"
)

// geocoding key in config file

func GetAndParseAPIData() (*sfmovies.APIData, error) {
	f, err := http.Get(sfmovies.TableURL)
	if err != nil {
		log.Fatal(err)
	}

	ad, err := ParseRows(csv.NewReader(f.Body))
	return ad, err
}

func ParseRows(r *csv.Reader) (*sfmovies.APIData, error) {
	ad := sfmovies.NewAPIData()

	// skip line 1. It contains the field descriptions
	_, err := r.Read()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for i:=0; i<10; i++ {
		fields, err := r.Read()
		if err != nil {
			break
		}

		movie, scene, err := ParseRow(fields)
		if err != nil {
			log.Println(err)
			continue
		}
		if _, ok := ad.Movies[movie.IMDBID]; !ok {
			ad.Movies[movie.IMDBID] = movie
		}

		hasher := fnv.New32()
		bts, err := json.Marshal(scene)
		if err != nil {
			log.Println("Hashing error")
			continue
		}
		_, err = hasher.Write(bts)
		if err != nil {
			log.Println("Hashing error")
			continue
		}
		hash := fmt.Sprintf("%x", hasher.Sum32())
		ad.Scenes[hash] = scene
	}
	ad.Time = time.Now()
	fmt.Println(ad.Time)

	return ad, nil
}

func ParseRow(record []string) (*sfmovies.Movie, *sfmovies.Scene, error) {
	movie := new(sfmovies.Movie)
	scene := new(sfmovies.Scene)

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

	scene = &sfmovies.Scene{movie.IMDBID, location}

	return movie, scene, nil
}

func GeoEncoding(location string) (*sfmovies.Location, error) {
	location = strings.Replace(location, " ", "+", -1)
	r, _ := http.Get(sfmovies.GeocodingURL + location + sfmovies.GeocodingURLSuffix)
	jsonbts, _ := ioutil.ReadAll(r.Body)

	response := struct {
		Results []struct {
			Geometry struct {
				Location *sfmovies.Location `json: location`
			} `json: geometry`
		} `json: results`
	}{}

	_ = json.Unmarshal(jsonbts, &response)
	rs := response.Results
	if len(rs) > 1 {
		log.Println("multiple geoencodings")
	}
	if len(rs) == 0 {
		log.Println("No results")
		return nil, errors.New("No geo encoding available")
	}
	return rs[0].Geometry.Location, nil
}

func GetOMDBMovieInfo(title string) (*sfmovies.Movie, error) {
	mi := new(sfmovies.Movie)
	title = strings.Replace(title, " ", "+", -1)
	r, err := http.Get(sfmovies.OmdbURL + "/?t=" + title)
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
