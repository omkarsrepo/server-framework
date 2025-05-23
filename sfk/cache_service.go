// Unpublished Work © 2024

package sfk

import (
	"github.com/maypok86/otter"
	"github.com/samber/lo"
	"sync"
	"time"
)

var (
	cacheServiceInstance *cacheService
	cacheServiceOnce     sync.Once
	cacheMapsMtx         sync.RWMutex
	variableCacheMapsMtx sync.RWMutex
)

type CacheService interface {
	New(capacity int, ttl time.Duration) otter.Cache[string, any]
	NewVariable(capacity int) otter.CacheWithVariableTTL[string, any]
	Close()
}

type cacheService struct {
	cacheMaps         []otter.Cache[string, any]
	variableCacheMaps []otter.CacheWithVariableTTL[string, any]
}

func Cache() CacheService {
	cacheServiceOnce.Do(func() {
		cacheServiceInstance = &cacheService{}
	})

	return cacheServiceInstance
}

func (c *cacheService) registerCache(cache otter.Cache[string, any]) {
	cacheMapsMtx.Lock()
	defer cacheMapsMtx.Unlock()

	c.cacheMaps = append(c.cacheMaps, cache)
}

func (c *cacheService) registerVariableCache(cache otter.CacheWithVariableTTL[string, any]) {
	variableCacheMapsMtx.Lock()
	defer variableCacheMapsMtx.Unlock()

	c.variableCacheMaps = append(c.variableCacheMaps, cache)
}

func (c *cacheService) New(capacity int, ttl time.Duration) otter.Cache[string, any] {
	cache, err := otter.MustBuilder[string, any](capacity).
		WithTTL(ttl).
		Build()

	if err != nil {
		panic(err)
	}

	c.registerCache(cache)

	return cache
}

func (c *cacheService) NewVariable(capacity int) otter.CacheWithVariableTTL[string, any] {
	cache, err := otter.MustBuilder[string, any](capacity).
		WithVariableTTL().Build()
	if err != nil {
		panic(err)
	}

	c.registerVariableCache(cache)

	return cache
}

func (c *cacheService) Close() {
	cacheMapsMtx.RLock()
	defer cacheMapsMtx.RUnlock()

	lo.ForEach(c.cacheMaps, func(cache otter.Cache[string, any], _ int) {
		cache.Close()
	})

	variableCacheMapsMtx.RLock()
	defer variableCacheMapsMtx.RUnlock()

	lo.ForEach(c.variableCacheMaps, func(cache otter.CacheWithVariableTTL[string, any], _ int) {
		cache.Close()
	})
}
