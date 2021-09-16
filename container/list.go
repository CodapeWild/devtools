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

// Append add new node after tail and return the new node pointer.
func (this *LinkedList) Append(v interface{}) *ListNode {
	this.Lock()
	defer this.Unlock()

	newTail := &ListNode{Data: v}
	if this.Head == nil {
		this.Head = newTail
		this.Tail = newTail
	} else {
		newTail.Previous = this.Tail
		this.Tail.Next = newTail
		this.Tail = newTail
	}
	this.total++

	return newTail
}

// Prepend add new node before head and return new node pointer.
func (this *LinkedList) Prepend(v interface{}) *ListNode {
	this.Lock()
	defer this.Unlock()

	newHead := &ListNode{Data: v}
	if this.Head == nil {
		this.Head = newHead
		this.Tail = newHead
	} else {
		newHead.Next = this.Head
		this.Head.Previous = newHead
		this.Head = newHead
	}
	this.total++

	return newHead
}

// Insert add new node after pos and return new node pointer,
// nil pointer will return if pos is not present in list.
func (this *LinkedList) Insert(pos *ListNode, v interface{}) *ListNode {
	if pos == nil {
		return nil
	}

	this.Lock()
	defer this.Unlock()

	iter := this.Head
	for iter != pos {
		iter = iter.Next
	}
	var newNode *ListNode
	if iter == pos {
		newNode = &ListNode{Data: v}
		if iter.Next != nil {
			newNode.Next = iter.Next
			newNode.Previous = iter
			iter.Next.Previous = newNode
			iter.Next = newNode
		} else {
			iter.Next = newNode
			newNode.Previous = iter
			this.Tail = newNode
		}
		this.total++
	}

	return newNode
}

// DeleteWithPosition will remove node referenced by pos and return true.
func (this *LinkedList) DeleteWithPosition(pos *ListNode) bool {
	if pos == nil {
		return false
	}

	this.Lock()
	defer this.Unlock()

	if pos == this.Head {
		this.Head.Next.Previous = nil
		this.Head = this.Head.Next
		this.total--

		return true
	} else if pos == this.Tail {
		this.Tail.Previous.Next = nil
		this.Tail = this.Tail.Previous
		this.total--

		return true
	}

	iter := this.Head
	for iter != pos {
		iter = iter.Next
	}
	if iter == pos {
		iter.Previous.Next = iter.Next
		iter.Next.Previous = iter.Previous
		this.total--

		return true
	}

	return false
}

// DeleteWithValue
func (this *LinkedList) DeleteWithValue(v interface{}, comp func(lopr, ropr interface{}) bool, count int) int {
	this.Lock()
	defer this.Unlock()

	iter := this.Head
	deleted := 0
	for iter != nil {
		if comp(iter.Data, v) {
			if count > 0 && deleted == count {
				break
			}

		}
	}
}

// Find return the first element find in the list.
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

// FindAll return all the elements find in the list.
func (this *LinkedList) FindAll(target interface{}, comp func(lopr, ropr interface{}) bool) []*ListNode {

}

func (this *LinkedList) Count() int {
	return this.total
}

func (this *LinkedList) Empty() bool {
	return this.Head == this.Tail
}

func (this *LinkedList) add(pos *ListNode, v interface{}) *ListNode {
	newNode := &ListNode{Data: v}
	if pos.Next == nil { // append
		newNode.Previous = pos
		pos.Next = newNode
		this.Tail = newNode
	} else if pos.Previous == nil { // prepend
		newNode.Next = pos
		pos.Previous = newNode
		this.Head = newNode
	} else { // insert
		newNode.Next = pos.Next
		newNode.Previous = pos
		pos.Next.Previous = newNode
		pos.Next = newNode
	}
	this.total++

	return newNode
}

func (this *LinkedList) remove(pos *ListNode) {
	if pos == this.Head { // remove head
		pos.Next.Previous = nil
		pos.Next = nil
		this.Head = pos
	} else if pos == this.Tail { // remove tail
		pos.Previous.Next = nil
		pos.Previous = nil
		this.Tail = pos
	} else { // remove
		pos.Previous.Next = pos.Next
		pos.Next.Previous = pos.Previous
		pos.Next = nil
		pos.Previous = nil
	}
	this.total--
}
