package main

import (
  "fmt"
  "reflect"
  "btree/errors"
)

// https://www.geeksforgeeks.org/b-tree-set-1-introduction-2/

/*
1) All leaves are at same level.
2) A B-Tree is defined by the term minimum degree ‘t’. The value of t depends upon disk block size.
3) Every node except root must contain at least t-1 keys. Root may contain minimum 1 key.
4) All nodes (including root) may contain at most 2t – 1 keys.
5) Number of children of a node is equal to the number of keys in it plus 1.
6) All keys of a node are sorted in increasing order. The child between two keys k1 and k2 contains all keys in the range from k1 and k2.
7) B-Tree grows and shrinks from the root which is unlike Binary Search Tree. Binary Search Trees grow downward and also shrink from downward.
8) Like other balanced Binary Search Trees, time complexity to search, insert and delete is O(Logn).
*/

type BTree struct {
  root   *node  // btree root
  degree int   // min number of keys
}

type node struct {
  keys     []interface{} // keys
  children []*node       // pointers to children
  degree   int           // min number of keys
  isLeaf   bool          // is node a leaf
}

type Comparable interface {
  Compare(other interface{}) int
}

func NewBTree(degree uint) *BTree {
  return &BTree{degree: int(degree)}
}

func newNode(degree int, isLeaf bool) (*node, error) {
  if degree < 2 {
    return nil, errors.New("degree should be > 1")
  }

  return &node{
    degree: degree,
    children: make([]*node, 0),
    keys: make([]interface{}, 0),
    isLeaf: isLeaf,
  }, nil
}

func (btree *BTree) Empty() bool {
  return btree.root == nil
}

func (btree *BTree) search(searchKey interface{}) *node {
  btree.checkKeyType(searchKey)

  if btree.root == nil {
    return nil
  }
  return btree.root.search(searchKey)
}

func (btree *BTree) Insert(key interface{})  {
  btree.checkKeyType(key)

  // create root if tree is empty
  if btree.root == nil {
    n, _ := newNode(btree.degree, true)
    n.keys = append(n.keys, key)
    btree.root = n
  } else {
    // if root is full then grow tree in size
    if btree.root.isFull() {
      // create new root
      newRoot, _ := newNode(btree.degree, true)

      // assign root as a child of new root
      newRoot.children = append(newRoot.children, btree.root)

      // split old root - should create two children
      newRoot.splitChild(0, btree.root)

      // insert into first or second node
      if compare(key, newRoot.keys[0]) <= 0 {
        newRoot.children[0].insert(key)
      } else {
        newRoot.children[1].insert(key)
      }

      // reassign root
      btree.root = newRoot
    } else {
      btree.root.insert(key)
    }
  }
}

func (btree *BTree) Delete(key interface{})  {
  btree.checkKeyType(key)

}

func (btree *BTree) checkKeyType(v interface{}) {
  if v == nil {
    panic("nil key value not allowed")
  }

  kind := reflect.TypeOf(v).Kind()

  isInt := int(kind) > 1 && int(kind) < 12
  isFloat := kind == reflect.Float32 || kind == reflect.Float64

  if isInt || isFloat {
    return
  }

  if _, ok := v.(Comparable); ok {
    return
  }

  panic(fmt.Sprintf("incompatible key value %s: int, float numbers and Comparable are supported",
    reflect.TypeOf(v)))
}

func (n *node) search(searchKey interface{}) *node {
  for i, key := range n.keys {
    if compare(searchKey, key) == 0 {
      if !n.isLeaf {
        return n.children[i].search(searchKey)
      }
      return n
    }

    if compare(searchKey, key) < 0 {
      return n.children[i].search(searchKey)
    }
  }

  return nil
}

func (n *node) insert(key interface{}) {
  // find index to insert
  idx := 0
  for i, k := range n.keys {
    if compare(key, k) >= 0 {
      idx = i + 1
    } else {
      break;
    }
  }

  if n.isLeaf {
    prev := key
    n.keys = append(n.keys, nil)
    for i, k := range n.keys {
      if i >= idx {
        n.keys[i] = prev
        prev = k
      }
    }
  } else {
    if n.children[idx].isFull() {
      n.splitChild(idx, n.children[idx])
      idx = idx + 1
    }
    n.children[idx].insert(key)
  }
}

func (n *node) splitChild(idx int, child *node)  {
  if !child.isFull() {
    panic("trying to split non full node")
  }

  if n.isFull() {
    panic("trying to insert element in full parent")
  }

  // not a leaf anymore
  n.isLeaf = false

  // make space for new key and child
  n.keys = append(n.keys, nil)
  n.children = append(n.children, nil)

  // free space in parent node
  for i := len(n.keys) - 1; i > idx ; i-- {
    n.keys[i] = n.keys[i-1]
  }

  for i := len(n.children) - 1; i > idx ; i-- {
    n.children[i] = n.children[i-1]
  }

  // pop up child mid elem to parent
  mid := child.degree - 1
  n.keys[idx] = child.keys[mid]

  // move right keys/children to new child
  otherChild, _ := newNode(n.degree, child.isLeaf)
  otherChild.keys = child.keys[mid + 1:]
  if len(child.children) > 0 {
    otherChild.children = child.children[mid+1:]
  }

  // leave left keys/children in left child
  child.keys = child.keys[:mid]
  if len(child.children) > 0 {
    child.children = child.children[:mid+1]
  }

  n.children[idx+1] = otherChild
}

func (n *node) isFull() bool {
  return len(n.keys) == (2 * n.degree - 1)
}

func compare(i1 interface{}, i2 interface{}) int {
  if reflect.TypeOf(i1) != reflect.TypeOf(i2) {
    panic(fmt.Sprintf("incompatible types [%s, %s] ", reflect.TypeOf(i1), reflect.TypeOf(i2)))
  }

  if _, ok := i1.(int); ok {
    if i1.(int) < i2.(int) {
      return -1
    } else if i1.(int) > i2.(int) {
      return 1
    }

    return 0
  }

  // TODO add its and floats

  if _, ok := i1.(Comparable); ok {
    return i1.(Comparable).Compare(i2)
  }

  panic(fmt.Sprintf("incomparable types %s", reflect.TypeOf(i1)))
}
