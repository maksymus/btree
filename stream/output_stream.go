package stream

import (
  "bytes"
  "encoding/binary"
  "sync"
)

type OutputStream interface {
  Write(data interface{}) error
}

type byteOutputStream struct {
  buffer *bytes.Buffer
  order  binary.ByteOrder

  lock sync.RWMutex
}

func NewByteOutputStream(order binary.ByteOrder) *byteOutputStream {
  return &byteOutputStream{
    buffer: new(bytes.Buffer),
    order:  order,
  }
}

func (stream *byteOutputStream) Write(data interface{}) error {
  stream.lock.Lock()
  defer stream.lock.Unlock()

  return binary.Write(stream.buffer, stream.order, data)
}

func (stream *byteOutputStream) Bytes() []byte {
  stream.lock.RLock()
  defer stream.lock.RUnlock()

  return stream.buffer.Bytes()
}
