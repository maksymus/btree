package main

// https://cstack.github.io/db_tutorial/

func main() {
  filename :=  "/tmp/btree1.dat"

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

  conf := DefaultConfig()

  paged := newPaged(filename, conf)

  paged.open()
  defer paged.close()

  // p, err := paged.getPage(0)
  // fmt.Println(p, err)
}
