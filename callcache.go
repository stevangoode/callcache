package callcache

import (
	"sync"
	"time"
)

// CallCache describes and contains what is necessary to run the polling and cache the output
type CallCache struct {
	interval time.Duration      // How often to fire the function in ms
	call     func() interface{} // The function to be called
	lock     sync.RWMutex       // The lock for the cached data
	cache    interface{}        // The cached response from the func
	stop     chan interface{}   // The stop channel, which when closed, will stop the polling
}

// Fetch gets the cached function response
func (c *CallCache) Fetch() interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.cache
}

// Start begins the polling. Returns a channel which, when closed, will stop the polling
func (c *CallCache) Start() chan interface{} {
	c.stop = make(chan interface{})

	// Immediately start fetching the results of the method
	c.update()

	// This func will cause the cache to be updated on a schedule
	go func() {
		timer := time.NewTicker(c.interval * time.Millisecond)

		for {
			select {
			case <-timer.C:
				c.update()

			case <-c.stop:
				timer.Stop()
				return
			}
		}
	}()

	return c.stop
}

// Update causes the call to be done, and the result cached
func (c *CallCache) update() {
	response := c.call()

	c.lock.Lock()
	c.cache = response
	c.lock.Unlock()
}
