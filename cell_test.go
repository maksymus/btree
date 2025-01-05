package main

import (
	"bytes"
	"reflect"
	"testing"
)

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
				bs: []byte{0, 0, 0, 3, 0, 0, 0, 1, 'k', 'e', 'y'},
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
