package main

import (
  "fmt"
  "os"
)

// https://cstack.github.io/db_tutorial/

func main() {
  paged, err := newPaged("/tmp/btree.dat", DefaultConfig())
  if err != nil {
    fmt.Println("error: ", err)
    os.Exit(1)
  }
  paged.open()
  defer paged.close()

}
