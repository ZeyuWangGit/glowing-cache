package singleflight

import "sync"

type singleCall struct {
	waitGroup sync.WaitGroup
	value     interface{}
	err       error
}

type Group struct {
	mu      sync.Mutex
	callMap map[string]*singleCall
}

func (group Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	group.mu.Lock()
	if group.callMap == nil {
		group.callMap = make(map[string]*singleCall)
	}

	if call, ok := group.callMap[key]; ok {
		call.waitGroup.Wait()
		return call.value, call.err
	}

	call := new(singleCall)
	call.waitGroup.Add(1)
	group.callMap[key] = call
	group.mu.Unlock()

	call.value, call.err = fn()
	call.waitGroup.Done()

	group.mu.Lock()
	delete(group.callMap, key)
	group.mu.Unlock()

	return call.value, call.err
}

