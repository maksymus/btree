package main

import (
  "btree/errors"
  . "github.com/smartystreets/goconvey/convey"
  "os"
  "testing"
)

func Test_newPaged(t *testing.T) {
  Convey("create new paged", t, func() {
    paged := newPaged("test.dat", DefaultConfig())

    Convey("paged created with default values", func() {
      So(paged.pageSize, ShouldEqual, DefaultPageSize)
      So(paged.pageHeaderSize, ShouldEqual, DefaultPageHeaderSize)
      So(paged.headerSize, ShouldEqual, DefaultHeaderSize)
      So(paged.maxKeySize, ShouldEqual, DefaultMaxKeySize)
      So(paged.pageCount, ShouldEqual, DefaultPageCount)

      So(paged.firstFreePage, ShouldEqual, 0)
      So(paged.lastFreePage, ShouldEqual, 0)
      So(paged.recordCount, ShouldEqual, 0)
      So(paged.totalCount, ShouldEqual, 0)
    })
  })
}

func Test_open_FileMissing(t *testing.T) {
  Convey("paged - create new paged file", t, func() {
    filename := "/tmp/test.dat"
    os.Remove(filename)

    defer func() {
      os.Remove(filename)
    }()

    paged := newPaged(filename, DefaultConfig())
    if err := paged.open(); err != nil {
      t.Errorf("failed to open paged with error %s\n%s", err, errors.Stack(err))
    }

    Convey("paged opened - file created", func() {
      if stat, err := os.Stat(filename); os.IsNotExist(err) {
        t.Error("file not created")
      } else {
        So(stat.Size(), ShouldEqual, 49)
      }
    })
  })
}

func Test_open_FileExists_NotOpen(t *testing.T) {
  Convey("paged - open existing file", t, func() {
    filename := "/tmp/test.dat"
    os.Remove(filename)

    paged := newPaged(filename, DefaultConfig())
    if err := paged.open(); err != nil {
      t.Errorf("failed to open paged with error: %s\n%s", err, errors.Stack(err))
    }

    if err := paged.close(); err != nil {
      t.Errorf("failed to close paged with error: %s\n%s", err, errors.Stack(err))
    }

    if err := paged.open(); err != nil {
      t.Errorf("failed to re-open paged with error: %s\n%s", err, errors.Stack(err))
    }

    Convey("open existing file", func() {
      if stat, err := os.Stat(filename); os.IsNotExist(err) {
        t.Error("failed to open existing file")
      } else {
        So(stat.Size(), ShouldEqual, 49)
      }
    })
  })
}

func Test_read_write_FileHeader(t *testing.T) {
  Convey("paged reopened", t, func() {
    filename := "/tmp/test.dat"
    os.Remove(filename)

    config := Config{
      headerSize:     1000,
      pageSize:       1001,
      pageCount:      1002,
      maxKeySize:     1003,
      pageHeaderSize: 32,
    }

    paged := newPaged(filename, config)
    if err := paged.open(); err != nil {
      t.Errorf("failed to open paged with error: %s\n%s", err, errors.Stack(err))
    }
    paged.close()

    paged1 := newPaged(filename, DefaultConfig())
    if err := paged1.open(); err != nil {
      t.Errorf("failed to open paged with error: %s\n%s", err, errors.Stack(err))
    }
    paged1.close()

    Convey("paged file header is populated from file", func() {
      if _, err := os.Stat(filename); os.IsNotExist(err) {
        t.Error("failed to open existing file")
      } else {
        So(paged1.headerSize, ShouldEqual, 1000)
        So(paged1.pageSize, ShouldEqual, 1001)
        So(paged1.pageCount, ShouldEqual, 1002)
        So(paged1.maxKeySize, ShouldEqual, 1003)
        So(paged1.pageHeaderSize, ShouldEqual, 32)
      }
    })
  })
}

