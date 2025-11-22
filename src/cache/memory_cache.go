// Package cache
package cache

import (
	"metar-provider/src/interfaces/cache"
	"metar-provider/src/utils"
	"sync"
	"time"
)

func isOutDate[T any](data *cache.CachedItem[T]) bool {
	return data.ExpiredAt.Before(time.Now())
}

type MemoryCache[T any] struct {
	cacheMap map[string]*cache.CachedItem[T]
	cleaner  *utils.IntervalActuator
	lock     sync.RWMutex
}

func NewMemoryCache[T any](cleanInterval time.Duration) *MemoryCache[T] {
	if cleanInterval <= 0 {
		cleanInterval = 30 * time.Minute
	}
	cached := &MemoryCache[T]{
		cacheMap: make(map[string]*cache.CachedItem[T]),
		lock:     sync.RWMutex{},
	}
	cached.cleaner = utils.NewIntervalActuator(cleanInterval, cached.CleanExpiredData)
	return cached
}

func (memoryCache *MemoryCache[T]) CleanExpiredData() {
	memoryCache.lock.Lock()
	defer memoryCache.lock.Unlock()

	for key, value := range memoryCache.cacheMap {
		if isOutDate(value) {
			delete(memoryCache.cacheMap, key)
		}
	}
}

func (memoryCache *MemoryCache[T]) Set(key string, value T, expiredAt time.Time) {
	if expiredAt.Before(time.Now()) {
		return
	}
	if key == "" {
		return
	}
	memoryCache.lock.Lock()
	memoryCache.cacheMap[key] = &cache.CachedItem[T]{CachedData: value, ExpiredAt: expiredAt}
	memoryCache.lock.Unlock()
}

func (memoryCache *MemoryCache[T]) SetWithTTL(key string, value T, ttl time.Duration) {
	expiredAt := time.Now().Add(ttl)
	memoryCache.Set(key, value, expiredAt)
}

func (memoryCache *MemoryCache[T]) Get(key string) (T, bool) {
	if key == "" {
		var zero T
		return zero, false
	}
	memoryCache.lock.RLock()
	defer memoryCache.lock.RUnlock()
	val, ok := memoryCache.cacheMap[key]
	if ok && isOutDate(val) {
		var zero T
		return zero, false
	}
	if val == nil {
		var zero T
		return zero, false
	}
	return val.CachedData, ok
}

func (memoryCache *MemoryCache[T]) Del(key string) {
	if key == "" {
		return
	}
	memoryCache.lock.Lock()
	delete(memoryCache.cacheMap, key)
	memoryCache.lock.Unlock()
}

func (memoryCache *MemoryCache[T]) Close() {
	memoryCache.cleaner.Stop()
}
