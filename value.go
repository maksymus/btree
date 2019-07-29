package main

import (
  "bytes"
  "hash/fnv"
)

type Value struct {
  data []byte
  pos  int
  len  int
  hash int
}

// create new value
func NewValue(bs []byte, pos int, len int) *Value {
  hash32 := fnv.New32()
  hash32.Write(bs)

  return &Value{
    data: bs,
    pos: pos,
    len: len,
    hash: hash32.Size(),
  }
}

// create buffer from value
func (value *Value) Buffer() *bytes.Buffer {
  return bytes.NewBuffer(value.data[value.pos : value.pos+value.len])
}