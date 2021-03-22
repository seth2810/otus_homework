package hw04lrucache

import "sync/atomic"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	count int32
}

func (l *list) Len() int {
	return int(l.count)
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: l.front}

	if l.front != nil {
		l.front.Prev = item
	}

	if l.back == nil {
		l.back = item
	}

	l.front = item
	l.count++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v, Prev: l.back}

	if l.back != nil {
		l.back.Next = item
	}

	if l.front == nil {
		l.front = item
	}

	l.back = item

	atomic.AddInt32(&l.count, 1)

	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if i == l.front {
		l.front = i.Next
	}

	if i == l.back {
		l.back = i.Prev
	}

	i.Next = nil
	i.Prev = nil

	atomic.AddInt32(&l.count, -1)
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
