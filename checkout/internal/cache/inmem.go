package cache

import (
	"sync"
	"time"
)

type Cache[T comparable, V any] interface {
	Set(key T, value V, ttl time.Duration)
	Get(key T) (V, bool)
	Delete(key T)
	Clear()
}

var _ Cache[string, string] = &inMemCache[string, string]{}

type inMemCache[T comparable, V any] struct {
	items map[T]*inMemCacheItem[V]
	mutex sync.RWMutex
}

type inMemCacheItem[V any] struct {
	value      V
	expiration int64
}

func NewInMemCache[T comparable, V any]() *inMemCache[T, V] {
	c := &inMemCache[T, V]{
		items: make(map[T]*inMemCacheItem[V]),
	}
	return c
}

// Set Установить значение ключа
func (c *inMemCache[T, V]) Set(key T, value V, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	expiration := time.Now().Add(ttl).UnixNano()
	c.items[key] = &inMemCacheItem[V]{
		value:      value,
		expiration: expiration,
	}
}

// Get Получить значение ключа
func (c *inMemCache[T, V]) Get(key T) (V, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, found := c.items[key]

	var noop V

	if !found {
		return noop, false
	}
	if time.Now().UnixNano() > item.expiration {
		delete(c.items, key)
		return noop, false
	}
	return item.value, true
}

// Delete Удалить ключ
func (c *inMemCache[T, V]) Delete(key T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.items, key)
}

// Clear Очистка кэша
func (c *inMemCache[T, V]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items = make(map[T]*inMemCacheItem[V])
}
