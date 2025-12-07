package pongo2

import (
	"reflect"
	"sync"
)

// structFieldCache caches struct field indices for fast lookup
type structFieldCache struct {
	mu    sync.RWMutex
	cache map[reflect.Type]map[string][]int // Type -> FieldName -> Field index path
}

var globalStructFieldCache = &structFieldCache{
	cache: make(map[reflect.Type]map[string][]int),
}

// getFieldIndex returns the cached field index path or computes and caches it
func (c *structFieldCache) getFieldIndex(typ reflect.Type, fieldName string) ([]int, bool) {
	c.mu.RLock()
	if fields, ok := c.cache[typ]; ok {
		if indices, ok := fields[fieldName]; ok {
			c.mu.RUnlock()
			return indices, true
		}
	}
	c.mu.RUnlock()

	// Not in cache, compute it
	field, ok := typ.FieldByName(fieldName)
	if !ok {
		return nil, false
	}

	// Cache the result
	c.mu.Lock()
	if c.cache[typ] == nil {
		c.cache[typ] = make(map[string][]int)
	}
	c.cache[typ][fieldName] = field.Index
	c.mu.Unlock()

	return field.Index, true
}
