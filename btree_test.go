package main


import (
  "fmt"
  . "github.com/smartystreets/goconvey/convey"
  "reflect"
  "testing"
)

func Test_BTree_Empty(t *testing.T) {
  Convey("Given btree", t, func() {
    btree := NewBTree(10)

    Convey("Initial btree should be empty", func() {
      empty := btree.Empty()
      So(empty, ShouldEqual, true)
    })
  })
}

func Test_BTree_Insert_OneNode_1Elem(t *testing.T) {
  Convey("Given btree", t, func() {
    btree := NewBTree(10)

    Convey("When elements are inserted", func() {

      btree.Insert(10)

      Convey("Initial btree should not be empty", func() {
        empty := btree.Empty()

        root := btree.root

        So(empty, ShouldEqual, false)
        So(root, ShouldNotBeNil)
        So(root.isLeaf, ShouldEqual, true)

        for i, value := range []int {10} {
          So(root.keys[i], ShouldEqual, value)
        }
      })
    })
  })
}

func Test_BTree_Insert_OneNode_KeysShouldBeOrdered(t *testing.T) {
  Convey("Given btree", t, func() {
    btree := NewBTree(10)

    Convey("When elements are inserted", func() {

      btree.Insert(14)
      btree.Insert(11)
      btree.Insert(334)
      btree.Insert(782)
      btree.Insert(-643)
      btree.Insert(-127)
      btree.Insert(-252)
      btree.Insert(-850)
      btree.Insert(561)
      btree.Insert(145)


      Convey("Initial btree should not be empty", func() {
        empty := btree.Empty()

        root := btree.root

        So(empty, ShouldEqual, false)
        So(root, ShouldNotBeNil)
        So(root.isLeaf, ShouldEqual, true)

        for i, value := range []int {-850, -643, -252, -127, 11, 14, 145, 334, 561, 782} {
          So(root.keys[i], ShouldEqual, value)
        }
      })
    })
  })
}

func Test_BTree_Insert_OneNode_AcceptStructs(t *testing.T) {
  Convey("Given btree", t, func() {
    btree := NewBTree(10)

    Convey("When elements are inserted", func() {

      btree.Insert(&TestKey{2})
      btree.Insert(&TestKey{1})
      btree.Insert(&TestKey{3})

      Convey("Initial btree should not be empty", func() {
        root := btree.root

        So(root, ShouldNotBeNil)
        So(root.isLeaf, ShouldEqual, true)

        for i, value := range []int {1, 2, 3} {
          So(root.keys[i], ShouldResemble, &TestKey{value})
        }
      })
    })
  })
}

func Test_BTree_Insert_SplitRoot(t *testing.T)  {
  Convey("Given btree with 3 keys in root", t, func() {
    btree := NewBTree(2)
    btree.Insert(1)
    btree.Insert(4)
    btree.Insert(3)

    Convey("When element is inserted", func() {

      btree.Insert(2)

      Convey("Root is split and two children are created", func() {

        root := btree.root

        So(root, ShouldNotBeNil)
        So(root.isLeaf, ShouldEqual, false)
        So(len(root.keys), ShouldEqual, 1)
        So(len(root.children), ShouldEqual, 2)

        So(root.keys[0], ShouldEqual, 3)
        So(len(root.children[0].keys), ShouldEqual, 2)
        So(len(root.children[1].keys), ShouldEqual, 1)

        So(root.children[0].keys[0], ShouldEqual, 1)
        So(root.children[0].keys[1], ShouldEqual, 2)
        So(root.children[1].keys[0], ShouldEqual, 4)
      })
    })
  })

}


type TestKey struct {
  i int
}

func (tk TestKey) Compare(other interface{}) int {
  if o, ok := other.(*TestKey); ok {
    if tk.i < o.i {
      return -1
    }

    if tk.i > o.i {
      return 1
    }

    return 0
  }

  panic(fmt.Sprintf("cannot compare TestKey to %v", reflect.TypeOf(other)))
}