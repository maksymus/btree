package main

import (
  "bytes"
  "encoding/binary"
  "fmt"
  sysio "io"
  "os"
  "sync"

  "btree/errors"
  "btree/io"
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

  *fileHeader

  lock sync.RWMutex
}

// create new paged
func newPaged(filename string, config Config) *paged {
  return &paged{
    filename:   filename,
    fileHeader: newFileHeader(config),
  }
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
  bs := make([]byte, paged.getHeaderSize())
  if err := read(paged, 0, uint32(paged.getHeaderSize()), bs); err != nil {
    return err
  }

  fh := fileHeader{}

  var errs *errors.Error
  bis := io.NewByteInputStream(bs, binary.BigEndian)
  errs = errors.Append(errs, bis.Read(&fh.headerSize))
  errs = errors.Append(errs, bis.Read(&fh.pageSize))
  errs = errors.Append(errs, bis.Read(&fh.pageCount))
  errs = errors.Append(errs, bis.Read(&fh.totalCount))
  errs = errors.Append(errs, bis.Read(&fh.firstFreePage))
  errs = errors.Append(errs, bis.Read(&fh.lastFreePage))
  errs = errors.Append(errs, bis.Read(&fh.pageHeaderSize))
  errs = errors.Append(errs, bis.Read(&fh.maxKeySize))
  errs = errors.Append(errs, bis.Read(&fh.recordCount))

  if errs.ErrorOrNil() == nil {
    log.Debug("read file header from paged file: ", fmt.Sprintf("%+v", fh))

    paged.fileHeader = &fh
    paged.fileHeader.dirty = true
  }

  return errs.ErrorOrNil()
}

