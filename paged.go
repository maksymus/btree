package main

import (
  "bytes"
  "encoding/binary"
  "fmt"
  "os"

  "github.com/hashicorp/go-multierror"
  log "github.com/sirupsen/logrus"
  "btree/errors"
)

// https://xml.apache.org/xindice/dev/guide-internals.html#3.+Data+storage
// https://github.com/myui/btree4j/blob/master/src/main/java/btree4j/Paged.java

// Paging provides efficient access to a random-access file by allowing parts of the file (pages) to be "mapped" to main memory for easy access.
// Pages have a fixed length.
// If data that must be stored is longer than the length of one page, subsequent pages in the file can be "linked" to the first.
// The file header is 4kb long, and each page is, by default, 4kb long.

const (
  DefaultHeaderSize     = 1024 * 4
  DefaultPageSize       = 1024 * 4
  DefaultPageCount      = 1024
  DefaultMaxKeySize     = 256
  DefaultPageHeaderSize = 64
)

func init() {
  log.SetFormatter(&log.TextFormatter{
    FullTimestamp: true,
  })
  log.SetReportCaller(true)
}

type Config struct {
  headerSize     int16 // header size
  pageSize       int32 // page size
  pageCount      int64 // init page count
  maxKeySize     int16 // max key size
  pageHeaderSize int8  // page header size
}

func DefaultConfig() Config {
  return Config{
    headerSize:     DefaultHeaderSize,
    pageSize:       DefaultPageSize,
    pageCount:      DefaultPageCount,
    maxKeySize:     DefaultMaxKeySize,
    pageHeaderSize: DefaultPageHeaderSize,
  }
}

type paged struct {
  isOpen   bool
  filename string
  file     *os.File

  config Config

  fileHeader *fileHeader
}

func newPaged(filename string, config Config) (*paged, error) {
  paged := &paged{}
  paged.filename = filename
  paged.config = config
  paged.fileHeader = newFileHeader(config)

  return paged, nil
}

func (paged *paged) open() error {
  if paged.isOpen {
    return errors.New("paged already opened")
  }

  var file *os.File
  var errors error

  if file, errors = os.OpenFile(paged.filename, os.O_RDWR, 0666); errors == nil {
    if errors = paged.read(); errors != nil {
      log.WithError(errors).Panic("failed to read file")
    }
  } else if os.IsNotExist(errors) {
    if file, errors = os.OpenFile(paged.filename, os.O_RDWR | os.O_CREATE, 0666); errors != nil {
      log.WithError(errors).Panic("failed to create file")
    }

    if errors = paged.create(); errors != nil {
      log.WithError(errors).Panic("failed to create file")
    }
  } else {
    log.WithError(errors).Panic("failed to open file")
  }

  paged.isOpen = true
  paged.file = file
  return nil
}

func (paged *paged) close() error {
  if !paged.isOpen {
    return fmt.Errorf("file is not open")
  }

  if err := paged.flush(); err != nil {
    return err
  }

  if err := paged.file.Close(); err != nil {
    return err
  }

  paged.isOpen = false
  return nil
}

func (paged *paged) create() error {
  return nil
}

func (paged *paged) flush() error {
  return nil
}

func (paged *paged) read() error {
  var fh *fileHeader

  if err := read(paged, 0, uint32(paged.fileHeader.HeaderSize), fh); err != nil {
    return errors.WrapMsg(err, "failed to read page header")
  }

  paged.fileHeader = fh
  paged.fileHeader.dirty = false
  return nil
}

func (paged *paged) write() error {
  return nil
}


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

  dirty bool
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

func write(p *paged, offset int64, data interface{}) error {
  var err error

  buf := new(bytes.Buffer)
  if err = binary.Write(buf, binary.BigEndian, data); err == nil {
    _, err = p.file.WriteAt(buf.Bytes(), offset)
  }

  return err
}

func read(p *paged, offset int64, size uint32, data interface{}) error {
  var err error

  bs := make([]byte, size)
  if _, err = p.file.ReadAt(bs, offset); err == nil {
    err = binary.Read(bytes.NewBuffer(bs), binary.BigEndian, data)
  }

  return err
}

// Each page contains a 64-byte header, followed by actual data. Pages are numbered.
// Whenever a particular page, say page n is needed, but not yet loaded into memory, the code can calculate the start
// address of the page as:
//
// offset = fileHeaderSize + (n * pageSize)
//
// At this address, it will then find the header of the wanted page, and 64 bytes further, the start of the page's data.
type pageHeader struct {
  Status       int8  // status (1 byte): pages in the data file are either used or unused. Used pages contain actual data.
  KeyLength    int16 // key length (2 bytes): pages have the possibility of storing a key just before their actual data.
  KeyHash      int16 // key hash (4 bytes): As the name suggests, this field stores a 32-bit hash value calculated from the key.
  DataLength   int32 // data len (4 bytes): The length of the data stored in this page.
  RecordLength int32 // record len (4 bytes): the total length of the data record of which part is stored in this page.
  NextPage     int64 // next page (8 bytes): page number of the page that contains subsequent data for the record stored in this page, if more data is available.
}

func newPageHeader(config Config) *pageHeader {
  return &pageHeader{}
}

type page struct {
  pageNumber int64
  offset     int64

  paged      *paged
  pageHeader *pageHeader

  data []byte
}

func newPage(paged *paged, pageNumber int64) *page {
  fileHeader := paged.fileHeader

  page := page{}

  page.paged = paged
  page.pageHeader = newPageHeader(paged.config)
  page.pageNumber = pageNumber
  page.offset = int64(fileHeader.HeaderSize) +
    (int64(pageNumber) * int64(fileHeader.PageSize))

  return &page
}

func (page *page) read() error {
  if len(page.data) > 0 {
    return nil
  }

  pageHeaderSize := page.paged.fileHeader.PageHeaderSize
  pageSize := page.paged.fileHeader.PageSize
  pageDataOffset := page.offset + int64(pageHeaderSize)
  pageDataSize := pageSize - int32(pageHeaderSize)

  var errors error

  if err := read(page.paged, page.offset, uint32(pageHeaderSize), page.pageHeader); err != nil {
    multierror.Append(errors, err)
  }

  if err := read(page.paged, pageDataOffset, uint32(pageDataSize), page.data); err != nil {
    multierror.Append(errors, err)
  }

  return errors
}

func (page *page) write() error {
  dataOffset := int64(page.offset) + int64(page.paged.fileHeader.PageHeaderSize)

  var errors error

  if err := write(page.paged, int64(page.offset), page.pageHeader); err != nil {
    multierror.Append(errors, err)
  }

  if err := write(page.paged, dataOffset, &page.data); err != nil {
    multierror.Append(errors, err)
  }

  return errors
}