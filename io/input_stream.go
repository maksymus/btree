package io

import (
  "bytes"
  "encoding/binary"
  "sync"
)

type InputStream interface {
  Read(data interface{}) error
}

type byteArrayInputStream struct {
  buffer *bytes.Buffer
  order  binary.ByteOrder

  lock sync.Mutex
}

func NewByteInputStream(bs []byte, order binary.ByteOrder) *byteArrayInputStream {
  return &byteArrayInputStream{
    buffer: bytes.NewBuffer(bs),
    order:  order,
  }
}

func (bis *byteArrayInputStream) Read(data interface{}) error {
  bis.lock.Lock()
  defer bis.lock.Unlock()

  return binary.Read(bis.buffer, bis.order, data)
}

