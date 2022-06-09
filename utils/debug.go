package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
)

// JsonFormatter should take in only a struct with exported fields, or a map..
// For debug purposes only.
func JsonFormatter[Any any](a ...Any) {
	for i := range a {
		s, _ := json.MarshalIndent(a[i], "", "  ")
		fmt.Println(string(s))
	}
}

// Destructor should be able to take in anything and print it out beautifully...
// But a struct is prefered.
func Destructor[Any any](a ...Any) {
	for i := range a {
		fmt.Printf("%#v\n", a[i])
	}
}

// PanicTrace tries to parse the stack frame and extract the actual useful info
// after the string literal "panic.go" in the stack frame.
// kb is the size of info returned in KB; usually 2 is enough.
func PanicTrace(kb int) string {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10)
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return string(stack)
}
