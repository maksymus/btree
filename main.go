package main

import "fmt"

func main() {
	finename := "btree.dat"

	bt := NewBTree(finename)

	bt.Open()
	defer bt.Close()

	find, _ := bt.Find([]byte{1})
	fmt.Println(string(find))

	bt.Insert([]byte{1}, []byte("hello world"))
	find, _ = bt.Find([]byte{1})
	fmt.Println(string(find))

	bt.Insert([]byte{3}, []byte("this is last string"))
	find, _ = bt.Find([]byte{3})
	fmt.Println(string(find))

	bt.Insert([]byte{2}, []byte("i'm in the middle"))
	find, _ = bt.Find([]byte{2})
	fmt.Println(string(find))
}
