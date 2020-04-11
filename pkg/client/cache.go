package client

import (
	"time"
)

type CacheInt interface {
	RetrieveKey(key string) (interface{}, bool)
	SetKey(key string, value interface{}, refresh int)
}

type cacheValue struct {
	val interface{}
	rt  int64
}

type simpleCache struct {
	flagCache map[string]*cacheValue
}

func newSimpleCache() *simpleCache {
	return &simpleCache{
		flagCache: make(map[string]*cacheValue)}
}

func (s *simpleCache) RetrieveKey(key string) (interface{}, bool) {
	val, ok := s.flagCache[key]
	if !ok {
		return nil, false
	}
	if val.rt > time.Now().Unix() {
		return val, true
	}
	// Expire key from cache
	delete(s.flagCache, key)
	return nil, false
}

func (s *simpleCache) SetKey(key string, value interface{}, refresh int) {
	rt := time.Now().Unix() + int64(refresh)
	s.flagCache[key] = &cacheValue{val: value, rt: rt}
}
