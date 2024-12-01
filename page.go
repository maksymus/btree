package main

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
	panic("TODO")
}

func NewLeafPage() *LeafPage {
	panic("TODO")
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

func (B BTreePage[T]) read(bs []byte) {
	//TODO implement me
	panic("implement me")
}

func (B BTreePage[T]) write(bs []byte) {
	panic("implement me")
}

func read(bs []byte) Page {
	panic("TODO")
}

func write(p Page, bs []byte) {
	p.write(bs)
}
