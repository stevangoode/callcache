// Package callcache will make a method call and cache the response. This can be useful
// for heavy operations, for which the output does not change very often, for example
// when making an API call to get a list of entities.
//
// Simple Example:
//
// cache := callcache.CallCache{
//   Interval: time.Duration(2 * time.Second)
//   Call: api.GetMyData
// }
// cache.Start()
//
// With Parameters:
//
// cache := callcache.CallCache{
//   Interval: time.Duration(2 * time.Second)
//   Call: func () []api.MyDataType {return api.GetMyData(aParameter)}
// }
// cache.Start()
//
// Fetching Data:
//
// myData := cache.Fetch()
//
// Stop Polling:
//
// cache.Stop()
//
// Change Duration:
//
// cache.Stop()
// cache.Interval = time.Duration(10 * time.Second)
// cache.Start()
package callcache

import (
	"sync"
	"time"
)

// CallCache describes and contains what is necessary to run the polling and cache the output
type CallCache struct {
	Interval time.Duration      // How often to fire the function in ms
	Call     func() interface{} // The function to be called
	lock     sync.RWMutex       // The lock for the cached data
	cache    interface{}        // The cached response from the func
	stop     chan interface{}   // The stop channel, which when closed, will stop the polling
}

// Fetch gets the cached function response
func (cache *CallCache) Fetch() interface{} {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	return cache.cache
}

// Start begins the polling. Returns a channel which, when closed, will stop the polling
func (cache *CallCache) Start() chan interface{} {
	cache.stop = make(chan interface{})

	// Immediately start fetching the results of the method
	cache.update()

	// This func will cause the cache to be updated on a schedule
	go func() {
		timer := time.NewTicker(cache.Interval * time.Millisecond)

		for {
			select {
			case <-timer.C:
				cache.update()

			case <-cache.stop:
				timer.Stop()
				return
			}
		}
	}()

	return cache.stop
}

// Stop will stop the polling and updating until Start is called again
func (cache *CallCache) Stop() {
	close(cache.stop)
	cache.stop = nil
}

// Update causes the call to be done, and the result cached
func (cache *CallCache) update() {
	response := cache.Call()

	cache.lock.Lock()
	cache.cache = response
	cache.lock.Unlock()
}
