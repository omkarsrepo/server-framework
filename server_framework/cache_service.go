package server_framework

import (
	"github.com/maypok86/otter"
	"github.com/samber/lo"
	"sync"
	"time"
)

var (
	cacheServiceInstance *cacheService
	cacheServiceOnce     sync.Once
	cacheMapsMtx         sync.Mutex
	variableCacheMapsMtx sync.Mutex
)

type CacheService interface {
	New(capacity int, ttl time.Duration) otter.Cache[string, any]
	NewVariable(capacity int) otter.CacheWithVariableTTL[string, any]
	Close()
}

type cacheService struct {
	cacheMaps         []*otter.Cache[string, any]
	variableCacheMaps []*otter.CacheWithVariableTTL[string, any]
}

func Cache() CacheService {
	cacheServiceOnce.Do(func() {
		cacheServiceInstance = &cacheService{}
	})

	return cacheServiceInstance
}

func (props *cacheService) registerCache(cache *otter.Cache[string, any]) {
	cacheMapsMtx.Lock()
	defer cacheMapsMtx.Unlock()

	props.cacheMaps = append(props.cacheMaps, cache)
}

func (props *cacheService) registerVariableCache(cache *otter.CacheWithVariableTTL[string, any]) {
	variableCacheMapsMtx.Lock()
	defer variableCacheMapsMtx.Unlock()

	props.variableCacheMaps = append(props.variableCacheMaps, cache)
}

func (props *cacheService) New(capacity int, ttl time.Duration) otter.Cache[string, any] {
	cache, err := otter.MustBuilder[string, any](capacity).
		WithTTL(ttl).
		Build()

	if err != nil {
		panic(err)
	}

	props.registerCache(&cache)

	return cache
}

func (props *cacheService) NewVariable(capacity int) otter.CacheWithVariableTTL[string, any] {
	cache, err := otter.MustBuilder[string, any](capacity).
		WithVariableTTL().Build()
	if err != nil {
		panic(err)
	}

	props.registerVariableCache(&cache)

	return cache
}

func (props *cacheService) Close() {
	cacheMapsMtx.Lock()
	defer cacheMapsMtx.Unlock()

	lo.ForEach(props.cacheMaps, func(cache *otter.Cache[string, any], _ int) {
		cache.Close()
	})

	variableCacheMapsMtx.Lock()
	defer variableCacheMapsMtx.Unlock()

	lo.ForEach(props.variableCacheMaps, func(cache *otter.CacheWithVariableTTL[string, any], _ int) {
		cache.Close()
	})
}
