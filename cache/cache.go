package cache

type Cache interface {
  Put(key, value interface{})
  Get(key interface{}) (interface{}, bool)
  Len() int
}
