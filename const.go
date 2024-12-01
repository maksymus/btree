package main

const (
	FileHeaderSize  = 100
	PageHeaderSize  = 12
	CellPointerSize = 4
	PageSize        = 4096
)

const (
	LeafCellFlag = 1 << 0
)
