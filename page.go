package main

import (
  "btree/errors"
  "bytes"
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

func newPageHeader() *pageHeader {
  return &pageHeader{}
}

// page stores page info and page data
type page struct {
  pageNumber int64 // page number
  offset     int64 // overall page offset if paged file

  paged      *paged      // reference to paged file
  pageHeader *pageHeader // page header with page info

  data []byte // data stores key and value or value only if key is missing
}

func newPage(paged *paged, pageNumber int64) *page {
  return &page{
    paged:      paged,
    pageNumber: pageNumber,
    pageHeader: newPageHeader(),
    offset:     int64(paged.getHeaderSize()) + (int64(pageNumber) * int64(paged.getPageSize())),
  }
}

// read page header and page data from paged file
func (page *page) read() error {
  if page.data == nil {
    return nil
  }

  pageHeaderSize := page.paged.getPageHeaderSize()
  pageSize := page.paged.getPageSize()
  pageDataOffset := page.offset + int64(pageHeaderSize)
  pageDataSize := pageSize - int32(pageHeaderSize)

  var errs *errors.Error

  err := read(page.paged, page.offset, uint32(pageHeaderSize), page.pageHeader)
  errs = errors.Append(errs, err)

  err = read(page.paged, pageDataOffset, uint32(pageDataSize), page.data)
  errs = errors.Append(errs, err)

  return errs.ErrorOrNil()
}

// write page header and page data to paged file
func (page *page) write() error {
  dataOffset := int64(page.offset) + int64(page.paged.getPageHeaderSize())

  var errs *errors.Error

  err := write(page.paged, int64(page.offset), page.pageHeader)
  errs = errors.Append(errs, err)

  err = write(page.paged, dataOffset, &page.data)
  errs = errors.Append(errs, err)

  return errs.ErrorOrNil()
}

// write page data to buffer
func (page *page) streamTo(buffer *bytes.Buffer) error {
  ph := page.pageHeader

  if ph.DataLength > 0 {
    if _, err := buffer.Write(page.data[ph.KeyLength:]); err != nil {
      return errors.Wrap(err)
    }
  }

  return nil
}

// read data from buffer to page
func (page *page) streamFrom(buffer *bytes.Buffer) error {
  paged := page.paged
  pageHeader := page.pageHeader

  // get key/data size of page
  workSize := paged.getPageSize() - int32(paged.getPageHeaderSize())

  // set data length based on length of data in buffer
  bufferLength  := int32(buffer.Len())
  page.pageHeader.DataLength = workSize - int32(pageHeader.KeyLength)
  if bufferLength < page.pageHeader.DataLength {
    page.pageHeader.DataLength = bufferLength
  }

  // read data from buffer
  if _, err := buffer.Read(page.data[page.pageHeader.KeyLength:]); err != nil {
    return err
  }

  return nil
}

func (page *page) getKey() (*Value, error) {
  panic("implement me")
}

func (page *page) setKey(value *Value) error {
  panic("implement me")
}
