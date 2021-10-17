package glowingcache

import (
	"fmt"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache *cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxMemory int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	group := &Group{
		name: name,
		getter: getter,
		mainCache: &cache{
			maxMemory: maxMemory,
		},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLocker()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("Key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		return v, nil
	}

	return g.loadFromSource(key)
}

func (g *Group) loadFromSource(key string) (ByteView, error) {
	return g.getFromLocal(key)
}

func (g *Group) getFromLocal(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{byteView: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

