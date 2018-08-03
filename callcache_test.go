package callcache

import (
	"math/rand"
	"testing"
	"time"
)

func TestStartCalledImmediate(t *testing.T) {

	calls := 0

	c := CallCache{}
	c.Call = func() interface{} { calls++; return 0 }
	c.Interval = 100
	c.Start()

	close(c.stop)

	if calls != 1 {
		t.Fail()
	}
}

func TestStartCalledTimeout(t *testing.T) {
	callCount := 0
	calls := &callCount

	c := CallCache{}
	c.Call = func() interface{} { *calls++; return 0 }
	c.Interval = time.Duration(5)
	c.Start()

	time.Sleep(20 * time.Millisecond)
	close(c.stop)

	// It's either 4 or 5, depending on timings
	if callCount < 4 || callCount > 5 {
		t.Log(callCount)
		t.Fail()
	}
}

func TestDataCached(t *testing.T) {
	c := CallCache{}

	for i := 0; i < 5; i++ {
		num := rand.Int()
		c.Call = func() interface{} { return num }
		c.Interval = 1000
		c.Start()

		close(c.stop)

		if num != c.Fetch().(int) {
			t.Fail()
		}
	}
}

func TestStop(t *testing.T) {

	calls := 0

	c := CallCache{}
	c.Call = func() interface{} { calls++; return 0 }
	c.Interval = 100
	c.Start()
	c.Stop()

	time.Sleep(20 * time.Millisecond)
	if calls != 1 {
		t.Fail()
	}
}
