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

    checkTreeInvariants(btree)
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

    checkTreeInvariants(btree)
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

    checkTreeInvariants(btree)
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

    checkTreeInvariants(btree)
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
        child1 := root.children[0]
        child2 := root.children[1]

        validateNode(root, false, 2, []int { 3 })
        validateNode(child1, true, 0, []int { 1, 2 })
        validateNode(child2, true, 0, []int { 4 })
      })
    })

    checkTreeInvariants(btree)
  })
}

func Test_BTree_Insert_Split(t *testing.T)  {
  Convey("Given btree with keys inserted", t, func() {
    btree := NewBTree(2)
    btree.Insert(1)
    btree.Insert(4)
    btree.Insert(3)
    btree.Insert(7)
    btree.Insert(5)
    btree.Insert(6)
    btree.Insert(8)
    btree.Insert(9)
    btree.Insert(0)
    btree.Insert(2)

    Convey("Nodes should be in order", func() {
      root := btree.root
      child1 := root.children[0]
      child2 := root.children[1]

      child3 := child1.children[0]
      child4 := child1.children[1]
      child5 := child2.children[0]
      child6 := child2.children[1]

      validateNode(root, false, 2, []int { 5 })
      validateNode(child1, false, 2, []int { 3 })
      validateNode(child2, false, 2, []int { 7 })

      validateNode(child3, true, 0, []int { 0, 1, 2 })
      validateNode(child4, true, 0, []int { 4 })
      validateNode(child5, true, 0, []int { 6 })
      validateNode(child6, true, 0, []int { 8, 9 })
    })

    checkTreeInvariants(btree)
  })
}

func Test_BTree_Insert_BigTest(t *testing.T)  {
  Convey("Given btree with keys inserted", t, func() {
    btree := NewBTree(3)
    btree.Insert(1)
    btree.Insert(3)
    btree.Insert(7)
    btree.Insert(10)
    btree.Insert(11)
    btree.Insert(13)
    btree.Insert(14)
    btree.Insert(15)
    btree.Insert(18)
    btree.Insert(16)
    btree.Insert(19)
    btree.Insert(24)
    btree.Insert(25)
    btree.Insert(26)
    btree.Insert(21)
    btree.Insert(4)
    btree.Insert(5)
    btree.Insert(20)
    btree.Insert(22)
    btree.Insert(2)
    btree.Insert(17)
    btree.Insert(12)
    btree.Insert(6)

    Convey("Nodes should be in order", func() {
      // root/first level
      root := btree.root

      // second level
      child1 := root.children[0]
      child2 := root.children[1]

      // third level
      child3 := child1.children[0]
      child4 := child1.children[1]
      child5 := child1.children[2]
      child6 := child1.children[3]
      child7 := child2.children[0]
      child8 := child2.children[1]
      child9 := child2.children[2]

      validateNode(root, false, 2, []int { 16 })
      validateNode(child1, false, 4, []int { 3, 7, 13 })
      validateNode(child2, false, 3, []int { 20, 24 })
      validateNode(child3, true, 0, []int { 1, 2 })
      validateNode(child4, true, 0, []int { 4, 5, 6 })
      validateNode(child5, true, 0, []int { 10, 11, 12 })
      validateNode(child6, true, 0, []int { 14, 15 })
      validateNode(child7, true, 0, []int { 17, 18, 19 })
      validateNode(child8, true, 0, []int { 21, 22 })
      validateNode(child9, true, 0, []int { 25, 26 })
    })

    checkTreeInvariants(btree)
  })
}

func Test_BTree_Delete_RootOneKey(t *testing.T) {
  Convey("Given btree with one key", t, func() {
    btree := NewBTree(2)
    btree.Insert(2)

    Convey("When elements are deleted", func() {
      deleted := btree.Delete(2)

      Convey("Root should be nil", func() {
        root := btree.root
        So(root, ShouldBeNil)
        So(deleted, ShouldBeTrue)
      })
    })
  })
}


func Test_BTree_Delete_RootOneKey_NotFound(t *testing.T) {
  Convey("Given btree with one key", t, func() {
    btree := NewBTree(2)
    btree.Insert(2)

    Convey("When elements are deleted and key not found", func() {
      deleted := btree.Delete(1)

      Convey("Root should be nil", func() {
        validateNode(btree.root, true, 0, []int{2})
        So(deleted, ShouldBeFalse)
      })
    })

    checkTreeInvariants(btree)
  })
}

// https://www.geeksforgeeks.org/b-tree-set-3delete/
// func Test_BTree_Delete_Scenarios(t *testing.T) {
//   Convey("Given btree with one key", t, func() {
//     btree := NewBTree(3)
//
//     Convey("When elements are deleted and key not found", func() {
//       deleted := btree.Delete(1)
//
//       Convey("Root should be nil", func() {
//         validateNode(btree.root, true, 0, []int{2})
//         So(deleted, ShouldBeFalse)
//       })
//     })
//
//     checkTreeInvariants(btree)
//   })
// }

func validateNode(n *node, isLeaf bool, numChildren int, keys []int)  {
  So(n.isLeaf, ShouldEqual, isLeaf)
  So(len(n.keys), ShouldEqual, len(keys))
  So(len(n.children), ShouldEqual, numChildren)

  for i, key := range keys {
    So(n.keys[i], ShouldEqual, key)
  }
}

func checkTreeInvariants(tree *BTree) {
  checkNodeInvariants(tree.root, true)
}

func checkNodeInvariants(node *node, isRoot bool) {
  shouldBeLessThanOrEqualTo := func(actual interface{}, expected ...interface{}) string {
    if compare(actual, expected[0]) <= 0 {
      return ""
    }

    return fmt.Sprintf("Expected '%v' to be less than or equal to '%v' (but it wasn't)!", actual, expected[0])
  }

  // return igf nil
  if node == nil { return }

  // check degree is correct
  if !isRoot {
    So(len(node.keys), ShouldBeGreaterThanOrEqualTo, node.degree - 1)
  }
  So(len(node.keys), ShouldBeLessThanOrEqualTo, 2 * node.degree - 1)

  // check node has len(keys)+1 children or leaf
  if len(node.children) == 0 {
    So(node.isLeaf, ShouldBeTrue)
  } else {
    So(node.isLeaf, ShouldBeFalse)
    So(len(node.keys), ShouldEqual, len(node.children) - 1)
  }

  if len(node.keys) > 1 {
    // check keys are ordered
    for i := 1; i < len(node.keys); i++ {
      So(node.keys[i-1], shouldBeLessThanOrEqualTo, node.keys[i])
    }
  }
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