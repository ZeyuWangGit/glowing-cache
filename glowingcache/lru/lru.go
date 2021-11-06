package lru

import "container/list"

// Cache LRU with
// A double List
// A map
// max memory size
// used memory size
// callback function when an entry is purged OnEvicted
type Cache struct {
	doubleList *list.List
	cacheMap map[string]*list.Element
	maxMemory int64
	usedMemory int64
	onEvicted func(key string, value CacheNodeValue)
}

// cacheNode with key, value pair
type cacheNode struct {
	key string
	value CacheNodeValue
}

type CacheNodeValue interface {
	Len() int
}

// New constructor of cache
func New(maxMemo int64, onEvicted func(key string, value CacheNodeValue)) *Cache{
	return &Cache{
		doubleList: list.New(),
		cacheMap: make(map[string]*list.Element),
		maxMemory: maxMemo,
		usedMemory: 0,
		onEvicted: onEvicted,
	}
}

// Get look ups a key's value
func (cache *Cache) Get(key string) (CacheNodeValue, bool) {
	if element, ok := cache.cacheMap[key]; ok {
		cache.doubleList.MoveToFront(element)
		node := element.Value.(*cacheNode)
		return node.value, true
	}
	return nil, false
}

// RemoveLeastRecently removes the oldest item
func (cache *Cache) RemoveLeastRecently()  {
	element := cache.doubleList.Back()
	if element != nil {
		cache.doubleList.Remove(element)
		node := element.Value.(*cacheNode)
		delete(cache.cacheMap, node.key)
		cache.usedMemory -= getCacheNodeMemory(node)
		if cache.onEvicted != nil {
			cache.onEvicted(node.key, node.value)
		}
	}
}

// Put adds or update value to the cache.
func (cache *Cache) Put(key string, value CacheNodeValue) {
	if element, ok := cache.cacheMap[key]; ok {
		cache.doubleList.MoveToFront(element)
		node := element.Value.(*cacheNode)
		cache.usedMemory += int64(value.Len()) - int64(node.value.Len())
		node.value = value
	} else {
		node := &cacheNode{
			key: key,
			value: value,
		}
		el := cache.doubleList.PushFront(node)
		cache.cacheMap[key] = el
		cache.usedMemory += getCacheNodeMemory(node)
	}
	for cache.maxMemory != 0 && cache.usedMemory > cache.maxMemory {
		cache.RemoveLeastRecently()
	}
}

func getCacheNodeMemory(node *cacheNode) int64 {
	return int64(len(node.key)) + int64(node.value.Len())
}






