package algorithm

import (
	"devtools/article"
)

type TrieWord interface {
	Len() int
	String() string
	Split(TrieWord) (commonPrefix, left TrieWord)
}

type TrieString string

func (this TrieString) Len() int {
	return len(this)
}

func (this TrieString) String() string {
	return string(this)
}

func (this TrieString) Split(other TrieWord) (commonPrefix, left TrieWord) {
	i := article.CommonPrefixLen(this.String(), other.String())
	commonPrefix = TrieString(this.String()[:i])
	left = TrieString(this.String()[i:])

	return
}

type TrieNode struct {
	Prefix TrieWord
	Count  int64
	Next   []*TrieNode
}

func NewTrieNode(word TrieWord) *TrieNode {
	return &TrieNode{Prefix: word}
}

/*
			a
	bc		def
*/
func AddIntoTrie(root *TrieNode, words ...TrieWord) {
	if root == nil {
		return
	}

	for _, w := range words {
		if root.Prefix.Len() != 0 {
			if cpx, _ := w.Split(root.Prefix); cpx.Len() == 0 {
				continue
			}
		}

	}
}

func RemoveFromTrie(root *TrieNode, words ...TrieWord) {
	if root == nil {
		return
	}
}

func ExistsInTrie(root *TrieNode, word TrieWord) bool {
	return false
}

func ShowTrie(root *TrieNode) {

}
