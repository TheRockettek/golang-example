package main

import (
	"fmt"
	"sync/atomic"
	"github.com/go-martini/martini"
)

var (
	Accumulator = int64(0)
)

func GetCount() string {
	return fmt.Sprintf("%d", atomic.LoadInt64(&Accumulator))
}

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
