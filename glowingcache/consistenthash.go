package glowingcache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func(data []byte) uint32

type Map struct {
	// Hash Algorithm Func
	hashFunc HashFunc
	// How many virtual nodes per node
	replicas int
	// hash keys circle
	hashCircle []int
	// map virtual node to real node
	hashMap map[int]string
}

func New(replicas int, fn HashFunc) *Map {
	m := &Map{
		replicas: replicas,
		hashFunc: fn,
		hashMap: make(map[int]string),
	}
	// default hash algorithm using crc32.ChecksumIEEE
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(nodeNames ...string) {
	for _, nodeName := range nodeNames {
		for i := 0; i < m.replicas; i++ {
			virtualNodeName := nodeName + "-" + strconv.Itoa(i)
			virtualNodeHash := int(m.hashFunc([]byte(virtualNodeName)))
			m.hashCircle = append(m.hashCircle, virtualNodeHash)
			m.hashMap[virtualNodeHash] = nodeName
		}
	}
	sort.Ints(m.hashCircle)
}

func (m *Map) Get(key string) string {
	if len(m.hashCircle) == 0 {
		return ""
	}

	givenHash := int(m.hashFunc([]byte(key)))
	index := sort.Search(len(m.hashCircle), func(i int) bool {
		return m.hashCircle[i] >= givenHash
	})
	virtualNodeHash := m.hashCircle[index % len(m.hashCircle)]
	realNodeName := m.hashMap[virtualNodeHash]
	return realNodeName
}