package main

import "encoding/binary"

/**
 * Page
 */

type (
	Page interface {
		PageReaderWriter
		Find(key []byte) Cell
	}

	PageReader interface {
		read(bs []byte)
	}

	PageWriter interface {
		write(bs []byte)
	}

	PageReaderWriter interface {
		PageReader
		PageWriter
	}

	PageHeader struct {
		flags    uint16
		lower    uint16
		upper    uint16
		numCells uint16
		left     uint32
	}

	BTreePage[T Cell] struct {
		PageHeader
		cells   []T
		cellGen func(pointer CellPointer) T
	}

	InternalPage struct {
		BTreePage[*InternalCell]
	}

	LeafPage struct {
		BTreePage[*LeafCell]
	}
)

func CreateInternalPage() *InternalPage {
	page := NewInternalPage()
	page.lower = PageHeaderSize
	page.upper = PageSize

	return page
}

func CreateLeafPage() *LeafPage {
	page := NewLeafPage()

	page.flags |= LeafCellFlag
	page.lower = PageHeaderSize
	page.upper = PageSize

	return page
}

func NewInternalPage() *InternalPage {
	return &InternalPage{
		BTreePage[*InternalCell]{
			cellGen: func(pointer CellPointer) *InternalCell {
				return &InternalCell{
					BTreeCell: BTreeCell{
						CellPointer: pointer,
					},
				}
			},
		},
	}
}

func NewLeafPage() *LeafPage {
	return &LeafPage{
		BTreePage[*LeafCell]{
			cellGen: func(pointer CellPointer) *LeafCell {
				return &LeafCell{
					BTreeCell: BTreeCell{
						CellPointer: pointer,
					},
				}
			},
		},
	}
}

func (p *InternalPage) Find(key []byte) Cell {
	panic("TODO")
}

func (p *LeafPage) Find(key []byte) Cell {
	panic("TODO")
}

func (p *InternalPage) Insert(key []byte, childPage uint32) {
	panic("TODO")
}

func (p *LeafPage) Insert(key []byte, data []byte) {
	panic("TODO")
}

func (p *BTreePage[T]) read(bs []byte) {
	header := PageHeader{}
	header.flags = binary.LittleEndian.Uint16(bs[:2])
	header.lower = binary.LittleEndian.Uint16(bs[2:4])
	header.upper = binary.LittleEndian.Uint16(bs[4:6])
	header.numCells = binary.LittleEndian.Uint16(bs[6:8])
	header.left = binary.LittleEndian.Uint32(bs[8:12])

	p.PageHeader = header
	for i := 0; i < int(header.numCells); i++ {
		pos := PageHeaderSize + (i * CellPointerSize)
		offset := binary.LittleEndian.Uint16(bs[pos+0 : pos+2])
		length := binary.LittleEndian.Uint16(bs[pos+2 : pos+4])

		pointer := CellPointer{offset, length}
		cell := p.cellGen(pointer)
		p.cells = append(p.cells, cell)
	}
}

func (p *BTreePage[T]) write(bs []byte) {
	binary.BigEndian.PutUint16(bs[0:2], p.flags)
	binary.BigEndian.PutUint16(bs[2:4], p.lower)
	binary.BigEndian.PutUint16(bs[4:6], p.upper)
	binary.BigEndian.PutUint16(bs[6:8], p.numCells)
	binary.BigEndian.PutUint32(bs[8:12], p.left)

	for i, cell := range p.cells {
		pos := PageHeaderSize + (i * CellPointerSize)

		pointer := cell.GetCellPointer()
		binary.BigEndian.PutUint16(bs[pos+0:pos+2], pointer.offset)
		binary.BigEndian.PutUint16(bs[pos+2:pos+4], pointer.length)
		cell.Write(bs[pointer.offset : pointer.offset+pointer.length])
	}
}

func read(bs []byte) Page {
	flags := binary.BigEndian.Uint16(bs[:2])

	if flags&LeafCellFlag == 1 {
		page := NewLeafPage()
		page.read(bs)
		return page
	}

	page := NewInternalPage()
	page.read(bs)
	return page
}

func write(p Page, bs []byte) {
	p.write(bs)
}
