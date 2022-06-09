package utils

// ElemInSlice takes a slice and a single element of one comparable type, and returns if the element is in the slice.
func ElemInSlice[T comparable](slice []T, ele T) bool {
	for i := range slice {
		if slice[i] == ele {
			return true
		}
	}
	return false
}

// AppendNX takes a pointer to a slice and a single element of one comparable type,
// and appends the element to the slice if the element is not in the slice.
func AppendNX[T comparable](slice []T, eles ...T) []T {
	for _, ele := range eles {
		if !ElemInSlice(slice, ele) {
			slice = append(slice, ele)
		}
	}
	return slice
}

// RemoveRedundant returns a slice containing only unique elements from the origin slice.
func RemoveRedundant[T comparable](slice []T) []T {
	s := make(map[T]struct{})
	ret := []T{}
	for i := range slice {
		s[slice[i]] = struct{}{}
	}
	for t := range s {
		ret = append(ret, t)
	}
	return ret
}

// SliceFilter filters the paassed-in slice by the specified filtering function f, and returns a new slice containing all filted elements.
func SliceFilter[t any](s []t, f func(t) bool) []t {
	var ret []t = make([]t, 0, len(s))
	for _, e := range s {
		if f(e) {
			ret = append(ret, e)
		}
	}
	return ret
}

// SliceFilter filters the paassed-in slice by the specified filtering function f, and returns a new slice containing all filted elements.
// 当 copy 传入 true, 过滤一个新复制的切片, 并返回新切片;
// 当 copy 为空或者传入 false, 过滤原切片, 返回nil, 此时可以直接忽略返回值;
func SliceFilterPtr[t any](s *[]t, f func(t) bool, copy ...bool) []t {
	if copy != nil {
		if copy[0] {
			var ret []t = make([]t, 0, len(*s))
			for _, e := range *s {
				if f(e) {
					ret = append(ret, e)
				}
			}
			return ret
		}
	}
	var ret []t = make([]t, 0, len(*s))
	for i := range *s {
		if f((*s)[i]) {
			ret = append(ret, (*s)[i])
		}
	}
	*s = ret
	return nil
}

// MatchAny matches all passed-in `s` with the filter function f,
// and returns true at once if and only if any of them matches the filter function.
func MatchAny[t any](f func(t) bool, s ...t) bool {
	for i := 0; i < len(s); i++ {
		if f(s[i]) {
			return true
		}
	}
	return false
}

// MatchAll matches all passed-in `s` with the filter function f,
// and returns true at once if and only if all of them match the filter function.
func MatchAll[t any](f func(t) bool, s ...t) bool {
	for i := 0; i < len(s); i++ {
		if !f(s[i]) {
			return false
		}
	}
	return true
}
