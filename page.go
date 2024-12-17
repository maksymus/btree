package main

import (
	"encoding/binary"
	"sort"
)

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
	// dummy cell with left page
	if p.cells[0].Compare(key) > 0 {
		return NewInternalCell(nil, p.left)
	}

	start, end := 0, len(p.cells)-1

	for start < end {
		mid := start + (end-start)/2
		cell := p.cells[mid]

		if cell.Compare(key) == 0 {
			return cell
		}

		if cell.Compare(key) < 0 {
			start = mid + 1
		} else {
			end = mid
		}
	}

	return p.cells[end]
}

func (p *LeafPage) Find(key []byte) Cell {
	start, end := 0, len(p.cells)-1

	for start <= end {
		mid := start + (end-start)/2
		cell := p.cells[mid]

		if cell.Compare(key) == 0 {
			return cell
		}

		if cell.Compare(key) < 0 {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}

	return nil
}

func (p *InternalPage) Insert(key []byte, childPage uint32) {
	cell := NewInternalCell(key, childPage)

	neededSize := cell.Size() + CellPointerSize
	if neededSize > uint(p.freeSpace()) {
		// TODO handle split
	}

	cell.CellPointer = CellPointer{p.upper - uint16(cell.Size()), uint16(cell.Size())}

	p.lower += CellPointerSize
	p.upper -= uint16(cell.Size())

	p.cells = append(p.cells, cell)

	sort.Slice(p.cells, func(i, j int) bool {
		return p.cells[i].Compare(p.cells[j].key) < 0
	})

	p.numCells++
}

func (p *LeafPage) Insert(key []byte, data []byte) {
	cell := NewLeafCell(key, data)

	neededSize := cell.Size() + CellPointerSize
	if neededSize > uint(p.freeSpace()) {
		// TODO handle split
	}

	cell.CellPointer = CellPointer{p.upper - uint16(cell.Size()), uint16(cell.Size())}

	p.lower += CellPointerSize
	p.upper -= uint16(cell.Size())

	p.cells = append(p.cells, cell)

	sort.Slice(p.cells, func(i, j int) bool {
		return p.cells[i].Compare(p.cells[j].key) < 0
	})

	p.numCells++
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

func (p PageHeader) freeSpace() uint16 {
	return p.upper - p.lower
}
