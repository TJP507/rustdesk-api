package cache

import (
	"errors"
	"reflect"
	"sync"
)

// This implements a simple cache for testing purposes
// SimpleCache is a simple cache implementation
type SimpleCache struct {
	data      map[string]interface{}
	mu        sync.Mutex
	maxBytes  int64
	usedBytes int64
}

func (s *SimpleCache) Get(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Use reflection to set the stored value into the provided pointer variable
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}
	v, ok := s.data[key]
	if !ok {
		// Set to the zero value
		val.Elem().Set(reflect.Zero(val.Elem().Type()))
		return nil
	}

	vval := reflect.ValueOf(v)
	if val.Elem().Type() != vval.Type() {
		// Set to the zero value
		val.Elem().Set(reflect.Zero(val.Elem().Type()))
		return nil
	}

	val.Elem().Set(reflect.ValueOf(v))
	return nil
}

func (s *SimpleCache) Set(key string, value interface{}, exp int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Check whether the provided value is a pointer; if so, dereference it
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	s.data[key] = val.Interface()
	return nil
}
func (s *SimpleCache) Gc() error {
	return nil
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{
		data: make(map[string]interface{}),
	}
}
