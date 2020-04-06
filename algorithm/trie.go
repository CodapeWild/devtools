package algorithm

import (
	"devtools/article"
	"devtools/comerr"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type TrieWord interface {
	Len() int
	String() string
	Equal(TrieWord) bool
	IsPrefixOf(TrieWord) bool
	SplitCommonPrefix(TrieWord) (commonPrefix, left TrieWord)
}

type TrieString string

func (this TrieString) Len() int {
	return len(this)
}

func (this TrieString) String() string {
	return string(this)
}

func (this TrieString) Equal(other TrieWord) bool {
	return this.String() == other.String()
}

func (this TrieString) IsPrefixOf(other TrieWord) bool {
	return this.Len() <= other.Len() && this.String() == other.String()[:this.Len()]
}

func (this TrieString) SplitCommonPrefix(other TrieWord) (commonPrefix, left TrieWord) {
	commonPrefix = TrieString(article.CommonPrefix(this.String(), other.String()))
	left = TrieString(strings.TrimPrefix(this.String(), commonPrefix.String()))

	return
}

type TrieNode struct {
	Prefix TrieWord
	Hit    int64
	Next   []*TrieNode
}

func NewTrieNode(word TrieWord, hit int64, next []*TrieNode) *TrieNode {
	return &TrieNode{
		Prefix: word,
		Hit:    hit,
		Next:   next,
	}
}

// root's prefix do not split, which means any word has no common prefix word with root can not add into trie
// unless root's prefix is empty.
func (this *TrieNode) Add(words ...TrieWord) int {
	var added = 0
	for _, word := range words {
		if word.Len() == 0 {
			continue
		}

		if this.Prefix.Len() != 0 && !this.Prefix.IsPrefixOf(word) {
			continue
		}

		_, left := word.SplitCommonPrefix(this.Prefix)
		if left.Len() == 0 {
			this.Hit++
			continue
		}
		addIntoTrie(this, left)
		added++
	}

	return added
}

func (this *TrieNode) Remove(words ...TrieWord) int {
	var removed = 0
	for _, word := range words {
		if word.Len() == 0 {
			continue
		}

		if node, ok := this.Find(word); ok {
			node.Hit = 0
			removed++
		}
	}

	return removed
}

func (this *TrieNode) Find(word TrieWord) (*TrieNode, bool) {
	if this.Prefix.Len() != 0 && !this.Prefix.IsPrefixOf(word) {
		return nil, false
	}

	var (
		_, left = word.SplitCommonPrefix(this.Prefix)
		found   = false
	)
	for left.Len() != 0 {
		found = false
		for _, node := range this.Next {
			if node.Prefix.IsPrefixOf(left) {
				_, left = left.SplitCommonPrefix(node.Prefix)
				this = node
				found = true
				break
			}
		}
		if !found {
			break
		}
	}

	if left.Len() == 0 && this.Hit > 0 {
		return this, true
	} else {
		return nil, false
	}
}

func addIntoTrie(cursor *TrieNode, word TrieWord) {
	var found = false
	for _, node := range cursor.Next {
		if word.Equal(node.Prefix) {
			node.Hit++

			return
		}

		compfx, left := word.SplitCommonPrefix(node.Prefix)
		if compfx.Len() == 0 {
			continue
		}

		if node.Prefix.Len() == compfx.Len() {
			cursor = node
			word = left
			found = true
			break
		} else {
			_, nl := node.Prefix.SplitCommonPrefix(compfx)
			node.Prefix = compfx
			n := NewTrieNode(nl, node.Hit, node.Next)
			node.Hit = 0
			node.Next = []*TrieNode{n}
			if left.Len() == 0 {
				node.Hit++
			} else {
				node.Next = append(node.Next, NewTrieNode(left, 1, nil))
			}

			return
		}
	}

	if !found {
		cursor.Next = append(cursor.Next, NewTrieNode(word, 1, nil))
	} else {
		addIntoTrie(cursor, word)
	}
}

func ShowTrie(cursor *TrieNode) {
	log.Printf("prefix:%s hit:%d next:%d\n", cursor.Prefix.String(), cursor.Hit, len(cursor.Next))
	for _, node := range cursor.Next {
		ShowTrie(node)
	}
}

type TrieNodeJson struct {
	Word string `json:"word"`
	Hit  int64  `json:"hit"`
}

type TrieJson []TrieNodeJson

func TrieToJson(root *TrieNode, w io.Writer) error {
	var tj TrieJson
	traverse(root, "", &tj)

	buf, err := json.Marshal(tj)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)

	return err
}

func traverse(cursor *TrieNode, word string, out *TrieJson) {
	word += cursor.Prefix.String()
	if cursor.Hit > 0 {
		*out = append(*out, TrieNodeJson{Word: word, Hit: cursor.Hit})
	}
	for _, node := range cursor.Next {
		traverse(node, word, out)
	}
}

func TrieFromJson(r io.Reader) (root *TrieNode, err error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var tj = TrieJson{}
	if err = json.Unmarshal(buf, &tj); err != nil {
		return nil, err
	}
	if len(tj) == 0 {
		return nil, comerr.EmptyData
	}

	root = NewTrieNode(TrieString(tj[0].Word), tj[0].Hit, nil)
	for _, node := range tj {
		root.Add(TrieString(node.Word))
		if n, ok := root.Find(TrieString(node.Word)); !ok {
			return nil, comerr.ProcessFailed
		} else {
			n.Hit = node.Hit
		}
	}

	return
}
