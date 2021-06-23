package algorithm

import (
	"devtools/comerr"
	"unicode/utf8"
)

type TrieNode struct {
	value   rune
	count   int
	succeed []*TrieNode
}

func (this *TrieNode) FindSucceed(r rune) *TrieNode {
	for _, node := range this.succeed {
		if node.value == r {
			return node
		}
	}

	return nil
}

func (this *TrieNode) AppendSucceed(r rune) *TrieNode {
	newNode := &TrieNode{value: r}
	this.succeed = append(this.succeed, newNode)

	return newNode
}

type Trie struct {
	origin *TrieNode
}

func NewTrie() *Trie {
	return &Trie{origin: &TrieNode{}}
}

func (this *Trie) Add(s string) error {
	if len(s) == 0 {
		return nil
	}

	return add(this.origin, []byte(s))
}

func (this *Trie) Remove(s string) {
	if node, ok := this.Find(s); ok {
		node.count = 0
	}
}

// if find s in Trie return the node and true,
// if find s but not a previously added string return node and false,
// if not find s in trie return nil and false.
func (this *Trie) Find(s string) (*TrieNode, bool) {
	if len(s) == 0 {
		return nil, false
	}

	target := find(this.origin, []byte(s))

	if target == nil {
		return nil, false
	} else {
		return target, target.count != 0
	}
}

func add(cur *TrieNode, bs []byte) error {
	r, size := utf8.DecodeRune(bs)
	if r == utf8.RuneError && size == 1 {
		return comerr.ErrEncodeInvalid
	}
	bs = bs[size:]

	succeed := cur.FindSucceed(r)
	if succeed == nil {
		succeed = cur.AppendSucceed(r)
	}
	if len(bs) == 0 {
		succeed.count++

		return nil
	}

	return add(succeed, bs)
}

func find(cur *TrieNode, bs []byte) *TrieNode {
	r, size := utf8.DecodeRune(bs)
	if r == utf8.RuneError && size == 1 {
		return nil
	}
	bs = bs[size:]

	succeed := cur.FindSucceed(r)
	if succeed == nil {
		return nil
	} else {
		if len(bs) == 0 {
			return succeed
		} else {
			return find(succeed, bs)
		}
	}
}
