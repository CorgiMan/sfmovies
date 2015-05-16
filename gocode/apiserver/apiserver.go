package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CorgiMan/sfmovies/gocode"
)

type Status struct {
	Machine      int
	RunningSince time.Time
	DataVersion  string
	Running      bool
}

var port = flag.String("port", "8080", "port that program listens on")
var machine = flag.String("machine", "1", "port that program listens on")

var ad *sfmovies.APIData
var trie *TrieNode
var status = Status{}

var usage = strings.Replace(
	`{
  "api_description": "sf movies api. Location and movie info of movies recorded in San Francisco",
  "api_examples": {
    "{{.}}/status":               "returns the status of the api server that handled the request",
    "{{.}}/movies/imdb_id":       "movie info",
    "{{.}}/scenes/scene_id":      "movie info",
    "{{.}}/complete?term=###":    "Auto complete results for query",
    "{{.}}/search?q=###":         "Search for movie title, film location, release year, director, production company, distributer, writer and actors",
    "{{.}}/near?lat=###&lng=###": "Search for film locations near the presented gps coordinates"
  }
}`, "{{.}}", sfmovies.HostName, -1)

func init() {
	flag.Parse()

	var err error
	ad, err = sfmovies.GetLatestAPIData()
	if err != nil {
		fmt.Println(err)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	status.Machine = rand.Int()
	status.RunningSince = time.Now()
	status.DataVersion = "1"
	status.Running = true

	trie = CreateTrie(ad)
}

//TODO: Cache movies and scenes

func moviesHandler(w http.ResponseWriter, r *http.Request) {
	imdbid := r.URL.Path[len("/movies/"):]
	movie, ok := ad.Movies[imdbid]
	if !ok {
		http.NotFound(w, r)
		return
	}
	enc := NewEncoder(w)
	err := enc.Encode(movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// bts, err := json.Marshal(movie)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
	// w.Write(bts)
}

func scenesHandler(w http.ResponseWriter, r *http.Request) {
	sceneid := r.URL.Path[len("/scenes/"):]
	scene, ok := ad.Scenes[sceneid]
	if !ok {
		http.NotFound(w, r)
		return
	}
	enc := NewEncoder(w)
	err := enc.Encode(scene)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func nearHandler(w http.ResponseWriter, r *http.Request) {
	lat, err := strconv.ParseFloat(r.FormValue("lat"), 64)
	lng, err := strconv.ParseFloat(r.FormValue("lng"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	loc := sfmovies.Location{"query", lat, lng}
	ds := make([]float64, 0)
	scs := make([]*sfmovies.Scene, 0)
	for _, scene := range ad.Scenes {
		ds = append(ds, loc.Distance(scene.Location))
		scs = append(scs, scene)
	}
	result := make([]*sfmovies.Scene, 0)
	// select the minimum distance, then remove them from arrays
	for i := 0; i < sfmovies.NearQuerySize; i++ {
		ix := mini(ds)
		if ix == -1 {
			break
		}
		result = append(result, scs[ix])
		ds[ix] = ds[len(ds)-1]
		scs[ix] = scs[len(scs)-1]
		ds = ds[:len(ds)-1]
		scs = scs[:len(scs)-1]
	}

	enc := NewEncoder(w)
	err = enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mini(a []float64) int {
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

func callbackHandler(fn Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		callback := r.FormValue("callback")
		if callback != "" {
			w.Header().Add("Content-Type", "application/javascript")
			w.Write([]byte(callback + "("))
		}
		fn(w, r)
		if callback != "" {
			w.Write([]byte(");"))
		}
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	// split q and intersect results
	result := trie.Get(q)
	if result == nil {
		result = new(TrieResults)
	}
	enc := NewEncoder(w)
	err := enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func completeHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("term")
	results := trie.GetFrom(q, sfmovies.AutoCompleteQuerySize)
	enc := NewEncoder(w)
	err := enc.Encode(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func statusHandler(w http.ResponseWriter, r *http.Request) {
	enc := NewEncoder(w)
	err := enc.Encode(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.String(), "/near") {
		nearHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.String(), "/search") {
		searchHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.String(), "/complete") {
		completeHandler(w, r)
		return
	}

	// bts, err := json.MarshalIndent(usage, "", "  ")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	_, err := io.WriteString(w, usage)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (e Encoder) Encode(v interface{}) error {
	bts, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = e.w.Write(bts)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// root handles near, search and complete queries as well as api description
	http.HandleFunc("/", callbackHandler(rootHandler))
	http.HandleFunc("/movies/", callbackHandler(moviesHandler))
	http.HandleFunc("/scenes/", callbackHandler(scenesHandler))
	http.HandleFunc("/status", callbackHandler(statusHandler))
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		fmt.Println(err)
	}
}
