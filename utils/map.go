package utils

// MapGet returns the stored value if the key exists, defaultValue otherwise.
func MapGet[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}

// MapFilter loops through the passed-in map by the filtering function f, deleting all key-value pairs that does not match the filter.
// 当 copy 传入 true, 过滤一个新复制的map, 并返回新map;
// 当 copy 为空或者传入 false, 过滤原map, 返回nil, 此时可以直接忽略返回值;
func MapFilter[K comparable, V any](m map[K]V, f func(key K, value V) bool, copy ...bool) map[K]V {
	if copy != nil {
		if copy[0] {
			ret := make(map[K]V)
			for k, v := range m {
				if f(k, v) {
					ret[k] = v
				}
			}
			return ret
		}
	}
	for k, v := range m {
		if !f(k, v) {
			delete(m, k)
		}
	}
	return nil
}

// MapKeys returns the slice of all Keys in the map.
func MapKeys[K comparable, V any](m map[K]V) []K {
	ret := make([]K, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

// MapValues returns the slice of all Values in the map.
func MapValues[K comparable, V any](m map[K]V) []V {
	ret := make([]V, 0, len(m))
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}
