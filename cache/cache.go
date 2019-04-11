package cache

type cache interface {
  put(uint, *struct{})
  get(uint) *struct{}
}

type lruCache struct {

}

func (lruCache) put(uint, *struct{}) {
  panic("implement me")
}

func (lruCache) get(uint) *struct{} {
  panic("implement me")
}

