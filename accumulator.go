package main

import (
	"sync"
	"time"
	"fmt"
	"log"
)

const (
	// The total number of samples/values to retain on a per-hour basis.
	DefaultTotalSamples = 3600
)

// An accumulatorImpl manages storing data on an hourly basis.
type Accumulator interface {
	// Increments the value of the accumulator.
	Increment()

	// Returns how many times this was incremented in the last second.
	GetLastSecond() string

	// Returns how many times this was incremented in the last minute.
	GetLastMinute() string

	// Returns how many times this was incremented in the last hour.
	GetLastHour() string
}

type accumulatorImpl struct {
	sync.Mutex

	// A store of all the values to date.
	values []int64

	// The current value.
	acc int64
}

// Increment the current second's counter.
func (self *accumulatorImpl) Increment() {
	self.Lock()
	defer self.Unlock()

	self.acc += 1
}

// Gets the counter from the last second.
func (self *accumulatorImpl) GetLastSecond() string {
	self.Lock()
	defer self.Unlock()

	var str string

	if len(self.values) > 0 {
		str = fmt.Sprintf("%d", self.values[len(self.values)-1])
	} else {
		str = "0"
	}

	return str
}

// Get the total count of increment calls in the last minute.
func (self *accumulatorImpl) GetLastMinute() string {
	self.Lock()
	defer self.Unlock()

	// Get (up to) the position of the last 60 seconds worth of data.
	sliceStart := len(self.values) - 60

	if sliceStart < 0 {
		sliceStart = 0
	}

	slice := self.values[sliceStart:len(self.values)]

	// The final value that we want to accumulate to.
	var acc int64

	for _, val := range slice {
		acc += val
	}

	return fmt.Sprintf("%d", acc)
}

// Get the total count of increment calls in the last hour.
func (self *accumulatorImpl) GetLastHour() string {
	self.Lock()
	defer self.Unlock()

	// The final value that we want to accumulate to.
	var acc int64

	// Since we can guarantee that we will never have more than 1 hour's worth of
	// data in our values slice, we can just roll with suming everything.
	for _, val := range self.values {
		acc += val
	}

	return fmt.Sprintf("%d", acc)
}

// Actually starts the accumulatorImpl, aggregating data in to the internal data
// structures that lets us track what each count is.
func (self *accumulatorImpl) run() {
	ticks := time.Tick(1 * time.Second)

	for _ = range ticks {
		self.Lock()
		log.Printf("Counter value: %d", self.acc)
		self.values = append(self.values, self.acc)

		if len(self.values) > DefaultTotalSamples {
			self.values = self.values[1:DefaultTotalSamples]
		}

		self.acc = 0
		self.Unlock()
	}
}

// Instantiates and starts a new accumulator.
func NewAccumulator() Accumulator {
	acc := &accumulatorImpl{}

	// start running the accumulator. Yay!
	go acc.run()

	return acc
}
