package main

import (
	"bytes"
	"encoding/binary"
)

/**
 * Internal Cell Layout
 * +---------------+--------------+------------------+
 * | Key Size      | Page         | Key              |
 * +---------------+--------------+------------------+
 *
 * Leaf Cell Layout
 * +---------------+--------------+-------------+----------------------+
 * | Key Size      | Data Size    | Key         | Data                 |
 * +---------------+--------------+-------------+----------------------+
 */

type (
	Cell interface {
		Size() uint
		Write(bs []byte)
		Read(bs []byte)
		Compare(bs []byte) int

		GetPointer() CellPointer
		SetPointer(CellPointer)

		GetPage() uint32
		GetData() []byte
	}

	CellPointer struct {
		offset uint16
		length uint16
	}

	BTreeCell struct {
		CellPointer
		keySize uint32
		key     []byte
	}

	InternalCell struct {
		BTreeCell
		// pointer to a child page
		page uint32
	}

	LeafCell struct {
		BTreeCell
		dataSize uint32
		data     []byte
	}
)

func NewInternalCell(key []byte, page uint32) *InternalCell {
	return &InternalCell{
		BTreeCell: BTreeCell{
			keySize: uint32(len(key)),
			key:     key,
		},
		page: page,
	}
}

func NewLeafCell(key []byte, data []byte) *LeafCell {
	return &LeafCell{
		BTreeCell: BTreeCell{
			keySize: uint32(len(key)),
			key:     key,
		},
		dataSize: uint32(len(data)),
		data:     data,
	}
}

func (c *BTreeCell) GetPointer() CellPointer {
	return c.CellPointer
}

func (c *BTreeCell) SetPointer(cp CellPointer) {
	c.CellPointer = cp
}

func (c *BTreeCell) Compare(bs []byte) int {
	return bytes.Compare(c.key, bs)
}

// Size returns the size of the cell - 4 bytes for key size + 4 bytes for page + key
func (c *InternalCell) Size() uint {
	return uint(8 + c.keySize)
}

func (c *InternalCell) Write(bs []byte) {
	binary.BigEndian.PutUint32(bs[0:4], c.keySize)
	binary.BigEndian.PutUint32(bs[4:8], c.page)
	copy(bs[8:], c.key)
}

func (c *InternalCell) Read(bs []byte) {
	c.keySize = binary.BigEndian.Uint32(bs[0:4])
	c.page = binary.BigEndian.Uint32(bs[4:8])
	c.key = bs[8 : 8+c.keySize]
}

func (c *InternalCell) GetPage() uint32 {
	return c.page
}

func (c *InternalCell) GetData() []byte {
	return nil
}

// Size returns the size of the cell - 4 bytes for key size + 4 bytes for data size + key + data
func (c *LeafCell) Size() uint {
	return uint(8 + c.keySize + c.dataSize)
}

func (c *LeafCell) Write(bs []byte) {
	binary.BigEndian.PutUint32(bs[0:4], c.keySize)
	binary.BigEndian.PutUint32(bs[4:8], c.dataSize)
	copy(bs[8:], c.key)
	copy(bs[8+c.keySize:], c.data)
}

func (c *LeafCell) Read(bs []byte) {
	c.keySize = binary.BigEndian.Uint32(bs[0:4])
	c.dataSize = binary.BigEndian.Uint32(bs[4:8])
	c.key = bs[8 : 8+c.keySize]
	c.data = bs[8+c.keySize : 8+c.keySize+c.dataSize]
}

func (c *LeafCell) GetPage() uint32 {
	return 0
}

func (c *LeafCell) GetData() []byte {
	return c.data
}
