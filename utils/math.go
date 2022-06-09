package utils

// Sum ...
func Sum[t any, retType ~float32 | ~float64 | intFamily](f func(x t) retType, s ...t) retType {
	var total retType = 0
	for i := range s {
		total += f(s[i])
	}
	return total
}
