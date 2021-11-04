package glowingcache

import (
	"fmt"
	"sync"
)

type Group struct {
	name      string
	getter    Getter
	mainCache *cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxSize int64, getter Getter) *Group {
	if getter == nil {
		panic("Nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	group := &Group{
		name:   name,
		getter: getter,
		mainCache: &cache{
			maxSize: maxSize,
		},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	group := groups[name]
	mu.Unlock()
	return group
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
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
		return ByteView{}, nil
	}
	value := ByteView{
		b: cloneBytes(bytes),
	}
	g.mainCache.add(key, value)
	return value, nil
}

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
