package main

import (
	"fmt"
	"sync/atomic"
	"github.com/go-martini/martini"
)

var (
	Accumulator = int64(0)
)

// This function handles getting the current value of the accumulator. It will
// return the current value (atomically) formatted as a string.
func GetCount() string {
	return fmt.Sprintf("%d", atomic.LoadInt64(&Accumulator))
}

// This function handles incrementing the counter atomically.
func GetIncrement() string {
	atomic.AddInt64(&Accumulator, 1)

	// Empty body is OK. Status 200 indicates that it all worked.
	return ""
}

func main() {
	m := martini.Classic()
	m.Get("/v1/increment", GetIncrement)
	m.Get("/v1/count", GetCount)
	m.Run()
}
