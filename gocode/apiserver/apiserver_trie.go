package main

import (
	"strings"
	"unicode"

	"github.com/CorgiMan/sfmovies/gocode"
)

type SearchResults struct {
	Movies []*sfmovies.Movie
	Scenes []*sfmovies.Scene
}

type TrieNode struct {
	next   []*TrieNode
	scenes []*sfmovies.Scene
	letter rune
	prev   *TrieNode
}

func CreateTrie(data *sfmovies.APIData) *TrieNode {
	root := NewTrieNode()

	for _, scene := range data.Scenes {
		if movie, ok := data.Movies[scene.IMDBID]; ok {
			root.AddMessyString(movie.Title, scene)
			root.AddMessyString(movie.Year, scene)
			root.AddMessyString(movie.Writer, scene)
			root.AddMessyString(movie.Director, scene)
			root.AddMessyString(movie.Actors, scene)
		}
		root.AddMessyString(scene.Name, scene)
	}

	return root
}

func NewTrieNode() *TrieNode {
	tn := new(TrieNode)
	tn.next = make([]*TrieNode, 36) //alphabet + digits
	tn.scenes = make([]*sfmovies.Scene, 0)
	return tn
}

func (t *TrieNode) AddMessyString(str string, scene *sfmovies.Scene) {
	str = CleanString(str)
	split := strings.Fields(str)
	for _, word := range split {
		t.Add(word, scene)
	}
}

func (t *TrieNode) Add(str string, scene *sfmovies.Scene) {
	str = CleanString(str)
	t.AddClean(str, scene)
}

func (t *TrieNode) AddClean(str string, scene *sfmovies.Scene) {
	if len(str) == 0 {
		t.scenes = append(t.scenes, scene)
		return
	}
	if ix, ok := toIndex(str[0]); ok {
		if t.next[ix] == nil {
			t.next[ix] = NewTrieNode()
			t.next[ix].prev = t
			t.next[ix].letter = rune(str[0])
		}
		t.next[ix].AddClean(str[1:], scene)
	}
}

func (t *TrieNode) Get(ad *sfmovies.APIData, str string) *SearchResults {
	scenes := t.recursiveGet(CleanString(str))
	if scenes == nil {
		return nil
	}

	r := new(SearchResults)
	r.Scenes = scenes

	// remove dups from movies
	M := make(map[string]bool)
	for _, scene := range scenes {
		M[scene.IMDBID] = true
	}
	r.Movies = make([]*sfmovies.Movie, 0)
	for imdbid := range M {
		if movie, ok := ad.Movies[imdbid]; ok {
			r.Movies = append(r.Movies, movie)
		}
	}

	return r
}

func (t *TrieNode) recursiveGet(str string) []*sfmovies.Scene {
	if len(str) == 0 {
		return t.scenes
	}
	if ix, ok := toIndex(str[0]); ok && t.next[ix] != nil {
		return t.next[ix].recursiveGet(str[1:])
	}
	return nil
}

func (t *TrieNode) GetFrom(str string, amount int) []string {
	str = CleanString(str)
	return t.recursiveGetFrom(str, amount)
}

func (t *TrieNode) recursiveGetFrom(str string, amount int) []string {
	r := make([]string, 0)
	if len(str) == 0 {
		t.BFS(&r, amount)
		return r
	}

	if ix, ok := toIndex(str[0]); ok && t.next[ix] != nil {
		return t.next[ix].recursiveGetFrom(str[1:], amount)
	}
	return r
}

func (t *TrieNode) BFS(r *[]string, amount int) {
	q := make([]*TrieNode, 0)
	q = append(q, t)
	for len(q) != 0 {
		n := q[0]
		q = q[1:]
		for _, m := range n.next {
			if m != nil {
				q = append(q, m)
			}
		}

		// don't add results smaller than 3 chars
		if n.prev == nil || n.prev.prev == nil || n.prev.prev.prev == nil {
			continue
		}
		if len(n.scenes) == 0 {
			continue
		}

		*r = append(*r, n.String())
		amount--
		if amount == 0 {
			return
		}
	}
}

func (t *TrieNode) String() string {
	if t.prev == nil || t.prev.prev == nil {
		return string(t.letter)
	} else {
		return t.prev.String() + string(t.letter)
	}
}

// remove all non a-z and non whitespace chars
func CleanString(str string) string {
	str = strings.ToLower(str)
	clean := make([]rune, 0)
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsSpace(r) || unicode.IsDigit(r) {
			clean = append(clean, r)
		}
	}
	return string(clean)
}

// returns index of [a-z, 0-9]
func toIndex(c byte) (int, bool) {
	ix := int(c - 'a')
	if ix >= 0 && ix < 26 {
		return ix, true
	}
	ix = int(c - '0')
	if ix >= 0 && ix < 10 {
		return ix + 26, true
	}
	return -1, false
}
