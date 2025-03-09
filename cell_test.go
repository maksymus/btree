package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestInternalCell_Size(t *testing.T) {
	type fields struct {
		cell BTreeCell
		page uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			name: "size of internal cell",
			fields: fields{
				cell: BTreeCell{
					keySize: 4,
					key:     []byte{0, 0, 0, 0xFF},
				},
				page: 132,
			},
			want: 12, // 4 bytes for key size + 4 bytes for page + 4 bytes for key
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &InternalCell{
				BTreeCell: tt.fields.cell,
				page:      tt.fields.page,
			}
			if got := c.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInternalCell_Read(t *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want *InternalCell
	}{
		{
			name: "read cell from bytes",
			args: args{
				bs: []byte{ /*key size*/ 0, 0, 0, 3 /*page*/, 0, 0, 0, 1 /*key*/, 'k', 'e', 'y'},
			},
			want: &InternalCell{
				BTreeCell: BTreeCell{
					CellPointer: CellPointer{},
					keySize:     3, // 3 bytes for "key"
					key:         []byte("key"),
				},
				page: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &InternalCell{}

			c.Read(tt.args.bs)

			if !reflect.DeepEqual(c, tt.want) {
				t.Errorf("Read() = %v, want %v", c, tt.want)
			}
		})
	}
}

func TestInternalCell_Write(t *testing.T) {
	type fields struct {
		BTreeCell BTreeCell
		page      uint32
	}
	type args struct {
		bs []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "write cell to bytes",
			fields: fields{
				BTreeCell: BTreeCell{
					CellPointer: CellPointer{},
					keySize:     11, // 11 bytes for "hello world"
					key:         []byte("hello world"),
				},
				page: 1,
			},
			args: args{
				bs: make([]byte, 19), // 4 bytes for key size + 4 bytes for page + 11 bytes for key
			},
			want: []byte{ /*key size*/ 0, 0, 0, 11 /*page*/, 0, 0, 0, 1 /*key*/, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &InternalCell{
				BTreeCell: tt.fields.BTreeCell,
				page:      tt.fields.page,
			}

			c.Write(tt.args.bs)

			if bytes.Compare(tt.args.bs, tt.want) != 0 {
				t.Errorf("Write() = %v, want %v", tt.args.bs, tt.want)
			}
		})
	}
}

func TestLeafCell_Size(t *testing.T) {
	type fields struct {
		BTreeCell BTreeCell
		dataSize  uint32
		data      []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			name: "size of leaf cell",
			fields: fields{
				BTreeCell: BTreeCell{
					keySize: 5,
					key:     []byte{0, 0, 0, 0, 0xFF},
				},
				dataSize: 6,
				data:     []byte("some data"),
			},
			want: 19, // 4 bytes for key size + 4 bytes for data size + 4 bytes for key + 6 bytes for data
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LeafCell{
				BTreeCell: tt.fields.BTreeCell,
				dataSize:  tt.fields.dataSize,
				data:      tt.fields.data,
			}
			if got := c.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLeafCell_Write(t *testing.T) {
	type fields struct {
		BTreeCell BTreeCell
		dataSize  uint32
		data      []byte
	}
	type args struct {
		bs []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LeafCell{
				BTreeCell: tt.fields.BTreeCell,
				dataSize:  tt.fields.dataSize,
				data:      tt.fields.data,
			}
			c.Write(tt.args.bs)
		})
	}
}

func TestLeafCell_Read(t *testing.T) {
	type fields struct {
		BTreeCell BTreeCell
		dataSize  uint32
		data      []byte
	}
	type args struct {
		bs []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LeafCell{
				BTreeCell: tt.fields.BTreeCell,
				dataSize:  tt.fields.dataSize,
				data:      tt.fields.data,
			}
			c.Read(tt.args.bs)
		})
	}
}
