package container

import (
	"sync"
)

type ListNode struct {
	Data     interface{}
	Previous *ListNode
	Next     *ListNode
}

type ListIter func() *ListNode

type List interface {
	Insert(ListIter, *ListNode) int
	Delete(ListIter) bool
	Count() int
}

type LinkList struct {
	Head  *ListNode
	Tail  *ListNode
	count int
	sync.Mutex
}

func NewLinkList() *LinkList {
	return &LinkList{
		Head:  &ListNode{},
		Tail:  nil,
		count: 0,
		Mutex: sync.Mutex{},
	}
}

func (this *LinkList) Insert(iter ListIter, newNode *ListNode) int {
	this.Lock()
	defer this.Unlock()

	if this.Tail != nil {
		prev := this.Tail
		if iter != nil {
			prev = iter()
		}
		newNode.Previous = prev
		if prev.Next == nil {
			newNode.Next = nil
			prev.Next = newNode
			this.Tail = newNode
		} else {
			newNode.Next = prev.Next
			prev.Next.Previous = newNode
			prev.Next = newNode
		}

		this.count++
	} else {
		this.Head.Next = newNode
		newNode.Next = nil
		newNode.Previous = this.Head
		this.Tail = newNode

		this.count = 1
	}

	return this.count
}

func (this *LinkList) Delete(iter ListIter) bool {
	this.Lock()
	defer this.Unlock()

	if this.Tail != nil {
		var target = this.Tail
		if iter != nil {
			target = iter()
		}
		if target.Next == nil {
			target.Previous.Next = nil
			this.Tail = target.Previous
		} else {
			target.Previous.Next = target.Next
			target.Next.Previous = target.Previous
		}
	} else {
		return false
	}

	this.count--
	if this.count == 0 {
		this.Tail = nil
	}

	return true
}

func (this *LinkList) Count() int {
	this.Lock()
	defer this.Unlock()

	return this.count
}
