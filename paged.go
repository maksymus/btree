package main

import (
  "bytes"
  "encoding/binary"
  "fmt"
  "io"
  "os"

  "btree/errors"
  log "github.com/sirupsen/logrus"
)

// https://xml.apache.org/xindice/dev/guide-internals.html#3.+Data+storage
// https://github.com/myui/btree4j/blob/master/src/main/java/btree4j/Paged.java

// Paging provides efficient access to a random-access file by allowing parts of the file (pages) to be "mapped" to main memory for easy access.
// Pages have a fixed length.
// If data that must be stored is longer than the length of one page, subsequent pages in the file can be "linked" to the first.
// The file header is 4kb long, and each page is, by default, 4kb long.

// paged abstraction for paged file
type paged struct {
  isOpen   bool
  filename string
  file     *os.File

  config Config

  fileHeader *fileHeader
}

// create new paged
func newPaged(filename string, config Config) (*paged, error) {
  paged := &paged{}
  paged.filename = filename
  paged.config = config
  paged.fileHeader = newFileHeader(config)

  return paged, nil
}

// try to read paged file or create if file doesn't exists
// returns error in following cases
// - file exists but not readable (format broken)
// - file doesn't exist and failed to create file
func (paged *paged) open() error {
  if paged.isOpen {
    return errors.New("paged already opened")
  }

  var errors error

  if paged.file, errors = os.OpenFile(paged.filename, os.O_RDWR, 0666); errors == nil {
    if errors = paged.read(); errors != nil {
      log.WithError(errors).Panic("failed to read file")
    }
  } else if os.IsNotExist(errors) {
    if paged.file, errors = os.OpenFile(paged.filename, os.O_RDWR | os.O_CREATE, 0666); errors != nil {
      log.WithError(errors).Panic("failed to open file for creation")
    }

    if errors = paged.create(); errors != nil {
      log.WithError(errors).Panic("failed to create file")
    }
  } else {
    log.WithError(errors).Panic("failed to open file")
  }

  paged.isOpen = true
  return nil
}

// close paged file
func (paged *paged) close() error {
  if !paged.isOpen {
    return fmt.Errorf("file is not open")
  }

  if err := paged.flush(); err != nil {
    return errors.Wrap(err)
  }

  if err := paged.file.Close(); err != nil {
    return errors.Wrap(err)
  }

  paged.isOpen = false
  return nil
}

// create paged file and populate file header
func (paged *paged) create() error {
  paged.fileHeader = newFileHeader(paged.config)

  if err := paged.write(); err != nil {
    return errors.Wrap(err)
  }

  if err := paged.flush(); err != nil {
    return errors.Wrap(err)
  }

  return nil
}

func (paged *paged) flush() error {
  return nil
}

// read paged file header
func (paged *paged) read() error {
  var fh fileHeader

  if err := read(paged, 0, uint32(paged.config.headerSize), &fh); err != nil {
    return errors.WrapMsg(err, "failed to read page header")
  }

  paged.fileHeader = &fh
  return nil
}

// write paged file header
func (paged *paged) write() error {
  return write(paged, 0, paged.fileHeader)
}

// get page by page number
func (paged *paged) getPage(pageNum int64) (*page, error) {
  if pageNum < 0 {
    return nil, errors.New("negative page number")
  }

  // todo use lru/other cache
  page := newPage(paged, pageNum)
  if err := page.read(); err != nil {
    return nil, err
  }

  return page, nil
}

// read page value
// iterates over pages if value spans through several pages
func (paged *paged) readValue(page *page) (*Value, error) {
  recordLength := page.pageHeader.RecordLength
  buffer := bytes.NewBuffer(make([]byte, recordLength))

  currentPage := page

  for {
    if err := page.streamTo(buffer); err != nil {
      return nil, err
    }

    nextPageNum := currentPage.pageHeader.NextPage
    if nextPageNum == NoPage {
      break
    }

    if nextPage, err := paged.getPage(nextPageNum); err == nil {
      currentPage = nextPage
    } else {
      return nil, err
    }
  }

  // return &Value{data: buffer.Bytes()}, nil
  return NewValue(buffer.Bytes(), 0, buffer.Len()), nil
}

// write value to page or pages
func (paged *paged) writeValue(page *page, value *Value) error {
  pageHeader := page.pageHeader
  pageHeader.RecordLength = int32(value.len)

  buffer := value.Buffer()

  // if more data left in buffer then write to page
  for buffer.Len() > 0 {
    pageHeader := page.pageHeader

    if err := page.streamFrom(buffer); err != nil {
      return err
    }

    // write current page
    if err := page.write(); err != nil {
      return err
    }

    // if no more data to write then break
    if buffer.Len() > 0 {
      break
    }

    if nextPageNum := pageHeader.NextPage; nextPageNum == NoPage {
      // TODO get free page
    } else {
      if nextPage, err := paged.getPage(nextPageNum); err == nil {
        page = nextPage
      } else {
        return err
      }
    }
  }

  // clean up unused overflow pages
  if page.pageHeader.NextPage != NoPage {
    // TODO unlink pages
  }

  page.pageHeader.NextPage = NoPage

  return nil
}

// The paged file header consists of a number of fixed-length fields. Fields which are longer than one byte, are always
// stored in Big Endian format, which means the most significant byte is written at the lowest address.
type fileHeader struct {
  HeaderSize     int16 // header size (2 bytes): set to 4096 (0x1000), the size of this header.
  PageSize       int32 // page size (4 bytes): set to 4096 (0x00001000), the page size.
  PageCount      int64 // page count (8 bytes): this field is not used consistently. It is present mainly for historical reason.
  TotalCount     int64 // total page count (8 bytes): total number of pages present in this file.
  FirstFreePage  int64 // first free page (8 bytes): page number of the first unused page in this file.
  LastFreePage   int64 // last free page (8 bytes): page number of the last unused page in this file.
  PageHeaderSize int8  // page header size (1 byte): size of each page header. Set to 64 (0x40) by default.
  MaxKeySize     int16
  RecordCount    int64 // record count (8 bytes): number of records stored in this file.
}

func newFileHeader(config Config) *fileHeader {
  return &fileHeader{
    HeaderSize:     config.headerSize,
    PageSize:       config.pageSize,
    PageCount:      config.pageCount,
    PageHeaderSize: config.pageHeaderSize,
    MaxKeySize:     config.maxKeySize,
  }
}

// write data to paged file starting from offset
func write(p *paged, offset int64, data interface{}) error {
  var err error

  buf := new(bytes.Buffer)
  if err = binary.Write(buf, binary.BigEndian, data); err == nil {
    _, err = p.file.WriteAt(buf.Bytes(), offset)
  }

  return errors.Wrap(err)
}

// read data from paged file starting from offset
func read(p *paged, offset int64, size uint32, data interface{}) error {
  bs := make([]byte, size)
  if _, err := p.file.ReadAt(bs, offset); err != nil {
    if err != io.EOF {
      return errors.Wrap(err)
    }
  }

  if err := binary.Read(bytes.NewBuffer(bs), binary.BigEndian, data); err != nil {
    return errors.Wrap(err)
  }

  return nil
}

