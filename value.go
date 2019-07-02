package main

import "hash/fnv"

type Value struct {
  data []byte
  pos  uint
  len  uint
  hash int
}

func NewValue(bs []byte, pos uint, len uint) *Value {
  hash32 := fnv.New32()
  hash32.Write(bs)

  return &Value{
    data: bs,
    pos: pos,
    len: len,
    hash: hash32.Size(),
  }
}

func (value Value) ByteAt(pos uint) byte {
  return (value.data)[pos]
}