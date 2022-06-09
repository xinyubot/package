package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
)

// literal is a type constraint that contains string, []byte, and their alias and variants.
type literal interface {
	~string | ~[]byte
}

// LiteralToMD5 returns a MD5 hash in string form of the passed-in `s`.
func LiteralToMD5[T literal](s T) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// LiteralToMD5 returns a MD5 hash in []byte] form of the passed-in `s`.
func LiteralToMD5Bytes[T literal](s T) []byte {
	r := md5.Sum([]byte(s))
	return r[:]
}

// StringConcat returns the concatenation of `base` and `strs` strings seperated by `sep`
func StringConcat(sep string, strs ...string) string {
	return strings.Join(strs, sep)
}

// intFamily contains all int types and their alias and variants.
type intFamily interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// TimeConcatTwoDigits ... 如果时间单位的值是个位数，则拼接一个0
func TimeConcatTwoDigits[T intFamily](v T) (ret string, err error) {
	switch {
	case v > 0 && v < 10:
		return fmt.Sprintf("0%d", v), nil
	case v < 100:
		return fmt.Sprintf("%d", v), nil
	default:
		return "", errors.New("wrong digits")
	}
}
