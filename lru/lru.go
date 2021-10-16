package lru

import (
	"container/list"
)

type Cache struct {
	doubleList *list.List
	cacheMap   map[string]*list.Element
	maxMemory  int64
	usedMemory int64
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value CacheNodeValue)
}

type cacheNode struct {
	key   string
	value CacheNodeValue
}

type CacheNodeValue interface {
	len() int
}

func NewLRUCache(maxMemo int64, onEvicted func(string, CacheNodeValue)) *Cache {
	return &Cache{
		doubleList: list.New(),
		cacheMap:   make(map[string]*list.Element),
		OnEvicted:  onEvicted,
		maxMemory:  maxMemo,
	}
}

func (cache *Cache) Get(key string) (value CacheNodeValue, ok bool) {
	if element, ok := cache.cacheMap[key]; ok {
		cache.doubleList.MoveToFront(element)
		node := element.Value.(*cacheNode)
		return node.value, true
	}
	return nil, false
}

func (cache *Cache) RemoveLeastRecently() {
	element := cache.doubleList.Back()
	if element != nil {
		cache.doubleList.Remove(element)
		node := element.Value.(*cacheNode)
		delete(cache.cacheMap, node.key)
		cache.usedMemory -= calculateCacheNodeMemory(node)
		if cache.OnEvicted != nil {
			cache.OnEvicted(node.key, node.value)
		}
	}
}

func (cache *Cache) Put(key string, value CacheNodeValue) {
	if element, ok := cache.cacheMap[key]; ok {
		cache.doubleList.MoveToFront(element)
		node := element.Value.(*cacheNode)
		cache.usedMemory += int64(value.len()) - int64(node.value.len())
		node.value = value
	} else {
		front := cache.doubleList.PushFront(&cacheNode{
			key:   key,
			value: value,
		})
		cache.cacheMap[key] = front
		cache.usedMemory += calculateCacheNodeMemory(front.Value.(*cacheNode))
	}
	for cache.maxMemory != 0 && cache.usedMemory > cache.maxMemory {
		cache.RemoveLeastRecently()
	}
}

func calculateCacheNodeMemory(node *cacheNode) int64 {
	return int64(len(node.key)) + int64(node.value.len())
}
