package cache

import (
	"container/list"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Put_MoveToFront(t *testing.T) {
	Convey("put value", t, func() {
		cache := &lruCache{
			maxSize:  100,
			list:     list.New(),
			elements: make(map[interface{}]*list.Element),
		}

		cache.Put(10, "ten")
		cache.Put(11, "eleven")
		cache.Put(12, "twelve")
		cache.Put(10, "value")

		Convey("should contain one value", func() {
			So(len(cache.elements), ShouldEqual, 3)
			So(cache.list.Len(), ShouldEqual, 3)

			l := list.New()
			l.PushBack("value")
			l.PushBack("twelve")
			l.PushBack("eleven")

			m := make(map[interface{}]*list.Element)
			m[10] = l.Front()
			m[12] = l.Front().Next()
			m[11] = l.Front().Next().Next()

			So(cache.list, ShouldResemble, l)
			So(cache.elements, ShouldResemble, m)
		})
	})
}


func Test_Get(t *testing.T) {
}