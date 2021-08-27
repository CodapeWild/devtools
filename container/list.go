package container

import (
	"sync"
)

type List interface {
	Insert(pos *ListNode, v interface{}) *ListNode
	Delete(pos *ListNode) bool
	Find(target interface{}, comp func(lopr, ropr interface{}) bool) *ListNode
}

type ListNode struct {
	Previous *ListNode
	Next     *ListNode
	Data     interface{}
}

type LinkedList struct {
	sync.Mutex
	Head  *ListNode
	Tail  *ListNode
	total int
}

func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

// Insert will create new list node and insert into list.
// Head insertion will happen if pos is nil.
// Nil will return if pos is not present in the list.
func (this *LinkedList) Insert(pos *ListNode, v interface{}) *ListNode {
	this.Lock()
	defer this.Unlock()

	newnode := &ListNode{Data: v}
	// list is empty
	if this.Head == nil {
		this.Head = newnode
		this.Tail = newnode
	} else {
		if pos == nil {
			newnode.Next = this.Head
			this.Head.Previous = newnode
			this.Head = newnode
		} else {
			iter := this.Head
			for iter != nil {
				if iter == pos {
					newnode.Next = iter.Next
					iter.Next.Previous = newnode
					newnode.Previous = iter
					iter.Next = newnode
					if iter == this.Tail {
						this.Tail = newnode
					}
				}
				iter = iter.Next
			}
			if iter == nil {
				newnode = nil
			}
		}
	}

	return newnode
}

// Delete will remove the node represented by pos.
// False will return if pos is not present in the list.
func (this *LinkedList) Delete(pos *ListNode) bool {
	this.Lock()
	defer this.Unlock()

}

func (this *LinkedList) Find(target interface{}, comp func(lopr, ropr interface{}) bool) *ListNode {
	iter := this.Head
	for iter != nil {
		if comp(iter.Data, target) {
			return iter
		}
		iter = iter.Next
	}

	return nil
}

func (this *LinkedList) Count() int {
	return this.total
}

func (this *LinkedList) Empty() bool {
	return this.Head == this.Tail
}
