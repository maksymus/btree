package main

import (
	"btree/cache"
	"testing"
)

func Test_cache(t *testing.T) {
	p := &page{
		pageNumber: 10,
	}

	pageCache := cache.NewLruCache(10)
	pageCache.Put(10, p)

	p.status = 123

	t.Run("check page is cached", func(t *testing.T) {
		if value, ok := pageCache.Get(10); ok {
			if cachedPage, ok := value.(*page); ok {
				if cachedPage != p {
					t.Errorf(`match wasn't as expected. a: %v, e: %v`, cachedPage, p)
				}
			} else {
				t.Error("cached value should be of type 'page'")
			}
		} else {
			t.Error("cached value is missing")
		}
	})
}