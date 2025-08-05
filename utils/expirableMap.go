package utils

import "time"

type Expirable[V any] struct {
	Value  V
	Expiry time.Time
}

type ExpirableMap[K comparable, T any] struct {
	Map map[K]Expirable[T]
	TTL time.Duration
}

func NewExpirableMap[K comparable, V any](ttl time.Duration) *ExpirableMap[K, V] {
	return &ExpirableMap[K, V]{
		Map: map[K]Expirable[V]{},
		TTL: ttl,
	}
}

func (e *ExpirableMap[K, V]) Set(key K, value V) {
	e.Map[key] = Expirable[V]{
		Value:  value,
		Expiry: time.Now().Add(e.TTL),
	}
}

func (e *ExpirableMap[K, V]) Get(key K) *V {
	value, ok := e.Map[key]
	if ok {
		if time.Now().Before(value.Expiry) {
			return &value.Value
		} else {
			delete(e.Map, key) // Clean up expired entries
		}
	}
	return nil
}

func (e *ExpirableMap[K, V]) Len() int {
	return len(e.Map)
}
