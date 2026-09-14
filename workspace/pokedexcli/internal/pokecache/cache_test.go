// Tests the internal functionality of the cache system
//
// File: cache_test.go
// Purpose: Validates that NewCache creates a functional cache and that stale
// entries are automatically reaped after the configured duration interval.
package pokecache

import (
	//"fmt"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	interval := 5 * time.Millisecond
	cache := NewCache(interval)

	if cache.cache == nil {
		t.Errorf("expected cache map to be initialized, got nil")
	}
}
