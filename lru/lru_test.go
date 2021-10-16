package lru

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type String string

func (d String) len() int {
	return len(d)
}

func TestCache_Get(t *testing.T) {
	lru := NewLRUCache(int64(200), nil)
	lru.Put("key1", String("1234"))

	v1, ok1 := lru.Get("key1")
	assert.True(t, ok1)
	assert.Equal(t, "1234", string(v1.(String)))

	_, ok2 := lru.Get("key2")
	assert.False(t, ok2)
}


func TestCache_RemoveLeastRecently(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewLRUCache(int64(cap), nil)
	lru.Put(k1, String(v1))
	lru.Put(k2, String(v2))
	lru.Put(k3, String(v3))

	_, ok := lru.Get("key1")
	assert.False(t, ok)
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value CacheNodeValue) {
		keys = append(keys, key)
	}
	lru := NewLRUCache(int64(10), callback)
	lru.Put("key1", String("123456"))
	lru.Put("k2", String("k2"))
	lru.Put("k3", String("k3"))
	lru.Put("k4", String("k4"))

	expect := []string{"key1", "k2"}

	assert.True(t, reflect.DeepEqual(expect, keys))
}