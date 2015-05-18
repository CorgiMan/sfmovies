// Implementation of a trie that stores scenes. The trie adds each scene a nodes representing
// Every word in the title, year, writer, director, actors and location name fields are stored in
// the trie. Each node that represent such a word is appended with a pointer to this scene.
// This means a single scene has multiple appearances in the trie!
// e.g. If we traverse the trie with the string "bill". We end up with a node that has
// scenes that are directed by Bill Guttentag, plus scenes that have the actor Bill Smitrovich
// plus scenes that are written by Bill Walsh. If there were movies with Bill in the title those
// would also be included.
package main

import (
	"strings"
	"unicode"

	"github.com/CorgiMan/sfmovies/gocode"
)

type TrieNode struct {
	next   []*TrieNode
	scenes []*sfmovies.Scene
	letter rune
	prev   *TrieNode
}

// given some APIData (a list of movies and scenes) this function construct a trie
// that can be queried by title, writers, directors, actors, location name. The
// trie returns a list of scenes containing the queried string.
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

// Cleans up and splits the string and the scene to every word of the split in the trie.
func (t *TrieNode) AddMessyString(str string, scene *sfmovies.Scene) {
	str = CleanString(str)
	split := strings.Fields(str)
	for _, word := range split {
		t.Add(word, scene)
	}
}

// Add a scene in the trie located at str.
func (t *TrieNode) Add(str string, scene *sfmovies.Scene) {
	str = CleanString(str)
	t.recursiveAdd(str, scene)
}

// Recursively traverse the try with str and append the scene to that node's scenes.
func (t *TrieNode) recursiveAdd(str string, scene *sfmovies.Scene) {
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
		t.next[ix].recursiveAdd(str[1:], scene)
	}
}

// The results of a search query
type SearchResults struct {
	Movies []*sfmovies.Movie
	Scenes []*sfmovies.Scene
}

// Traverse the trie with str and lists the scenes stored at this node.
// From these scenes a list of movies is composed which is also part of the result.
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

// Recursively traverses the trie with str and returns the scenes stored in this node.
func (t *TrieNode) recursiveGet(str string) []*sfmovies.Scene {
	if len(str) == 0 {
		return t.scenes
	}
	if ix, ok := toIndex(str[0]); ok && t.next[ix] != nil {
		return t.next[ix].recursiveGet(str[1:])
	}
	return nil
}

// Returns a list of words in the try starting with str.
func (t *TrieNode) GetFrom(str string, amount int) []string {
	str = CleanString(str)
	return t.recursiveGetFrom(str, amount)
}

// Recursively traverses the trie with str. The amount parameter is passed along.
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

// Does a breadth first search from the node. Stores full words in r and stops when the amount is reached.
func (t *TrieNode) BFS(r *[]string, amount int) {
	q := make([]*TrieNode, 0)
	q = append(q, t)
	for len(q) != 0 {
		// pop first element from the queue
		n := q[0]
		q = q[1:]

		// add all next nodes to the queue
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

		// add string to the result
		*r = append(*r, n.String())
		amount--
		if amount == 0 {
			return
		}
	}
}

// Constructs the string leading to this node by traversing the parent nodes
func (t *TrieNode) String() string {
	if t.prev == nil || t.prev.prev == nil {
		return string(t.letter)
	} else {
		return t.prev.String() + string(t.letter)
	}
}

// Remove all non a-z and non whitespace chars
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

// Returns index of [a-z, 0-9].
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
