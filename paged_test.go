package main

import (
  . "github.com/smartystreets/goconvey/convey"
  "testing"
)

func Test_newPaged(t *testing.T) {
  Convey("create new paged", t, func() {

    paged, err := newPaged("test.dat", DefaultConfig())

    Convey("paged created with default values", func() {
      So(err, ShouldBeNil)
      So(paged.fileHeader.PageSize, ShouldEqual, DefaultPageSize)
      So(paged.fileHeader.PageHeaderSize, ShouldEqual, DefaultPageHeaderSize)
      So(paged.fileHeader.HeaderSize, ShouldEqual, DefaultHeaderSize)
      So(paged.fileHeader.MaxKeySize, ShouldEqual, DefaultMaxKeySize)
      So(paged.fileHeader.PageCount, ShouldEqual, DefaultPageCount)

      So(paged.fileHeader.FirstFreePage, ShouldEqual, 0)
      So(paged.fileHeader.LastFreePage, ShouldEqual, 0)
      So(paged.fileHeader.RecordCount, ShouldEqual, 0)
      So(paged.fileHeader.TotalCount, ShouldEqual, 0)
    })
  })
}
