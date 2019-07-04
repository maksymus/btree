package main

import (
  "btree/errors"
  "bytes"
  "github.com/hashicorp/go-multierror"
)

const (
  NoPage = -1
)

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

// read page header and page data
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

// write page header and page data
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

// write page data to buffer
func (page *page) streamTo(buffer *bytes.Buffer) error {
  ph := page.pageHeader

  if ph.DataLength > 0 {
    if _, err := buffer.Write(page.data[:ph.DataLength]); err != nil {
      return errors.Wrap(err)
    }
  }

  return nil
}
