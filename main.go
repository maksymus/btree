package main

import (
  "fmt"
  "os"
)

// https://cstack.github.io/db_tutorial/

func main() {
  filename :=  "/tmp/btree.dat"

  // file, _ := os.OpenFile(filename, os.O_RDWR, 0666)
  //
  // bs := make([]byte, 4096)
  // if n, err := file.ReadAt(bs, 0); err != nil {
  //   fmt.Println("test", bs)
  //   // fmt.Printf("err: %s\n%s", err, errors.Stack(err))
  // } else {
  //   fmt.Printf("%n bytes read \n", n)
  // }

  // os.Remove(filename)

  paged, err := newPaged(filename, DefaultConfig())
  if err != nil {
    fmt.Println("error: ", err)
    os.Exit(1)
  }

  paged.open()
  defer paged.close()
}
