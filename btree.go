package main

import (
  "btree/errors"
  "fmt"
  "reflect"
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
  root   *node // btree root
  degree int   // min number of keys
}

type node struct {
  keys     []interface{} // keys
  children []*node       // pointers to children
  degree   int           // min number of keys
  isLeaf   bool          // is node a leaf
}

// Comparable interface
type Comparable interface {
  Compare(other interface{}) int
}

// Create new BTree
func NewBTree(degree uint) *BTree {
  return &BTree{degree: int(degree)}
}

func newNode(degree int, isLeaf bool, numKeys int, numChild int) (*node, error) {
  if degree < 2 {
    return nil, errors.New("degree should be > 1")
  }

  if numChild < 0 {
    numChild = 0
  }

  return &node{
    degree:   degree,
    keys:     make([]interface{}, numKeys),
    children: make([]*node, numChild),
    isLeaf:   isLeaf,
  }, nil
}

// Delete key from tree
// @return true if key found and deleted
func (btree *BTree) Delete(searchKey interface{}) bool {
  btree.checkKeyType(searchKey)

  if btree.Empty() {
    return false
  }

  deleted := btree.root.delete(searchKey)

  // if no keys left in root
  if len(btree.root.keys) == 0 {
    // if not a leaf then replace root by first root child
    // else set root to nil
    if !btree.root.isLeaf {
      btree.root = btree.root.children[0]
    } else {
      btree.root = nil
    }
  }

  return deleted
}

func (btree *BTree) Empty() bool {
  return btree.root == nil
}

func (btree *BTree) Search(searchKey interface{}) *node {
  btree.checkKeyType(searchKey)

  if btree.root == nil {
    return nil
  }
  return btree.root.search(searchKey)
}

