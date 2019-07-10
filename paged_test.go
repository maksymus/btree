package main

import (
  "btree/errors"
  . "github.com/smartystreets/goconvey/convey"
  "os"
  "testing"
)

func Test_newPaged(t *testing.T) {
  Convey("create new paged", t, func() {
    paged, err := newPaged("test.dat", DefaultConfig())

    Convey("paged created with default values", func() {
      So(err, ShouldBeNil)

      header := paged.fileHeader

      So(header.PageSize, ShouldEqual, DefaultPageSize)
      So(header.PageHeaderSize, ShouldEqual, DefaultPageHeaderSize)
      So(header.HeaderSize, ShouldEqual, DefaultHeaderSize)
      So(header.MaxKeySize, ShouldEqual, DefaultMaxKeySize)
      So(header.PageCount, ShouldEqual, DefaultPageCount)

      So(header.FirstFreePage, ShouldEqual, 0)
      So(header.LastFreePage, ShouldEqual, 0)
      So(header.RecordCount, ShouldEqual, 0)
      So(header.TotalCount, ShouldEqual, 0)
    })
  })
}

func Test_open_FileMissing(t *testing.T) {
  Convey("create new paged", t, func() {
    filename := "/tmp/test.dat"
    os.Remove(filename)

    defer func() {
      os.Remove(filename)
    }()

    paged, err := newPaged(filename, DefaultConfig())
    if err != nil {
      t.Errorf("failed to init paged with error %s\n%s", err, errors.Stack(err))
    }

    if err := paged.open(); err != nil {
      t.Errorf("failed to open paged with error %s\n%s", err, errors.Stack(err))
    }

    Convey("paged opened - file created", func() {
      if stat, err := os.Stat(filename); os.IsNotExist(err) {
        t.Error("file not created")
      } else {
        So(stat.Size(), ShouldEqual, 50)
      }
    })
  })
}

func Test_open_FileExists_NotOpen(t *testing.T) {
  Convey("create new paged", t, func() {
    filename := "/tmp/test.dat"
    os.Remove(filename)

    // defer func() {
    //   os.Remove(filename)
    // }()

    paged, err := newPaged(filename, DefaultConfig())
    if err != nil {
      t.Errorf("failed to init paged with error: %s\n%s", err, errors.Stack(err))
    }

    if err := paged.open(); err != nil {
      t.Errorf("failed to open paged with error: %s\n%s", err, errors.Stack(err))
    }

    if err := paged.close(); err != nil {
      t.Errorf("failed to close paged with error: %s\n%s", err, errors.Stack(err))
    }

    if err := paged.open(); err != nil {
      t.Errorf("failed to re-open paged with error: %s\n%s", err, errors.Stack(err))
    }

    Convey("paged opened - file created", func() {
      if stat, err := os.Stat(filename); os.IsNotExist(err) {
        t.Error("file not created")
      } else {
        So(stat.Size(), ShouldEqual, 50)
      }
    })
  })
}
