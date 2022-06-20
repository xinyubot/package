package mlist

import (
	"sync"
)

type mlist[V any] struct {
	lock *sync.Mutex // to keep the Mutex-related fields and methods from being exported
	m    *List[V]    // underlying slice that actually stores the values
}

// Init initializes or clears list l.
func (l *mlist[V]) InitList() *List[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.m.Init()
}

// NewSlice returns a pointer to a newly instantiated Slice,
// with the underlying slice allocated with enough space to hold the specified number of elements.
func NewList[V any]() *mlist[V] {
	ret := &mlist[V]{lock: new(sync.Mutex), m: new(List[V]).Init()}
	return ret
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *mlist[V]) Len() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.m.Len()
}

// Front returns the first element of list l or nil if the list is empty.
func (l *mlist[V]) Front() *Element[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.m.Front()
}

// Back returns the last element of list l or nil if the list is empty.
func (l *mlist[V]) Back() *Element[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.m.Back()
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *mlist[V]) PushFront(value V) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.m.PushFront(value)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *mlist[V]) PushBack(value V) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.m.PushBack(value)
}

// PushBackList adds a list of type V to the end of the list l.
func (l *mlist[V]) PushFrontList(list *List[V]) {
	l.lock.Lock()
	l.m.PushFrontList(list)
	l.lock.Unlock()
}

// PushBackList adds a list of type V to the end of the list l.
func (l *mlist[V]) PushBackList(list *List[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.m.PushBackList(list)
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *mlist[V]) Remove(ele *Element[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.m.Remove(ele)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *mlist[V]) InsertBefore(v V, mark *Element[V]) *Element[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.m.InsertBefore(v, mark)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *mlist[V]) InsertAfter(v V, mark *Element[V]) *Element[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.m.InsertAfter(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *mlist[V]) MoveToFront(e *Element[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.m.MoveToFront(e)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *mlist[V]) MoveToBack(e *Element[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.m.MoveToBack(e)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *mlist[V]) MoveBefore(e, mark *Element[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.m.MoveBefore(e, mark)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *mlist[V]) MoveAfter(e, mark *Element[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.m.MoveAfter(e, mark)
}
