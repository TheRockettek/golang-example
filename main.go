package main

import (
	"github.com/bradhe/golang-examples/accumulator"
	"github.com/go-martini/martini"
)

func main() {
	acc := accumulator.NewAccumulator()

	m := martini.Classic()
	m.Get("/v1/increment", acc.Increment)
	m.Get("/v1/counts/last-second", acc.GetLastSecond)
	m.Get("/v1/counts/last-minute", acc.GetLastMinute)
	m.Get("/v1/counts/last-hour", acc.GetLastHour)
	m.Run()
}