func (btree *BTree) Insert(key interface{}) {
  btree.checkKeyType(key)

  // create root if tree is empty
  if btree.root == nil {
    n, _ := newNode(btree.degree, true, 0, 0)
    n.keys = append(n.keys, key)
    btree.root = n
  } else {
    // if root is full then grow tree in size
    if btree.root.isFull() {
      // create new root
      newRoot, _ := newNode(btree.degree, true, 0, 0)

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

// Check key if key type is supported
// Currently int, float and Comparable types are supported
// @panic if key type is not supported
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

// Recursively delete node by key
// @return true if key found and key is deleted
func (n *node) delete(key interface{}) bool {
  // find index to delete
  idx := n.findKey(key)

  // if key is present in this node then remove from node
  // else search in child nodes
  if idx < len(n.keys) && compare(n.keys[idx], key) == 0 {
    if n.isLeaf {
      n.deleteFromLeaf(idx)
    } else {
      n.deleteFromInternal(idx)
    }
  } else {
    // if key not found and it is a leaf node then return
    if n.isLeaf {
      return false
    }

    // num of key in node should be more than (degree - 1)
    // if number of keys is less/equals to degree - 1)
    // then borrow key from left/right neighbour or merge
    if len(n.keys) < n.degree {
      n.fill(idx)
    }

    if idx > len(n.keys) {
      idx--
    }

    n.children[idx].delete(key)
  }

  return true
}

// remove key from nodes key list
func (n *node) deleteFromLeaf(idx int) {
  n.keys = append(n.keys[:idx], n.keys[idx+1:]...)
}

func (n *node) deleteFromInternal(idx int) {
  getPredecessor := func() interface{} {
    curr := n.children[idx]
    // keep moving to the rightmost node until leaf is reached
    for ; !curr.isLeaf ;  {
      curr = curr.children[len(curr.children)-1]
    }
    // return last key
    return curr.keys[len(curr.keys)-1]
  }

  getSuccessor := func() interface{} {
    curr := n.children[idx+1]
    // keep moving to the leftmost node until leaf is reached
    for ; !curr.isLeaf ;  {
      curr = curr.children[0]
    }
    // return first key
    return curr.keys[0]
  }

  if len(n.children[idx].keys) >= n.degree {
    // if child that precedes idx (children[idx]) has at least degree keys
    // then find predecessor of key in children[idx] tree
    // replace key[idx] with predecessor and delete predecessor in children[idx]
    pred := getPredecessor()
    n.keys[idx] = pred
    n.children[idx].delete(pred)
  } else if len(n.children[idx+1].keys) >= n.degree {
    // if child that succeeds idx (children[idx+1]) has at least degree keys
    // then find successor of key in children[idx+1] tree
    // replace key[idx] with predecessor and delete successor in children[idx+1]
    succ := getSuccessor()
    n.keys[idx] = succ
    n.children[idx+1].delete(succ)
  } else {
    n.merge(idx)
    n.children[idx].delete(n.keys[idx])
  }
}

func (n *node) fill(idx int) {
  // last key of sibling goes up to parent node n (this node)
  // idx-1 key from parent node n goes down to child as first node
  // sibling's last child becomes child's first child
  borrowPrev := func() {
    child := n.children[idx]
    sibling := n.children[idx-1]

    // move down current node's idx-1 key to child first key
    child.keys = append([]interface{}{n.keys[idx-1]}, child.keys...)

    // move sibling's last child to child's first child
    if !child.isLeaf {
      child.children = append([]*node{sibling.children[len(sibling.children)-1]}, child.children...)
      sibling.children = sibling.children[:len(sibling.children)-1]
    }

    // move up siblings last key to current node's idx-1
    n.keys[idx-1] = sibling.keys[len(sibling.keys)-1]
    sibling.keys = sibling.keys[:len(sibling.keys)-1]
  }

  // first key of sibling goes up to parent node n (this node)
  // idx+1 key from parent node n goes down to child as last node
  // sibling's first child becomes child's last child
  borrowNext := func() {
    child := n.children[idx]
    sibling := n.children[idx+1]

    // move down current node's idx+1 key to child last key
    child.keys = append(child.keys, n.keys[idx+1])

    // move sibling's fist child to child's last child
    if !child.isLeaf {
      child.children = append(child.children, sibling.children[0])
      sibling.children = sibling.children[1:]
    }

    // moving up siblings first key to current node's idx+1
    n.keys[idx+1] = sibling.keys[0]
    sibling.keys = sibling.keys[1:]
  }

  // ================ LOGIC starts here
  // if left child has more than (degree - 1) then borrow from left
  // else if right child has more than (degree - 1) then borrow from right
  // else merge nodes
  if idx != 0 && len(n.children[idx-1].keys) >= n.degree {
    borrowPrev()
  } else if idx != n.degree && len(n.children[idx+1].keys) >= n.degree {
    borrowNext()
  } else {
    if idx >= n.degree {
      idx--
    }
    n.merge(idx)
  }
}

// merge idx and idx+1 children
// idx+1 child is freed
func (n* node) merge(idx int) {
  child := n.children[idx]
  sibling := n.children[idx+1]

  // move down idx key to child
  child.keys = append(child.keys, child.keys[idx])

  // append siblings key to child
  child.keys = append(child.keys, sibling.keys...)

  // append siblings children to child
  if !child.isLeaf {
    child.children = append(child.children, sibling.children...)
  }

  // remove key/child from current node
  n.keys = append(n.keys[:idx], n.keys[idx+1:]...)
  n.children = append(n.children[:idx+1], n.children[idx+2:]...)
}


// Search key by key
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

// Insert key
func (n *node) insert(key interface{}) {
  // find index to insert
  idx := n.findKey(key)

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
      // split child
      n.splitChild(idx, n.children[idx])
      if compare(key, n.keys[idx]) > 0 {
        idx = idx + 1
      }
    }
    n.children[idx].insert(key)
  }
}

// Split child key during insert if child key is full
func (n *node) splitChild(idx int, child *node) {
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

  // free space in parent node for popped key
  for i := len(n.keys) - 1; i > idx; i-- {
    n.keys[i] = n.keys[i-1]
  }

  // free space in parent node one more child
  for i := len(n.children) - 1; i > idx; i-- {
    n.children[i] = n.children[i-1]
  }

  // pop up child mid elem to parent
  mid := child.degree - 1
  n.keys[idx] = child.keys[mid]

  numKeys := len(child.keys) - mid - 1
  numChild := len(child.children) - mid - 1

  // move right keys/children to new child
  otherChild, _ := newNode(n.degree, child.isLeaf, numKeys, numChild)
  copy(otherChild.keys, child.keys[mid+1:])
  if len(child.children) > 0 {
    copy(otherChild.children, child.children[mid+1:])
  }

  // leave left keys/children in left child
  child.keys = child.keys[:mid]
  if len(child.children) > 0 {
    child.children = child.children[:mid+1]
  }

  // add new child to parent node
  n.children[idx+1] = otherChild
}

// Find matching key in node's keys
// @return index of the key grater or equals to search key
func (n *node) findKey(key interface{}) int {
  idx := 0
  for i, k := range n.keys {
    if compare(key, k) > 0 {
      idx = i + 1
    } else {
      break
    }
  }
  return idx
}

// True if number of keys is "2 * degree - 1"
func (n *node) isFull() bool {
  return len(n.keys) == (2*n.degree - 1)
}

// // True if number of keys is "degree - 1"
// func (n *node) isEmpty() bool {
//   return len(n.keys) == n.degree - 1
// }

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

  if _, ok := i1.(float32); ok {
    if i1.(float32) < i2.(float32) {
      return -1
    } else if i1.(float32) > i2.(float32) {
      return 1
    }

    return 0
  }

  if _, ok := i1.(float64); ok {
    if i1.(float64) < i2.(float64) {
      return -1
    } else if i1.(float64) > i2.(float64) {
      return 1
    }

    return 0
  }

  if _, ok := i1.(Comparable); ok {
    return i1.(Comparable).Compare(i2)
  }

  panic(fmt.Sprintf("incomparable types %s", reflect.TypeOf(i1)))
}
