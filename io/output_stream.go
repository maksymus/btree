package io

import (
  "bytes"
  "encoding/binary"
  "sync"
)

type OutputStream interface {
  Write(data interface{}) error
}

type byteArrayOutputStream struct {
  buffer *bytes.Buffer
  order  binary.ByteOrder

  lock sync.RWMutex
}

func NewByteOutputStream(order binary.ByteOrder) *byteArrayOutputStream {
  return &byteArrayOutputStream{
    buffer: new(bytes.Buffer),
    order:  order,
  }
}

func (stream *byteArrayOutputStream) Write(data interface{}) error {
  stream.lock.Lock()
  defer stream.lock.Unlock()

  return binary.Write(stream.buffer, stream.order, data)
}

func (stream *byteArrayOutputStream) Bytes() []byte {
  stream.lock.RLock()
  defer stream.lock.RUnlock()

  return stream.buffer.Bytes()
}
