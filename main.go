package main

import (
  "fmt"
  "os"
)

// https://cstack.github.io/db_tutorial/

func main() {
  filename :=  "/tmp/btree.dat"

  os.Remove(filename)

  paged, err := newPaged(filename, DefaultConfig())
  if err != nil {
    fmt.Println("error: ", err)
    os.Exit(1)
  }

  paged.open()
  defer paged.close()
}
