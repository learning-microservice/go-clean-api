package syncmap

import "sync"

type Map[T any] struct {
	mu sync.RWMutex
	kv map[string]*T
}

func New[T any](capacity int, initKeyValues ...map[string]*T) *Map[T] {
	syncMap := &Map[T]{
		kv: make(map[string]*T, capacity),
	}
	for i := range initKeyValues {
		for key, value := range initKeyValues[i] {
			syncMap.Store(key, value)
		}
	}
	return syncMap
}

// Load -.
func (sm *Map[T]) Load(key string) (*T, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	value, ok := sm.kv[key]
	if !ok {
		return nil, false
	}
	return value, true
}

// Store -.
func (sm *Map[T]) Store(key string, value *T) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.kv[key] = value
}
