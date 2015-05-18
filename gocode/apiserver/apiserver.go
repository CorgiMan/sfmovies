// Web server that handles requests for the API service
// API handles auto-complete, search and location based queries.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CorgiMan/sfmovies/gocode"
)

// listens on port 80 by default
var port = flag.String("port", "80", "port that program listens on")

// Data used by API server
var (
	appData *sfmovies.APIData
	trie    *TrieNode
	status  Status
)

// The status of an API server. This is used when "/status" is requested
type Status struct {
	APIVersion   string
	RunningSince time.Time
	DataVersion  time.Time
}

// API servers errors are served by encoding this struct to JSON
type Error struct {
	Error string
}

// Download the latest APIData from MongoDB, calculate the search trie and determine status of the server.
func init() {
	flag.Parse()

	var err error
	appData, err = sfmovies.GetLatestAPIData()
	if err != nil {
		log.Fatal(err)
	}
	if appData == nil {
		log.Fatal(errors.New("No API data received from mongodb"))
	}
	status = Status{}
	status.APIVersion = sfmovies.APIVersion
	status.RunningSince = time.Now()
	status.DataVersion = appData.Time

	trie = CreateTrie(appData)
}

// Sets up the webserver on a port specified by the --port flag
func main() {
	// root handles near, search and complete queries as well as API description
	http.HandleFunc("/", jsonpHandler(rootHandler))
	http.HandleFunc("/movies/", jsonpHandler(moviesHandler))
	http.HandleFunc("/status", jsonpHandler(statusHandler))
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Handles near, search, complete and root (usage) queries
func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/near":
		nearHandler(w, r)
	case "/search":
		searchHandler(w, r)
	case "/complete":
		completeHandler(w, r)
	case "/":
		_, err := io.WriteString(w, sfmovies.Usage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Encodes an object into JSON and writes to the response writer
func writeResult(w http.ResponseWriter, v interface{}) {
	bts, err1 := json.MarshalIndent(v, "", "  ")
	_, err2 := w.Write(bts)
	if err1 != nil || err2 != nil {
		http.Error(w, "failed to marshal and write json", http.StatusInternalServerError)
	}
}

// Writes the status of the API server
func statusHandler(w http.ResponseWriter, r *http.Request) {
	writeResult(w, status)
}

// Handles queries for a specific IMDB movie ID.
func moviesHandler(w http.ResponseWriter, r *http.Request) {
	imdbid := r.URL.Path[len("/movies/"):]
	if movie, ok := appData.Movies[imdbid]; ok {
		writeResult(w, movie)
	} else {
		writeResult(w, Error{"Recource not found"})
	}
}

// Handles auto-complete queries.
func completeHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("term")
	result := trie.GetFrom(q, sfmovies.AutoCompleteQuerySize)
	writeResult(w, result)
}

// Handles queries that search for a complete word.
func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	if result := trie.Get(appData, q); result != nil {
		writeResult(w, result)
	} else {
		writeResult(w, Error{"Recource not found"})
	}
}

// Handles near queries. Returns the closest NearQuerySize points-of-interest.
func nearHandler(w http.ResponseWriter, r *http.Request) {
	lat, err1 := strconv.ParseFloat(r.FormValue("lat"), 64)
	lng, err2 := strconv.ParseFloat(r.FormValue("lng"), 64)
	if err1 != nil || err2 != nil {
		http.Error(w, "failed to parse lat and lng parameters", http.StatusInternalServerError)
		return
	}
	loc := sfmovies.Location{"", lat, lng}
	ds := make([]float64, 0)
	scs := make([]*sfmovies.Scene, 0)
	for _, scene := range appData.Scenes {
		ds = append(ds, loc.Distance(scene.Location))
		scs = append(scs, scene)
	}

	// select closest element NearQuerySize times
	result := make([]*sfmovies.Scene, 0)
	for i := 0; i < sfmovies.NearQuerySize; i++ {
		// select nearest element
		ix := minix(ds)
		if ix == -1 {
			break
		}
		result = append(result, scs[ix])

		// remove it from the array so that the next closest element will be found next time
		ds[ix] = ds[len(ds)-1]
		scs[ix] = scs[len(scs)-1]
		ds = ds[:len(ds)-1]
		scs = scs[:len(scs)-1]
	}

	writeResult(w, result)
}

// Returns the index of the smallest element in a.
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

// Wraps around all other handlers and adds JSONP padding only if the callback parameter is set.
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
