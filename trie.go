package sfmovies

import (
	"strings"
	"unicode"
)

type TrieResults struct {
	Movies []*Movie
	Scenes []*Scene
}

type TrieNode struct {
	next   []*TrieNode
	result *TrieResults // can be movie or
	letter rune
	prev   *TrieNode
}

func CreateTrie(data *APIData) *TrieNode {
	root := NewTrieNode()
	for _, movie := range data.Movies {
		root.AddMovieFields(movie, movie)
	}

	for _, scene := range data.Scenes {
		if movie, ok := data.Movies[scene.IMDBID]; ok {
			root.AddMovieFields(movie, scene)
		}
		root.AddMessyString(scene.Name, scene)
	}

	return root
}

func NewTrieNode() *TrieNode {
	tn := new(TrieNode)
	tn.next = make([]*TrieNode, 36) //alphabet + digits
	tn.result = new(TrieResults)
	return tn
}

func (t *TrieNode) AddMovieFields(movie *Movie, eltToAdd interface{}) {
	t.AddMessyString(movie.Title, eltToAdd)
	t.AddMessyString(movie.Year, eltToAdd)
	t.AddMessyString(movie.Writer, eltToAdd)
	t.AddMessyString(movie.Director, eltToAdd)
	t.AddMessyString(movie.Actors, eltToAdd)
}

func (t *TrieNode) AddMessyString(str string, elt interface{}) {
	str = CleanString(str)
	split := strings.Fields(str)
	for _, word := range split {
		t.Add(word, elt)
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

func (t *TrieNode) Add(str string, elt interface{}) {
	_, ismovie := elt.(*Movie)
	_, isscene := elt.(*Scene)
	if !ismovie && !isscene {
		return
	}
	str = CleanString(str)
	t.AddClean(str, elt)
}

func (t *TrieNode) AddClean(str string, elt interface{}) {
	if len(str) == 0 {
		if movie, ok := elt.(*Movie); ok {
			t.result.Movies = append(t.result.Movies, movie)
		}
		if scene, ok := elt.(*Scene); ok {
			t.result.Scenes = append(t.result.Scenes, scene)
		}
		return
	}
	if ix, ok := toIndex(str[0]); ok {
		if t.next[ix] == nil {
			t.next[ix] = NewTrieNode()
			t.next[ix].prev = t
			t.next[ix].letter = rune(str[0])
		}
		t.next[ix].AddClean(str[1:], elt)
	}
}

func (t *TrieNode) Get(str string) *TrieResults {
	return t.GetClean(CleanString(str))
}

func (t *TrieNode) GetClean(str string) *TrieResults {
	if len(str) == 0 {
		return t.result
	}
	if ix, ok := toIndex(str[0]); ok && t.next[ix] != nil {
		return t.next[ix].GetClean(str[1:])
	}
	return nil
}

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

func (t *TrieNode) GetFrom(str string, amount int) []string {
	str = CleanString(str)
	return t.GetFromClean(str, amount)
}

func (t *TrieNode) GetFromClean(str string, amount int) []string {
	r := make([]string, 0)
	if len(str) == 0 {
		t.BFS(&r, amount)
		return r
	}

	if ix, ok := toIndex(str[0]); ok {
		return t.next[ix].GetFromClean(str[1:], amount)
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
		if n.result.Movies != nil || n.result.Scenes != nil {
			*r = append(*r, n.String())
			amount--
			if amount == 0 {
				return
			}
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
