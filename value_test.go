package main

import "testing"
import . "github.com/smartystreets/goconvey/convey"

func Test_Buffer(t *testing.T) {
  Convey("get buffer from value", t, func() {
    bs := []byte{ 0, 1, 2, 3, 4, 5}

    value := NewValue(bs, 1, 3)
    buffer := value.Buffer()

    Convey("buffer should contain bytes from 1 to 3", func() {
      So(buffer, ShouldNotBeNil)

      So(buffer.Len(), ShouldEqual, 3)
      So(buffer.Bytes(), ShouldResemble, []byte{ 1, 2, 3 })
    })
  })
}