// write paged file header
func (paged *paged) write() error {
  var errs *errors.Error

  if paged.dirty {
    bos := io.NewByteOutputStream(binary.BigEndian)
    errs = errors.Append(errs, bos.Write(paged.getHeaderSize()))
    errs = errors.Append(errs, bos.Write(paged.getPageSize()))
    errs = errors.Append(errs, bos.Write(paged.getPageCount()))
    errs = errors.Append(errs, bos.Write(paged.getTotalCount()))
    errs = errors.Append(errs, bos.Write(paged.getFirstFreePage()))
    errs = errors.Append(errs, bos.Write(paged.getLastFreePage()))
    errs = errors.Append(errs, bos.Write(paged.getPageHeaderSize()))
    errs = errors.Append(errs, bos.Write(paged.getMaxKeySize()))
    errs = errors.Append(errs, bos.Write(paged.getRecordCount()))

    if errs.ErrorOrNil() == nil {
      return write(paged, 0, bos.Bytes())
    }
  }

  return errs.ErrorOrNil()
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
  recordLength := page.recordLength
  buffer := bytes.NewBuffer(make([]byte, recordLength))

  currentPage := page

  for {
    if err := page.streamTo(buffer); err != nil {
      return nil, err
    }

    nextPageNum := currentPage.nextPage
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
  page.setRecordLength(int32(value.len))

  buffer := value.Buffer()

  // if more data left in buffer then write to page
  for buffer.Len() > 0 {
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

    if nextPageNum := page.nextPage; nextPageNum == NoPage {
      if freePage, err := paged.getFreePage(); err == nil {
        page = freePage
      } else {
        return err
      }
    } else {
      if nextPage, err := paged.getPage(nextPageNum); err == nil {
        page = nextPage
      } else {
        return err
      }
    }
  }

  // clean up unused overflow pages
  if page.nextPage != NoPage {
    if nextPage, err := paged.getPage(page.nextPage); err == nil {
      paged.unlinkPages(nextPage)
    } else {
      return err
    }
  }

  page.setNextPage(NoPage)

  return nil
}

func (paged *paged) getFreePage() (*page, error) {
  paged.lock.Lock()
  defer paged.lock.Unlock()

  var freePage *page
  var err error

  if paged.getFirstFreePage() != NoPage {
    if freePage, err = paged.getPage(paged.getFirstFreePage()); err == nil {
      paged.setFirstFreePage(freePage.nextPage)
      if freePage.nextPage == NoPage {
        paged.setLastFreePage(NoPage)
      }
    }
  } else {
    freePage, err = paged.getPage(paged.incrementTotalCount())
  }

  if err != nil {
    return nil, err
  }

  freePage.setNextPage(NoPage)
  freePage.setStatus(0) // unused

  return freePage, nil
}

// add overflow pages to list of free pages
func (paged *paged) unlinkPages(p* page) error {
  first, last := p, p

  for ;last.nextPage != NoPage; {
    if next, err := paged.getPage(last.nextPage); err == nil {
      last = next
    } else {
      return err
    }
  }
  
  if paged.getLastFreePage() != NoPage {
    if lastPage, err := paged.getPage(paged.getLastFreePage()); err == nil {
      lastPage.setNextPage(p.pageNumber)
    } else {
      return err 
    }
  } 
  
  if paged.getFirstFreePage() == NoPage {
    paged.setFirstFreePage(first.pageNumber)
  }
  
  paged.setLastFreePage(last.pageNumber)

  return nil
}

// The paged file header consists of a number of fixed-length fields. Fields which are longer than one byte, are always
// stored in Big Endian format, which means the most significant byte is written at the lowest address.
type fileHeader struct {
  headerSize     int16 // header size (2 bytes): set to 4096 (0x1000), the size of this header.
  pageSize       int32 // page size (4 bytes): set to 4096 (0x00001000), the page size.
  pageCount      int64 // page count (8 bytes): this field is not used consistently. It is present mainly for historical reason.
  totalCount     int64 // total page count (8 bytes): total number of pages present in this file.
  firstFreePage  int64 // first free page (8 bytes): page number of the first unused page in this file.
  lastFreePage   int64 // last free page (8 bytes): page number of the last unused page in this file.
  pageHeaderSize int8  // page header size (1 byte): size of each page header. Set to 64 (0x40) by default.
  maxKeySize     int16
  recordCount    int64 // record count (8 bytes): number of records stored in this file.

  dirty bool // non transient
}

func newFileHeader(config Config) *fileHeader {
  return &fileHeader{
    headerSize:     config.headerSize,
    pageSize:       config.pageSize,
    pageCount:      config.pageCount,
    pageHeaderSize: config.pageHeaderSize,
    maxKeySize:     config.maxKeySize,
    dirty:          true,
  }
}

func (fh *fileHeader) getHeaderSize() int16 {
  return fh.headerSize
}

func (fh *fileHeader) getPageSize() int32 {
  return fh.pageSize
}

func (fh *fileHeader) getPageCount() int64 {
  return fh.pageCount
}

func (fh *fileHeader) getTotalCount() int64 {
  return fh.totalCount
}

func (fh *fileHeader) getPageHeaderSize() int8 {
  return fh.pageHeaderSize
}

func (fh *fileHeader) getMaxKeySize() int16 {
  return fh.maxKeySize
}

func (fh *fileHeader) incrementTotalCount() int64 {
  fh.totalCount++
  fh.dirty = true

  return fh.totalCount - 1
}

func (fh *fileHeader) getFirstFreePage() int64 {
  return fh.firstFreePage
}

func (fh *fileHeader) setFirstFreePage(pageNumber int64) {
  fh.dirty = true
  fh.firstFreePage = pageNumber
}

func (fh *fileHeader) getLastFreePage() int64 {
  return fh.lastFreePage
}

func (fh *fileHeader) setLastFreePage(pageNumber int64) {
  fh.dirty = true
  fh.lastFreePage = pageNumber
}

func (fh *fileHeader) getRecordCount() int64 {
  return fh.recordCount
}

func (fh *fileHeader) setRecordCount(count int64) {
  fh.dirty = true
  fh.recordCount = count
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
    if err != sysio.EOF {
      return errors.Wrap(err)
    }
  }

  if err := binary.Read(bytes.NewBuffer(bs), binary.BigEndian, data); err != nil {
    return errors.Wrap(err)
  }

  return nil
}

