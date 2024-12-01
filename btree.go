package main

import "os"

/**
 * B+ Tree disk based implementation
 */

type PageNum struct {
	num  uint32
	page Page
}

type BTreeHeader struct {
}

type BTree struct {
	BTreeHeader
	file     *os.File
	filename string
}

func NewBTree(filename string) *BTree {
	return &BTree{filename: filename}
}

// Open file, read header and root page
func (b *BTree) Open() {
	var file *os.File

	// create file if not exists
	if _, err := os.Stat(b.filename); os.IsNotExist(err) {
		file, err = os.Create(b.filename)
		if err != nil {
			panic(err)
		}

		// create 10M file
		file.Truncate(int64(1024 * 102 * 10))

		// create root page
		b.PutPage(1, CreateLeafPage())
	} else {
		if b.file, err = os.OpenFile(b.filename, os.O_RDWR, 0666); err != nil {
			panic(err)
		}
	}
}

// Close file
func (b *BTree) Close() {
	b.file.Close()
}

func (b *BTree) Find(key []byte) ([]byte, error) {
	// read root page
	page := b.GetPage(1)

	for internal, ok := page.(*InternalPage); ok; {
		if cell := internal.Find(key); cell != nil {
			page = b.GetPage(cell.GetPage())
		}
	}

	// search the key
	if cell := page.Find(key); cell != nil {
		return cell.GetData(), nil
	}

	return nil, nil
}

func (b *BTree) Insert(key []byte, data []byte) error {
	stack := NewStack[PageNum]()

	// read root page
	page := b.GetPage(1)
	stack.Push(PageNum{1, page})

	// find leaf page
	for internal, ok := page.(*InternalPage); ok; {
		if cell := internal.Find(key); cell != nil {
			page = b.GetPage(cell.GetPage())
			stack.Push(PageNum{cell.GetPage(), page})
		}
	}

	// search the key
	if leaf, ok := page.(*LeafPage); ok {
		if cell := leaf.Find(key); cell == nil {
			leaf.Insert(key, data)
		} else {
			// TODO key exists
		}
	}

	// TODO fix me
	b.PutPage(1, page)

	return nil
}

func (b *BTree) GetPage(pageNum uint32) Page {
	bs, _ := b.readPage(pageNum)
	return read(bs)
}

func (b *BTree) PutPage(pageNum uint32, page Page) error {
	bs := make([]byte, PageSize)
	write(page, bs)
	return b.writePage(pageNum, bs)
}

func (b *BTree) readPage(pageNum uint32) ([]byte, error) {
	bs := make([]byte, PageSize)
	offset := FileHeaderSize + PageSize*(pageNum-1)
	_, err := b.file.ReadAt(bs, int64(offset))
	return nil, err
}

func (b *BTree) writePage(pageNum uint32, bs []byte) error {
	offset := FileHeaderSize + PageSize*(pageNum-1)
	_, err := b.file.WriteAt(bs, int64(offset))
	return err
}
