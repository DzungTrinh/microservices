package main

import (
	"github.com/sony/gobreaker"
	"sync"
	"time"
)

var breakers = struct {
	sync.RWMutex
	instances map[string]*gobreaker.CircuitBreaker
}{instances: make(map[string]*gobreaker.CircuitBreaker)}

func getBreaker(host string) *gobreaker.CircuitBreaker {
	breakers.RLock()
	cb, exists := breakers.instances[host]
	breakers.RUnlock()
	if exists {
		return cb
	}

	settings := gobreaker.Settings{
		Name:        host,
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	}

	cb = gobreaker.NewCircuitBreaker(settings)
	breakers.Lock()
	breakers.instances[host] = cb
	breakers.Unlock()
	return cb
}
