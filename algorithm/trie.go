package algorithm

import (
	"devtools/article"
	"log"
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
	i := article.CommonPrefixLen(this.String(), other.String())
	commonPrefix = TrieString(this.String()[:i])
	left = TrieString(this.String()[i:])

	return
}

type TrieNode struct {
	Prefix TrieWord
	Hit    int64
	Next   []*TrieNode
}

func NewTrieNode(word TrieWord) *TrieNode {
	return &TrieNode{Prefix: word}
}

type TrieRoot struct {
	*TrieNode
}

func NewTrieRoot(node *TrieNode) *TrieRoot {
	return &TrieRoot{node}
}

func (this *TrieRoot) Add(words ...TrieWord) int {
	added := 0
	for _, word := range words {
		if word.Len() == 0 {
			continue
		}
		// we supposed that root's prefix can't split if it's not empty then every word can't add into trie
		// if it do not has a common prefix with the root
		if this.Prefix.Len() != 0 && !this.Prefix.IsPrefixOf(word) {
			continue
		}

		if addIntoTrie(this.TrieNode, word) {
			added++
		}
	}

	return added
}

func (this *TrieRoot) Remove(words ...TrieWord) int {
	removed := 0
	for _, word := range words {
		if word.Len() == 0 {
			continue
		}

		if node := this.Find(word); node != nil {
			node.Hit--
			removed++
		}
	}

	return removed
}

func (this *TrieRoot) Find(word TrieWord) *TrieNode {
	// we supposed that root's prefix can't split if it's not empty then every word can't add into trie
	// if it has't a common prefix with the root
	if this.Prefix.Len() != 0 && !this.Prefix.IsPrefixOf(word) {
		return nil
	}

	var (
		cursor  = this.TrieNode
		_, left = word.SplitCommonPrefix(cursor.Prefix)
	)
	for left.Len() != 0 {
		found := false
		for _, node := range cursor.Next {
			if node.Prefix.IsPrefixOf(left) {
				_, left = left.SplitCommonPrefix(node.Prefix)
				cursor = node
				found = true
				break
			}
		}
		if !found {
			break
		}
	}

	if left.Len() == 0 && cursor.Hit > 0 {
		return cursor
	} else {
		return nil
	}
}

func addIntoTrie(cursor *TrieNode, word TrieWord) bool {
	if word.Equal(cursor.Prefix) {
		cursor.Hit++

		return true
	}

	var (
		cpx, left = word.SplitCommonPrefix(cursor.Prefix)
		added     = false
	)
	if cpx.Len() == cursor.Prefix.Len() {
		found := false
		for _, node := range cursor.Next {
			if cpx, _ = left.SplitCommonPrefix(node.Prefix); cpx.Len() != 0 {
				addIntoTrie(node, left)
				found = true
				break
			}
		}
		if !found {
			nn := NewTrieNode(left)
			nn.Hit++
			cursor.Next = append(cursor.Next, nn)
		}

		added = true
	} else if cpx.Len() < cursor.Prefix.Len() {
		_, cursorLeft := cursor.Prefix.SplitCommonPrefix(cpx)
		ln := NewTrieNode(cursorLeft)
		ln.Hit = cursor.Hit
		ln.Next = cursor.Next

		cursor.Prefix = cpx
		if left.Len() != 0 {
			rn := NewTrieNode(left)
			rn.Hit++
			if cursor.Hit > 0 {
				cursor.Hit--
			}
			cursor.Next = []*TrieNode{ln, rn}
		} else {
			cursor.Hit = 1
			cursor.Next = []*TrieNode{ln}
		}

		added = true
	}

	return added
}

func ShowTrie(cursor *TrieNode) {
	if cursor == nil {
		return
	}
	log.Printf("prefix:%s hit: %d", cursor.Prefix.String(), cursor.Hit)
	for _, node := range cursor.Next {
		ShowTrie(node)
	}
}
