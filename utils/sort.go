package utils

import (
	"sort"
)

// SortSlice ... 不稳定排序
// 当 copy 传入 true, 排序一个新复制的切片, 并返回新切片;
// 当 copy 为空或者传入 false, 排序原切片, 返回nil, 此时可以直接忽略返回值;
func SortSlice[t any](s []t, f func(x, y t) bool, copy ...bool) []t {
	if copy != nil {
		if copy[0] {
			ret := make([]t, 0, len(s))
			for i := range s {
				ret = append(ret, s[i])
			}
			sort.Slice(ret, func(i, j int) bool { return f(ret[i], ret[j]) })
			return ret
		}
	}
	sort.Slice(s, func(i, j int) bool { return f(s[i], s[j]) })
	return nil
}

// SortSliceStable ... 稳定排序
// 当 copy 传入 true, 排序一个新复制的切片, 并返回新切片;
// 当 copy 为空或者传入 false, 排序原切片, 返回nil, 此时可以直接忽略返回值;
func SortSliceStable[t any](s []t, f func(x, y t) bool, copy ...bool) []t {
	if copy != nil {
		if copy[0] {
			ret := make([]t, 0, len(s))
			for i := range s {
				ret = append(ret, s[i])
			}
			sort.SliceStable(ret, func(i, j int) bool { return f(ret[i], ret[j]) })
			return ret
		}
	}
	sort.SliceStable(s, func(i, j int) bool { return f(s[i], s[j]) })
	return nil
}
