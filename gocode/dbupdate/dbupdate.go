// Updates the database with the latest API Data. This should be run once a day to keep the database up to date.
// Fetches the source table, consults OMDB and Google Geo-Encoding API, and stores data in MongoDB
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

func main() {
	apidata, err := GetAndParseAPIData()
	if err != nil {
		log.Fatal(err)
	}

	sfmovies.StoreAPIData(apidata)
}

// Fetches source table
func GetAndParseAPIData() (*sfmovies.APIData, error) {
	f, err := http.Get(sfmovies.TableURL)
	if err != nil {
		log.Fatal(err)
	}

	ad, err := ParseRows(csv.NewReader(f.Body))
	return ad, err
}

// For each row/scene in the source table we add it to the APIData.Scenes,
// if we encounter a new movie we store it in APIData.Movies.
func ParseRows(r *csv.Reader) (*sfmovies.APIData, error) {
	ad := sfmovies.NewAPIData()

	// skip line 1. It contains the field descriptions
	_, err := r.Read()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Println("Parsing...")
	// 30 is an arbitrary limit set to not exhaust API key (in production set to endless)
	for i := 0; i < 30; i++ {
		if i%100 == 99 {
			fmt.Println(i+1, "done")
		}
		fields, err := r.Read()
		if err != nil {
			// No more lines in source table
			break
		}

		movie, scene, err := ParseRow(fields)
		if err != nil {
			log.Println(err)
			continue
		}

		// check if movie is new
		if _, ok := ad.Movies[movie.IMDBID]; !ok {
			ad.Movies[movie.IMDBID] = movie
		}

		// calculate a new scene id by hashing
		hasher := fnv.New32()
		bts, err := json.Marshal(scene)
		if err != nil {
			log.Println("Hashing error: failed to marshal", scene)
			continue
		}
		_, err = hasher.Write(bts)
		if err != nil {
			log.Println("Hashing error", string(bts))
			continue
		}
		hash := fmt.Sprintf("%x", hasher.Sum32())
		ad.Scenes[hash] = scene
	}
	fmt.Println("Done parsing")
	ad.Time = time.Now()

	return ad, nil
}

// Finds movie and location information from the fields in a table.
func ParseRow(record []string) (*sfmovies.Movie, *sfmovies.Scene, error) {
	movie := new(sfmovies.Movie)
	scene := new(sfmovies.Scene)

	if len(record) < 3 {
		return movie, scene, errors.New("Not enough record fields in " + strings.Join(record, " "))
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

// Consults Google Geo-Encoding API to find the coords of the location name somewhere near San Francisco.
func GeoEncoding(location string) (*sfmovies.Location, error) {
	location = strings.Replace(location, " ", "+", -1)
	location = strings.Replace(location, "&", "", -1)
	r, err := http.Get(sfmovies.GeocodingURL + location + sfmovies.GeocodingURLSuffix)
	for err != nil {
		return nil, err
	}

	jsonbts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	response := struct {
		Results []struct {
			Geometry struct {
				Location *sfmovies.Location `json: location`
			} `json: geometry`
		} `json: results`
	}{}

	err = json.Unmarshal(jsonbts, &response)
	if err != nil {
		fmt.Println(string(jsonbts))
		return nil, err
	}
	rs := response.Results
	for _, r := range rs {
		if r.Geometry.Location.IsInBounds() {
			return r.Geometry.Location, nil
		}
	}
	return nil, errors.New("No geo encoding available for " + location)
}

// Consults OMDB API to get the movie info.
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
