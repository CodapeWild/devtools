package structure

import (
	"log"
	"testing"
)

func showlist(list *LinkList, reverse bool) {
	cursor := list.Head.Next
	next := func(cursor *ListNode) *ListNode {
		return cursor.Next
	}
	if reverse {
		cursor = list.Tail
		next = func(cursor *ListNode) *ListNode {
			return cursor.Previous
		}
	}
	for ; cursor != nil; cursor = next(cursor) {
		log.Println(cursor.Data)
	}
	log.Println("#####################", list.Count())
}

func TestList(t *testing.T) {
	list := NewLinkList()
	for i := 0; i < 6; i++ {
		list.Insert(nil, &ListNode{Data: i})
	}

	showlist(list, false)
	log.Println(list.Delete(nil))
	showlist(list, true)
	list.Insert(nil, &ListNode{Data: 321})
	showlist(list, false)
	list.Insert(func() *ListNode { return list.Head.Next.Next }, &ListNode{Data: 111})
	showlist(list, true)
	log.Println("delete", list.Head.Next.Next.Next.Data)
	log.Println(list.Delete(func() *ListNode { return list.Head.Next.Next.Next }))
	showlist(list, false)
	for i := 0; i < 6; i++ {
		list.Delete(nil)
	}
	showlist(list, false)
}
