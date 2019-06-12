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
func Test_BTree_Delete_Scenarios(t *testing.T) {
  createNode := func(degree int, keys []int32, children []*node) *node {
    node, _ := newNode(degree, len(children) == 0, len(keys), len(children))
    node.children = children

    for i, key := range  keys {
      node.keys[i] = key
    }

    return node
  }

  Convey("Given A-Z btree", t, func() {
    degree := 3

    btree := NewBTree(uint(degree))

    btree.root = createNode(degree, []int32{'P' }, []*node {
      createNode(degree, []int32{'C', 'G', 'M' }, []*node{
        createNode(degree, []int32{'A', 'B' }, nil),
        createNode(degree, []int32{'D', 'E', 'F' }, nil),
        createNode(degree, []int32{'J', 'K', 'L' }, nil),
        createNode(degree, []int32{'N', 'O' }, nil),
      }),
      createNode(degree, []int32{'T', 'X' }, []*node{
        createNode(degree, []int32{'Q', 'R', 'S' }, nil),
        createNode(degree, []int32{'U', 'V' }, nil),
        createNode(degree, []int32{'Y', 'Z' }, nil),
      }),
    })
    checkTreeInvariants(btree)

    // Delete F - delete key from leaf
    Convey("When F is deleted", func() {
      deleted := btree.Delete('F')

      Convey("Validate tree- key deleted from leaf node", func() {
        child1, child2 := btree.root.children[0], btree.root.children[1]
        child11, child12, child13, child14 := child1.children[0], child1.children[1], child1.children[2], child1.children[3]
        child21, child22, child23 := child2.children[0], child2.children[1], child2.children[2]

        So(deleted, ShouldBeTrue)

        validateNodeChar(btree.root, false, 2, []int32{ 'P' })

        validateNodeChar(child1, false, 4, []int32{ 'C', 'G', 'M' })
        validateNodeChar(child2, false, 3, []int32{ 'T', 'X' })

        validateNodeChar(child11, true, 0, []int32{ 'A', 'B' })
        validateNodeChar(child12, true, 0, []int32{ 'D', 'E' })
        validateNodeChar(child13, true, 0, []int32{ 'J', 'K', 'L' })
        validateNodeChar(child14, true, 0, []int32{ 'N', 'O' })
        validateNodeChar(child21, true, 0, []int32{ 'Q', 'R', 'S' })
        validateNodeChar(child22, true, 0, []int32{ 'U', 'V' })
        validateNodeChar(child23, true, 0, []int32{ 'Y', 'Z' })
      })

      checkTreeInvariants(btree)
    })


    // Delete M - delete key from internal node
    Convey("When M is deleted ", func() {
      deleted := btree.Delete('F')
      deleted = btree.Delete('M')

      Convey("Validate tree - key deleted from internal node with enough capacity (no merge/fill)", func() {
        child1, child2 := btree.root.children[0], btree.root.children[1]
        child11, child12, child13, child14 := child1.children[0], child1.children[1], child1.children[2], child1.children[3]
        child21, child22, child23 := child2.children[0], child2.children[1], child2.children[2]

        So(deleted, ShouldBeTrue)

        validateNodeChar(btree.root, false, 2, []int32{ 'P' })

        validateNodeChar(child1, false, 4, []int32{ 'C', 'G', 'L' })
        validateNodeChar(child2, false, 3, []int32{ 'T', 'X' })

        validateNodeChar(child11, true, 0, []int32{ 'A', 'B' })
        validateNodeChar(child12, true, 0, []int32{ 'D', 'E' })
        validateNodeChar(child13, true, 0, []int32{ 'J', 'K' })
        validateNodeChar(child14, true, 0, []int32{ 'N', 'O' })
        validateNodeChar(child21, true, 0, []int32{ 'Q', 'R', 'S' })
        validateNodeChar(child22, true, 0, []int32{ 'U', 'V' })
        validateNodeChar(child23, true, 0, []int32{ 'Y', 'Z' })
      })

      checkTreeInvariants(btree)
    })

    Convey("When G is deleted ", func() {
      deleted := btree.Delete('F')
      deleted = btree.Delete('M')
      deleted = btree.Delete('G')

      Convey("Validate tree - key deleted from internal node with merge", func() {
        child1, child2 := btree.root.children[0], btree.root.children[1]
        child11, child12, child13 := child1.children[0], child1.children[1], child1.children[2]
        child21, child22, child23 := child2.children[0], child2.children[1], child2.children[2]

        So(deleted, ShouldBeTrue)

        validateNodeChar(btree.root, false, 2, []int32{ 'P' })

        validateNodeChar(child1, false, 3, []int32{ 'C', 'L' })
        validateNodeChar(child2, false, 3, []int32{ 'T', 'X' })

        validateNodeChar(child11, true, 0, []int32{ 'A', 'B' })
        validateNodeChar(child12, true, 0, []int32{ 'D', 'E', 'J', 'K' })
        validateNodeChar(child13, true, 0, []int32{ 'N', 'O' })
        validateNodeChar(child21, true, 0, []int32{ 'Q', 'R', 'S' })
        validateNodeChar(child22, true, 0, []int32{ 'U', 'V' })
        validateNodeChar(child23, true, 0, []int32{ 'Y', 'Z' })
      })

      checkTreeInvariants(btree)
    })


    Convey("When D is deleted ", func() {
      deleted := btree.Delete('F')
      deleted = btree.Delete('M')
      deleted = btree.Delete('G')
      deleted = btree.Delete('D')

      Convey("Validate tree - key deleted with merge/shrink", func() {
        children := btree.root.children

        child1, child2, child3, child4, child5, child6 :=
          children[0], children[1], children[2], children[3], children[4], children[5]

        So(deleted, ShouldBeTrue)

        validateNodeChar(btree.root, false, 6, []int32{ 'C', 'L', 'P', 'T', 'X'})

        validateNodeChar(child1, true, 0, []int32{ 'A', 'B' })
        validateNodeChar(child2, true, 0, []int32{ 'E', 'J', 'K' })
        validateNodeChar(child3, true, 0, []int32{ 'N', 'O' })
        validateNodeChar(child4, true, 0, []int32{ 'Q', 'R', 'S' })
        validateNodeChar(child5, true, 0, []int32{ 'U', 'V' })
        validateNodeChar(child6, true, 0, []int32{ 'Y', 'Z' })
      })

      checkTreeInvariants(btree)
    })


    // TODO add more scenarios
  })
}

