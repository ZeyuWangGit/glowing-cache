package glowingcache

import (
	"glowing-cache/glowingcache/lru"
	"sync"
)

type cache struct {
	mutexLock sync.Mutex
	lruCache  *lru.Cache
	maxMemory int64
}

func (c *cache) add(key string, value ByteView) {
	c.mutexLock.Lock()
	defer c.mutexLock.Unlock()
	if c.lruCache == nil {
		c.lruCache = lru.NewLRUCache(c.maxMemory, nil)
	}
	c.lruCache.Put(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mutexLock.Lock()
	defer c.mutexLock.Unlock()
	if c.lruCache == nil {
		return
	}
	if v, ok := c.lruCache.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
