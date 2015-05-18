package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CorgiMan/sfmovies/gocode"
)

type Status struct {
	APIVersion   string
	RunningSince time.Time
	DataVersion  time.Time
}

type Error struct {
	Error string
}

var (
	port   = flag.String("port", "80", "port that program listens on")
	ad     *sfmovies.APIData
	trie   *TrieNode
	status Status

	usage = strings.Replace(fmt.Sprintf(
		`{
		  "api_description": "San Francisco Movies Api %s. Location and movie info of films recorded in San Francisco",
		  "api_examples": {
		    "{{.}}/status":                     "the status of the api server that handled the request",
		    "{{.}}/movies/tt0028216":           "movie info of the specified IMDB ID",
		    "{{.}}/complete?term=franc":        "auto complete results for the specified term parameter",
		    "{{.}}/search?q=francisco":         "searches for movie title, film location, release year, director, production company, distributer, writer and actors",
		    "{{.}}/near?lat=37.76&lng=-122.39": "searches for film locations near the presented gps coordinates"
		    "{{.}}/?callback=XXX":              "use the callback parameter on any request to return JSONP in stead of just JSON"
		  }
		}`, sfmovies.APIVersion), "{{.}}", sfmovies.HostName, -1)
)

func init() {
	flag.Parse()

	var err error
	ad, err = sfmovies.GetLatestAPIData()
	if err != nil {
		log.Fatal(err)
	}
	if ad == nil {
		log.Fatal(errors.New("No API data received from mongodb"))
	}
	status = Status{}
	status.APIVersion = sfmovies.APIVersion
	status.RunningSince = time.Now()
	status.DataVersion = ad.Time

	trie = CreateTrie(ad)
}

func main() {
	// root handles near, search and complete queries as well as api description
	http.HandleFunc("/", jsonpHandler(rootHandler))
	http.HandleFunc("/movies/", jsonpHandler(moviesHandler))
	http.HandleFunc("/status", jsonpHandler(statusHandler))
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/near":
		nearHandler(w, r)
	case "/search":
		searchHandler(w, r)
	case "/complete":
		completeHandler(w, r)
	case "/":
		_, err := io.WriteString(w, usage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func writeResult(w http.ResponseWriter, v interface{}) {
	bts, err1 := json.MarshalIndent(v, "", "  ")
	_, err2 := w.Write(bts)
	if err1 != nil || err2 != nil {
		http.Error(w, "failed to marshal and write json", http.StatusInternalServerError)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	writeResult(w, status)
}

func moviesHandler(w http.ResponseWriter, r *http.Request) {
	imdbid := r.URL.Path[len("/movies/"):]
	if movie, ok := ad.Movies[imdbid]; ok {
		writeResult(w, movie)
	} else {
		writeResult(w, Error{"Recource not found"})
	}
}

func completeHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("term")
	result := trie.GetFrom(q, sfmovies.AutoCompleteQuerySize)
	writeResult(w, result)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	if result := trie.Get(ad, q); result != nil {
		writeResult(w, result)
	} else {
		writeResult(w, Error{"Recource not found"})
	}
}

func nearHandler(w http.ResponseWriter, r *http.Request) {
	lat, err1 := strconv.ParseFloat(r.FormValue("lat"), 64)
	lng, err2 := strconv.ParseFloat(r.FormValue("lng"), 64)
	if err1 != nil || err2 != nil {
		http.Error(w, "failed to parse lat lng parameters", http.StatusInternalServerError)
		return
	}
	loc := sfmovies.Location{"", lat, lng}
	ds := make([]float64, 0)
	scs := make([]*sfmovies.Scene, 0)
	for _, scene := range ad.Scenes {
		ds = append(ds, loc.Distance(scene.Location))
		scs = append(scs, scene)
	}

	// select closest element NearQuerySize times
	result := make([]*sfmovies.Scene, 0)
	for i := 0; i < sfmovies.NearQuerySize; i++ {
		ix := minix(ds)
		if ix == -1 {
			break
		}
		result = append(result, scs[ix])
		ds[ix] = ds[len(ds)-1]
		scs[ix] = scs[len(scs)-1]
		ds = ds[:len(ds)-1]
		scs = scs[:len(scs)-1]
	}

	writeResult(w, result)
}

func minix(a []float64) int {
	if len(a) == 0 {
		return -1
	}
	mini := 0
	for i := range a {
		if a[i] < a[mini] {
			mini = i
		}
	}
	return mini
}

type Handler func(http.ResponseWriter, *http.Request)

// wraps around a handler and adds JSONP padding if the callback parameter is set
func jsonpHandler(fn Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		callback := r.FormValue("callback")
		if callback != "" {
			w.Header().Add("Content-Type", "application/javascript")
			_, err1 := w.Write([]byte(callback + "("))
			fn(w, r)
			_, err2 := w.Write([]byte(");"))
			if err1 != nil || err2 != nil {
				http.Error(w, "failed to write JSONP padding", http.StatusInternalServerError)
			}
		} else {
			fn(w, r)
		}

	}
}