// Helper functions ===============================================================================================

func validateNode(n *node, isLeaf bool, numChildren int, keys []int)  {
  So(n.isLeaf, ShouldEqual, isLeaf)
  So(len(n.keys), ShouldEqual, len(keys))
  So(len(n.children), ShouldEqual, numChildren)

  for i, key := range keys {
    So(n.keys[i], ShouldEqual, key)
  }
}

func validateNodeChar(n *node, isLeaf bool, numChildren int, keys []int32)  {
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

  // return if nil
  if node == nil {
    return
  }

  // check degree is correct
  if !isRoot {
    So(len(node.keys), ShouldBeGreaterThanOrEqualTo, node.degree-1)
  }
  So(len(node.keys), ShouldBeLessThanOrEqualTo, 2*node.degree-1)

  // check keys are ordered
  if len(node.keys) > 1 {
    for i := 1; i < len(node.keys); i++ {
      So(node.keys[i-1], shouldBeLessThanOrEqualTo, node.keys[i])
    }
  }

  // check node has len(keys)+1 children or leaf
  if len(node.children) == 0 {
    So(node.isLeaf, ShouldBeTrue)
  } else {
    So(node.isLeaf, ShouldBeFalse)
    So(len(node.keys), ShouldEqual, len(node.children)-1)
  }

  // check child keys are grater/equal than previous key and less/equal than next key
  for i, child := range node.children {
    if i > 0 {
      prevKey := node.keys[i-1]
      So(prevKey, shouldBeLessThanOrEqualTo, child.keys[0])
    }

    if i < len(node.children)-1 {
      nextKey := node.keys[i]
      So(child.keys[len(child.keys)-1], shouldBeLessThanOrEqualTo, nextKey)
    }
  }

  // validate each child recursively
  for _, child := range node.children {
    checkNodeInvariants(child, false)
  }
}

// Test structure ==============================================================================================
type TestKey struct {
  i int
}

func (tk TestKey) Compare(other interface{}) int {
  var otherTk *TestKey

  if reflect.ValueOf(other).Kind() == reflect.Ptr {
    if o, ok := other.(*TestKey); ok {
      otherTk = o
    }
  } else {
    if o, ok := other.(TestKey); ok {
      otherTk = &o
    }
  }

  if otherTk != nil {
    if tk.i < otherTk.i {
      return -1
    }

    if tk.i > otherTk.i {
      return 1
    }

    return 0
  }


  panic(fmt.Sprintf("cannot compare TestKey to %v", reflect.TypeOf(other)))
}