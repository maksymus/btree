package cache

import "container/list"

type lruCache struct {
  // max size of cache
  maxSize int

  // list data stuct
  list *list.List

  // map of elements
  elements map[interface{}] *list.Element
}

func NewLruCache(maxSize int) Cache {
  return &lruCache{
    maxSize:  maxSize,
    list:     list.New(),
    elements: make(map[interface{}]*list.Element),
  }
}

func (cache *lruCache) Put(key, value interface{}) {
  if elem, ok := cache.elements[key]; ok {
    cache.list.MoveToFront(elem)
    elem.Value = value
  } else {
     if cache.list.Len() >= cache.maxSize {
       cache.list.Remove(cache.list.Back())
       delete(cache.elements, key)
     }

    elem := cache.list.PushFront(value)
    cache.elements[key] = elem
  }
}

func (cache *lruCache) Get(key interface{}) (interface{}, bool) {
  if elem, ok := cache.elements[key]; ok {
    cache.list.MoveToFront(elem)
    return elem.Value, true
  }

  return nil, false
}

func (cache *lruCache) Len() int {
  return cache.list.Len()
}


