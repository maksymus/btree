package main

import (
  "btree/errors"
  "btree/io"
  "bytes"
  "encoding/binary"
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
  status       int8   // status (1 byte): pages in the data file are either used or unused. Used pages contain actual data.
  keyLength    int16  // key length (2 bytes): pages have the possibility of storing a key just before their actual data.
  keyHash      uint32 // key hash (4 bytes): As the name suggests, this field stores a 32-bit hash value calculated from the key.
  dataLength   int32  // data len (4 bytes): The length of the data stored in this page.
  recordLength int32  // record len (4 bytes): the total length of the data record of which part is stored in this page.
  nextPage     int64  // next page (8 bytes): page number of the page that contains subsequent data for the record stored in this page, if more data is available.

  dirty bool // transient
}

func (page *page) setStatus(status int8) {
  page.status = status
  page.dirty = true
}

func (page *page) setKeyLength(keyLength int16) {
  page.keyLength = keyLength
  page.dirty = true
}

func (page *page) setKeyHash(keyHash uint32) {
  page.keyHash = keyHash
  page.dirty = true
}

func (page *page) setDataLength(dataLength int32) {
  page.dataLength = dataLength
  page.dirty = true
}

func (page *page) setRecordLength(recordLength int32) {
  page.recordLength = recordLength
  page.dirty = true
}

func (page *page) setNextPage(nextPage int64) {
  page.nextPage = nextPage
  page.dirty = true
}

// page stores page info and page data
type page struct {
  pageNumber int64 // page number
  offset     int64 // overall page offset if paged file

  pageHeader // page header with page info

  data []byte // data stores key and value or value only if key is missing
}

func newPage(paged *paged, pageNumber int64) *page {
  pageHeader := pageHeader{}
  pageHeader.nextPage = NoPage

  return &page{
    pageNumber: pageNumber,
    pageHeader: pageHeader,
    offset:     int64(paged.getHeaderSize()) + (pageNumber * int64(paged.getPageSize())),
  }
}

// read page header and page data from paged file
func (paged *paged) readPage(page *page) error {
  // if page is not in file then return new page
  if !paged.isFilePage(page.pageNumber) {
    return nil
  }

  pageHeaderSize := paged.getPageHeaderSize()
  pageDataOffset := page.offset + int64(pageHeaderSize)

  var errs *errors.Error

  // read page header
  bs := make([]byte, paged.getPageHeaderSize())
  errs = errors.Append(errs, read(paged, page.offset, uint32(pageHeaderSize), bs))

  bis := io.NewByteInputStream(bs, binary.BigEndian)
  errs = errors.Append(errs, bis.Read(&page.status))
  errs = errors.Append(errs, bis.Read(&page.keyLength))
  errs = errors.Append(errs, bis.Read(&page.keyHash))
  errs = errors.Append(errs, bis.Read(&page.dataLength))
  errs = errors.Append(errs, bis.Read(&page.recordLength))
  errs = errors.Append(errs, bis.Read(&page.nextPage))

  // read page data
  page.data = make([]byte, page.dataLength)
  errs = errors.Append(errs, read(paged, pageDataOffset, uint32(page.dataLength), &page.data))

  return errs.ErrorOrNil()
}

// write page header and page data to paged file
func (paged *paged) writePage(page *page) error {
  var errs *errors.Error

  if page.dirty {
    headerBos := io.NewByteOutputStream(binary.BigEndian)
    errs = errors.Append(errs, headerBos.Write(page.status))
    errs = errors.Append(errs, headerBos.Write(page.keyLength))
    errs = errors.Append(errs, headerBos.Write(page.keyHash))
    errs = errors.Append(errs, headerBos.Write(page.dataLength))
    errs = errors.Append(errs, headerBos.Write(page.recordLength))
    errs = errors.Append(errs, headerBos.Write(page.nextPage))

    dataOffset := int64(page.offset) + int64(paged.getPageHeaderSize())

    if errs.ErrorOrNil() == nil {
      errs = errors.Append(errs, write(paged, int64(page.offset), headerBos.Bytes()))
      errs = errors.Append(errs, write(paged, dataOffset, &page.data))
    }
  }

  return errs.ErrorOrNil()
}

// write page data to buffer
func (page *page) streamTo(buffer *bytes.Buffer) error {
  if page.dataLength > 0 {
    if _, err := buffer.Write(page.data[page.keyLength:]); err != nil {
      return errors.Wrap(err)
    }
  }

  return nil
}

// read data from buffer to page
func (page *page) streamFrom(buffer *bytes.Buffer, paged *paged) error {
  // get key/data size of page
  workSize := paged.getPageSize() - int32(paged.getPageHeaderSize())

  // set data length based on length of data in buffer
  bufferLength  := int32(buffer.Len())
  page.dataLength = workSize - int32(page.keyLength)
  if bufferLength < page.dataLength {
    page.dataLength = bufferLength
  }

  data := make([]byte, page.dataLength)
  page.data = append(page.data, data...)

  // read data from buffer
  if _, err := buffer.Read(page.data[page.keyLength:]); err != nil {
    return errors.Wrap(err)
  }

  page.dirty = true

  return nil
}

func (page *page) getKey() *Value {
  return NewValue(page.data, 0, int(page.keyLength))
}

// setKey wipes out data
func (page *page) setKey(value *Value) {
  page.setKeyLength(int16(len(value.data)))
  page.setKeyHash(value.hash)
  page.setRecordLength(0)

  page.data = make([]byte, page.keyLength)

  copy(page.data[:page.keyLength], value.data)
}
