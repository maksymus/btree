package main

import (
  "btree/errors"
  . "github.com/smartystreets/goconvey/convey"
  "os"
  "testing"
)

var filename = "/tmp/test.dat"

func Test_newPaged(t *testing.T) {
  Convey("create new paged", t, func() {
    paged := newPaged("test.dat", DefaultConfig())

    Convey("paged created with default values", func() {
      So(paged.pageSize, ShouldEqual, DefaultPageSize)
      So(paged.pageHeaderSize, ShouldEqual, DefaultPageHeaderSize)
      So(paged.headerSize, ShouldEqual, DefaultHeaderSize)
      So(paged.maxKeySize, ShouldEqual, DefaultMaxKeySize)
      So(paged.pageCount, ShouldEqual, DefaultPageCount)

      So(paged.firstFreePage, ShouldEqual, NoPage)
      So(paged.lastFreePage, ShouldEqual, NoPage)
      So(paged.recordCount, ShouldEqual, 0)
      So(paged.totalCount, ShouldEqual, 0)
    })
  })
}

func Test_open_FileMissing(t *testing.T) {
  Convey("paged - create new paged file", t, func() {
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

func Test_writeValue(t *testing.T) {
  Convey("paged opened and data written to page", t, func() {
    os.Remove(filename)

    paged := newPaged(filename, DefaultConfig())
    if err := paged.open(); err != nil {
      t.Errorf("failed to open paged with error: %s\n%s", err, errors.Stack(err))
    }

    freePage, err := paged.getPage(0)
    if err != nil {
      t.Errorf("failed to get page: %s\n%s", err, errors.Stack(err))
    }

    freePage.setStatus(127)
    err = paged.writeValue(freePage, NewValue([]byte{ 1, 2, 3}, 0, 3))
    if err != nil {
      t.Errorf("failed to write data: %s\n%s", err, errors.Stack(err))
    }
    paged.close()

    Convey("paged reopened and data is persisted", func() {
      paged1 := newPaged(filename, DefaultConfig())
      if err := paged1.open(); err != nil {
        t.Errorf("failed to open paged with error: %s\n%s", err, errors.Stack(err))
      }
      defer paged1.close()

      // check page 0
      page, err := paged1.getPage(0)
      if err != nil {
        t.Errorf("failed to fetch page: %s\n%s", err, errors.Stack(err))
      }


      So(page.pageHeader.status, ShouldEqual, 127)
      So(page.pageHeader.nextPage, ShouldEqual, -1)
      So(page.pageHeader.keyLength, ShouldEqual, 0)
      So(page.pageHeader.keyHash, ShouldEqual, 0)
      So(page.pageHeader.recordLength, ShouldEqual, 3)
      So(page.pageHeader.dataLength, ShouldEqual, 3)

      value, err := paged1.readValue(page)
      So(err, ShouldBeNil)

      So(value.data, ShouldResemble, []byte{ 1, 2, 3})
    })
  })
}


