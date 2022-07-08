package mlist

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

// mlist embeds a lock and a generic doubly linked list
type mlist[V any] struct {
	lock *sync.Mutex // to keep the Mutex-related fields and methods from being exported
	dll  *List[V]    // underlying slice that actually stores the values
}

// Init initializes or clears list l.
func (l *mlist[V]) InitList() *List[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.dll.Init()
}

// NewList returns a pointer to a newly instantiated List,
// embedding a mutex lock and
// a doubly linked list by Go Authors team tweaked to fit generics.
func NewList[V any]() *mlist[V] {
	ret := &mlist[V]{lock: new(sync.Mutex), dll: new(List[V]).Init()}
	return ret
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *mlist[V]) Len() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.dll.Len()
}

// Front returns the first element of list l or nil if the list is empty.
func (l *mlist[V]) Front() *Element[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.dll.Front()
}

// Back returns the last element of list l or nil if the list is empty.
func (l *mlist[V]) Back() *Element[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.dll.Back()
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *mlist[V]) PushFront(value V) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.dll.PushFront(value)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *mlist[V]) PushBack(value V) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.dll.PushBack(value)
}

// PushBackList adds a list of type V to the end of the list l.
func (l *mlist[V]) PushFrontList(list *List[V]) {
	l.lock.Lock()
	l.dll.PushFrontList(list)
	l.lock.Unlock()
}

// PushBackList adds a list of type V to the end of the list l.
func (l *mlist[V]) PushBackList(list *List[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.dll.PushBackList(list)
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *mlist[V]) Remove(ele *Element[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.dll.Remove(ele)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *mlist[V]) InsertBefore(v V, mark *Element[V]) *Element[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.dll.InsertBefore(v, mark)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *mlist[V]) InsertAfter(v V, mark *Element[V]) *Element[V] {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.dll.InsertAfter(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *mlist[V]) MoveToFront(e *Element[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.dll.MoveToFront(e)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *mlist[V]) MoveToBack(e *Element[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.dll.MoveToBack(e)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *mlist[V]) MoveBefore(e, mark *Element[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.dll.MoveBefore(e, mark)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *mlist[V]) MoveAfter(e, mark *Element[V]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.dll.MoveAfter(e, mark)
}

// Print prints the contents of the list.
func (l *mlist[V]) Print() (ret string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	var str strings.Builder
	str.WriteString("{}")
	str.WriteString(reflect.TypeOf(l.dll.root.Value).String())
	str.WriteString("[")
	s := ", "
	for e := l.dll.Front(); e != nil; e = e.Next() {
		str.WriteString(fmt.Sprint(e.Value))
		if e.Next() != nil {
			str.WriteString(s)
		}
	}
	str.WriteString("]")

	ret = str.String()
	fmt.Println(ret)
	return ret
}

// Sort is a trivial, unstable, and poorly-designed implementation that
// sorts the list according to specified function LessThan.
// LessThan(v1, v2 V) bool is expected to be a function that returns true iff v1 < v2.
// If LessThan is spcified as nil, Sort is a no-op.
func (l *mlist[V]) Sort(LessThan func(v1, v2 V) bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if LessThan != nil {
		tmp := make([]V, 0, l.dll.Len())
		for e := l.dll.Front(); e != nil; e = e.Next() {
			tmp = append(tmp, e.Value)
		}
		sort.Slice(tmp, func(i, j int) bool { return LessThan(tmp[i], tmp[j]) })
		l.dll.Init()
		for i := range tmp {
			l.dll.PushBack(tmp[i])
		}
	}
}
