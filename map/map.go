package mmap

import (
	"encoding/json"
	"fmt"
	"sync"
)

// mMap embeds a lock and a generic map[K]V.
type mMap[K comparable, V any] struct {
	lock *sync.Mutex // to keep the Mutex-related fields and methods from being exported
	m    map[K]V     // underlying map that actually stores the K-V pairs
}

// keyval is a simple structs to hold the K-V pairs.
type keyval[K comparable, V any] struct {
	Key   K
	Value V
}

// NewMap returns a pointer to a newly instantiated Map,
// with the underlying map allocated with enough space to hold the specified number of elements.
func NewMap[K comparable, V any](size ...uint64) *mMap[K, V] {
	ret := &mMap[K, V]{lock: new(sync.Mutex)}
	switch len(size) {
	case 0:
		ret.m = make(map[K]V)
	default:
		ret.m = make(map[K]V, size[0])
	}
	return ret
}

// Set sets a K-V pair into the mMap, and returns the pointer to the mMap.
func (mm *mMap[K, V]) Set(key K, value V) *mMap[K, V] {
	mm.lock.Lock()
	mm.m[key] = value
	mm.lock.Unlock()
	return mm
}

// Get returns the value V from the mMap given a key K and a bool indicating whether the key K exists.
// If the key does exists, Get returns the value V, and true.
// Otherwise Get returns the specified default value defval, and false.
func (mm *mMap[K, V]) Get(key K, defval V) (V, bool) {
	mm.lock.Lock()
	ret, ok := mm.m[key]
	mm.lock.Unlock()
	if !ok {
		ret = defval
	}
	return ret, ok
}

// Del removes a K-V pair from the mMap.
// If there is no such key K, Del is a no-op.
func (mm *mMap[K, V]) Del(key K) *mMap[K, V] {
	mm.lock.Lock()
	delete(mm.m, key)
	mm.lock.Unlock()
	return mm
}

// Len returns the number of K-V pairs in the mMap.
func (mm *mMap[K, V]) Len() int {
	mm.lock.Lock()
	ret := len(mm.m)
	mm.lock.Unlock()
	return ret
}

// JsonFormatter prints the mMap in a pretty JSON format with indent.
// For debug purposes only.
func (mm *mMap[K, V]) JsonFormatter() {
	mm.lock.Lock()
	s, _ := json.MarshalIndent(mm.m, "", "  ")
	mm.lock.Unlock()
	fmt.Println(string(s))
}

// MapKeys returns the slice of all Keys in the map.
func (mm *mMap[K, V]) Keys() []K {
	ret := make([]K, 0, len(mm.m))
	for k := range mm.m {
		ret = append(ret, k)
	}
	mm.lock.Unlock()
	return ret
}

// MapValues returns the slice of all Values in the map.
func (mm *mMap[K, V]) Values(V) []V {
	mm.lock.Lock()
	ret := make([]V, 0, len(mm.m))
	for _, v := range mm.m {
		ret = append(ret, v)
	}
	mm.lock.Unlock()
	return ret
}

// Filter filters out the K-V pairs from the mMap by the specified filtering fuction f, and returns the original mMap.
// f should return false if the K-V pair is desired to be removed from the mMap.
func (mm *mMap[K, V]) Filter(f func(key K, value V) bool) *mMap[K, V] {
	mm.lock.Lock()
	for k, v := range mm.m {
		if !f(k, v) {
			delete(mm.m, k)
		}
	}
	mm.lock.Unlock()
	return mm
}

// FilterCopy returns a filtered mMap copied from the original mMap.
// f should return false if the K-V pair is desired to be removed from the mMap.
func (mm *mMap[K, V]) FilterCopy(f func(key K, value V) bool) *mMap[K, V] {
	mm.lock.Lock()
	ret := NewMap[K, V](uint64(len(mm.m)))
	for k, v := range mm.m {
		if f(k, v) {
			ret.m[k] = v
		}
	}
	mm.lock.Unlock()
	return ret
}

// ToSlice returns a slice of K-V pairs in the mMap.
// Notice that the order of the K-V pairs is not guaranteed.
func (mm *mMap[K, V]) ToSlice() []keyval[K, V] {
	mm.lock.Lock()
	ret := make([]keyval[K, V], 0, len(mm.m))
	for k, v := range mm.m {
		ret = append(ret, keyval[K, V]{Key: k, Value: v})
	}
	mm.lock.Unlock()
	return ret
}
