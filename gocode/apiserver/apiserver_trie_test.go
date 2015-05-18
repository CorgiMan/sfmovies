// Tests for apiserver_trie.go
package main

import (
	"testing"

	"github.com/CorgiMan/sfmovies/gocode"
)

func TestCleanString(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{"", ""},
		{"hello", "hello"},
		{"abc def ghi xyz 123 456 ab12", "abc def ghi xyz 123 456 ab12"},
		{"ABC DEF GhI xYz 123 456 AB12", "abc def ghi xyz 123 456 ab12"},
		{"me^&sf asf7 H3r 78fn#$^ewn4#^yu f d88+_\\:}{\"j", "mesf asf7 h3r 78fnewn4yu f d88j"},
	}
	for _, c := range cases {
		got := CleanString(c.in)
		if got != c.out {
			t.Errorf("CleanString(%v) == %v, want %v", c.in, got, c.out)
		}
	}
}

func TestToIndex(t *testing.T) {
	cases := []struct {
		in   byte
		out1 int
		out2 bool
	}{
		{'a', 0, true},
		{'f', 5, true},
		{'x', 23, true},
		{'y', 24, true},
		{'z', 25, true},
		{'0', 26, true},
		{'5', 31, true},
		{'9', 35, true},
		{' ', -1, false},
		{'*', -1, false},
		{'_', -1, false},
	}
	for _, c := range cases {
		got1, got2 := toIndex(c.in)
		if got1 != c.out1 {
			t.Errorf("toIndex(%v) == %v, want %v", c.in, got1, c.out2)
		}
		if got2 != c.out2 {
			t.Errorf("toIndex(%v) == %v, want %v", c.in, got2, c.out2)
		}
	}
}

func TestTrieNodeString(t *testing.T) {
	// Filled trie
	n := NewTrieNode()
	scene1, scene2 := &sfmovies.Scene{}, &sfmovies.Scene{}
	n.Add("abc", scene1)
	n.Add("abcdef", scene2)
	n.Add("abcABC", scene1)
	n.Add("1234", scene2)
	n.Add("abcABC", scene2)

	cases := []struct {
		in  *TrieNode
		out string
	}{
		{n, ""},
		{n.next[0], "a"},
		{n.next[27], "1"},
		{n.next[0].next[1], "ab"},
		{n.next[0].next[1].next[2], "abc"},
		{n.next[0].next[1].next[2].next[3].next[4].next[5], "abcdef"},
		{n.next[0].next[1].next[2].next[3].next[4].next[5], "abcdef"},
		{n.next[0].next[1].next[2].next[0].next[1].next[2], "abcabc"},
		{n.next[27].next[28].next[29].next[30], "1234"},
	}
	for _, c := range cases {
		got := c.in.String()
		if got != c.out {
			t.Errorf("n.String(%v) == %v, want %v", c.in, got, c.out)
		}
	}
}

func TestNewTrieNode(t *testing.T) {
	// Empty trie
	n := NewTrieNode()
	for _, node := range n.next {
		if node != nil {
			t.Errorf("Node has children")
		}
	}
}

func TestTrieNodeAdd(t *testing.T) {
	// Filled trie
	n := NewTrieNode()
	scene1, scene2 := &sfmovies.Scene{}, &sfmovies.Scene{}
	n.Add("abc", scene1)
	n.Add("abcdef", scene2)
	n.Add("abcABC", scene1)
	n.Add("1234", scene2)
	n.Add("abcabc", scene2)

	cases := []struct {
		s1, s2 *sfmovies.Scene
	}{
		{n.next[0].next[1].next[2].scenes[0], scene1},
		{n.next[0].next[1].next[2].next[3].next[4].next[5].scenes[0], scene2},
		{n.next[0].next[1].next[2].next[0].next[1].next[2].scenes[0], scene1},
		{n.next[27].next[28].next[29].next[30].scenes[0], scene2},
		{n.next[0].next[1].next[2].next[0].next[1].next[2].scenes[1], scene2},
	}

	for _, c := range cases {
		if c.s1 != c.s2 {
			t.Errorf("In trie: %p, want %p", c.s1, c.s2)
		}
	}
}

func TestTrieNodeAddMessyString(t *testing.T) {
	// Filled trie
	n := NewTrieNode()
	scene1, scene2 := &sfmovies.Scene{}, &sfmovies.Scene{}
	n.AddMessyString("abc def GhI x23", scene1)
	n.AddMessyString("abc defg ghi 321", scene2)

	cases := []struct {
		s1, s2 *sfmovies.Scene
	}{
		{n.next[0].next[1].next[2].scenes[0], scene1},
		{n.next[0].next[1].next[2].scenes[1], scene2},
		{n.next[3].next[4].next[5].scenes[0], scene1},
		{n.next[3].next[4].next[5].next[6].scenes[0], scene2},
		{n.next[6].next[7].next[8].scenes[0], scene1},
		{n.next[6].next[7].next[8].scenes[1], scene2},
		{n.next[29].next[28].next[27].scenes[0], scene2},
	}

	for _, c := range cases {
		if c.s1 != c.s2 {
			t.Errorf("In trie: %v, want %v", c.s1, c.s2)
		}
	}
}

func TestTrieNodeBFS(t *testing.T) {
	// Empty trie
	n := NewTrieNode()
	strs := make([]string, 0)
	n.BFS(&strs, 100)
	if len(strs) != 0 {
		t.Errorf("Got: %v, want %v", len(strs), 0)
	}

	// Filled trie
	n = NewTrieNode()
	scene1, scene2 := &sfmovies.Scene{}, &sfmovies.Scene{}
	n.Add("abc", scene1)
	n.Add("abcdef", scene2)
	n.Add("abcABC", scene1)
	n.Add("1234", scene2)
	n.Add("abcabc", scene2)

	strs = make([]string, 0)
	n.BFS(&strs, 100)

	M := make(map[string]bool)
	for _, str := range strs {
		M[str] = true
	}

	if len(M) != 4 {
		t.Errorf("Got: %v, want %v", len(M), 4)
	}

	strs2 := []string{"abc", "abcdef", "abcabc", "1234"}
	for _, s := range strs2 {
		if _, ok := M[s]; !ok {
			t.Errorf("Does not contain %v", s)
		}
	}
}
