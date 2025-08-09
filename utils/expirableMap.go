package utils

import (
	"sync"
	"time"
)

type Expirable[V any] struct {
	Value  V
	Expiry time.Time
}

type ExpirableMap[K comparable, T any] struct {
	Map  map[K]Expirable[T]
	TTL  time.Duration
	Lock sync.RWMutex
}

func NewExpirableMap[K comparable, V any](ttl time.Duration) *ExpirableMap[K, V] {
	return &ExpirableMap[K, V]{
		Map: map[K]Expirable[V]{},
		TTL: ttl,
	}
}

func (e *ExpirableMap[K, V]) Set(key K, value V) {
	e.Lock.Lock()
	e.Map[key] = Expirable[V]{
		Value:  value,
		Expiry: time.Now().Add(e.TTL),
	}
	e.Lock.Unlock()
}

func (e *ExpirableMap[K, V]) Get(key K) *V {
	e.Lock.RLock()
	value, ok := e.Map[key]
	e.Lock.RUnlock()
	if ok {
		if time.Now().Before(value.Expiry) {
			return &value.Value
		} else {
			e.Lock.Lock()
			delete(e.Map, key) // Clean up expired entries
			e.Lock.Unlock()
		}
	}
	return nil
}

func (e *ExpirableMap[K, V]) Len() int {
	e.Lock.RLock()
	length := len(e.Map)
	e.Lock.RUnlock()
	return length
}
