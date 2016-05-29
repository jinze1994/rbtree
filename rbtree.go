package rbtree

/*
A pure Golang implementation of a red-black tree as described
by Thomas H in introduction to algorithms 3rd edition

Red-Black tree properties:  http://en.wikipedia.org/wiki/Rbtree

 1) A node is either red or black
 2) The root is black
 3) All leaves (NULL) are black
 4) Both children of every red node are black
 5) Every simple path from root to leaves contains the same number
    of black nodes.
*/

import (
	"bytes"
	"fmt"
)

/*
	Rbnode represents a Red-Black tree's single node
*/
type Rbnode struct {
	left   *Rbnode
	right  *Rbnode
	parent *Rbnode

	prev *Rbnode
	next *Rbnode

	/*
		A interface contains key and value
		Must implement Less method (like package 'sort' need)
		In red-black tree, the key must be unique
		Less(a, b) == false && Less(b, a) == false means a equal b
		Note: the node's Item is type interface, like the pointer. You shouldn't change the item's key before you delete them in red-black tree.
	*/
	Item

	color bool
}

const (
	RED   = false
	BLACK = true
)

type Item interface {
	Less(than Item) bool
}

/*
	Rbtree represents a Red-Black tree.
*/
type Rbtree struct {
	nill  *Rbnode
	root  *Rbnode
	count int
	first *Rbnode
	last  *Rbnode
}

/*
	New returns an pointer to initialized Red-Black tree
*/
func NewRbtree() *Rbtree {
	nillNode := &Rbnode{nil, nil, nil, nil, nil, nil, BLACK}
	return &Rbtree{
		nill:  nillNode,
		root:  nillNode,
		count: 0,
		first: nillNode,
		last:  nillNode,
	}
}

// ===================== Main API Method ==========================

func (t *Rbtree) Init() {
	nillNode := &Rbnode{nil, nil, nil, nil, nil, nil, BLACK}

	t.nill = nillNode
	t.root = nillNode
	t.count = 0
	t.first = nillNode
	t.last = nillNode
}

/*
	Return curent number of nodes in the tree.
*/
func (t *Rbtree) Count() int {
	l := t.count
	return int(l)
}

/*
	Insert 'item' to red-black tree
	when argument 'item' == nil, return 'nil, false'
	when returned 'ok' == true
		returned 'node' is the inserted node with 'item'
	when returned 'ok' == false
		means there already is one node equally with 'item' by twice Less method comparison
		returned 'node' is that node
*/
func (t *Rbtree) Insert(item Item) (node *Rbnode, ok bool) {
	if item == nil {
		return nil, false
	}

	// Always insert a RED node
	return t.insert(&Rbnode{t.nill, t.nill, t.nill, t.nill, t.nill, item, RED})
}

/*
	Remove the node equally with argument item in red-black tree
	when returned 'ok' == true
		returned 'i' is that Item equally with argument item
	when returned 'ok' == false
		means there isn't any node equally with argument item
		returned 'i' == nil
*/
func (t *Rbtree) Remove(item Item) (i Item, ok bool) {
	if item == nil {
		return nil, false
	}

	// The `color` field here is nobody
	var node *Rbnode
	if node, ok = t.remove(&Rbnode{t.nill, t.nill, t.nill, t.nill, t.nill, item, RED}); ok {
		return node.Item, true
	}
	return nil, false
}

func (t *Rbtree) Remove_raw(z *Rbnode) (i Item, ok bool) {
	if z == nil {
		return nil, false
	}
	return t.remove_raw(z)
}

/*
	Find the node equally with argument item in red-black tree
	Return that node if found, or return nil
	Note: the node's Item is type interface, like the pointer. You shouldn't change the item's key before you delete them in red-black tree.
*/
func (t *Rbtree) Get(item Item) *Rbnode {
	if item == nil {
		return nil
	}

	// The `color` field here is nobody
	ret := t.search(&Rbnode{t.nill, t.nill, t.nill, t.nill, t.nill, item, RED})
	if ret == t.nill {
		return nil
	} else {
		return ret
	}
}

/*
	Get the First rbnode
	Note: this is not a thread-safe mothod
*/
func (t *Rbtree) First() *Rbnode {
	return t.first
}

/*
	Get the Last rbnode
	Note: this is not a thread-safe mothod
*/
func (t *Rbtree) Last() *Rbnode {
	return t.last
}

/*
	Get the Next rbnode
	Note: this is not a thread-safe mothod
*/
func (node *Rbnode) Next() *Rbnode {
	return node.next
}

/*
	Get the Previous rbnode
	Note: this is not a thread-safe mothod
*/
func (node *Rbnode) Prev() *Rbnode {
	return node.prev
}

// ==================== Private Method for Internal Support ===================

func (t *Rbtree) leftRotate(x *Rbnode) {
	// Since we are doing the left rotation, the right child should *NOT* nil.
	if x.right == t.nill {
		return
	}

	//
	// The illation of left rotation
	//
	//          |                                  |
	//          X                                  Y
	//         / \         left rotate            / \
	//        α  Y       ------------->         X   γ
	//           / \                            / \
	//          β  γ                         α  β
	//
	// It should be note that during the rotating we do not change
	// the Rbnodes' color.
	//
	y := x.right
	x.right = y.left
	if y.left != t.nill {
		y.left.parent = x
	}
	y.parent = x.parent

	if x.parent == t.nill {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.left = x
	x.parent = y
}

func (t *Rbtree) rightRotate(x *Rbnode) {
	// Since we are doing the right rotation, the left child should *NOT* nil.
	if x.left == t.nill {
		return
	}

	//
	// The illation of right rotation
	//
	//          |                                  |
	//          X                                  Y
	//         / \         right rotate           / \
	//        Y   γ      ------------->         α  X
	//       / \                                    / \
	//      α  β                                 β  γ
	//
	// It should be note that during the rotating we do not change
	// the Rbnodes' color.
	//
	y := x.left
	x.left = y.right
	if y.right != t.nill {
		y.right.parent = x
	}
	y.parent = x.parent

	if x.parent == t.nill {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.right = x
	x.parent = y
}

func (t *Rbtree) insert(z *Rbnode) (*Rbnode, bool) {
	x := t.root
	y := t.nill

	for x != t.nill {
		y = x
		if z.Item.Less(x.Item) {
			x = x.left
		} else if x.Item.Less(z.Item) {
			x = x.right
		} else {
			return x, false
		}
	}

	z.parent = y
	if y == t.nill {
		t.root = z
	} else if z.Item.Less(y.Item) {
		y.left = z
	} else {
		y.right = z
	}

	t.count++
	t.insertFixup(z)

	z.next = t.successor(z)
	if z.next != t.nill {
		z.prev = z.next.prev
		z.next.prev = z
	} else {
		z.prev = z.parent
		t.last = z
	}
	if z.prev != t.nill {
		z.prev.next = z
	} else {
		t.first = z
	}

	return z, true
}

func (t *Rbtree) insertFixup(z *Rbnode) {
	for z.parent.color == RED {
		//
		// Howerver, we do not need the assertion of non-nil grandparent
		// because
		//
		//  2) The root is black
		//
		// Since the color of the parent is RED, so the parent is not root
		// and the grandparent must be exist.
		//
		if z.parent == z.parent.parent.left {
			// Take y as the uncle, although it can be nill, in that case
			// its color is BLACK
			y := z.parent.parent.right
			if y.color == RED {
				//
				// Case 1:
				// parent and uncle are both RED, the grandparent must be BLACK
				// due to
				//
				//  4) Both children of every red node are black
				//
				// Since the current node and its parent are all RED, we still
				// in violation of 4), So repaint both the parent and the uncle
				// to BLACK and grandparent to RED(to maintain 5)
				//
				//  5) Every simple path from root to leaves contains the same
				//     number of black nodes.
				//
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					//
					// Case 2:
					// parent is RED and uncle is BLACK and the current node
					// is right child
					//
					// A left rotation on the parent of the current node will
					// switch the roles of each other. This still leaves us in
					// violation of 4).
					// The continuation into Case 3 will fix that.
					//
					z = z.parent
					t.leftRotate(z)
				}
				//
				// Case 3:
				// parent is RED and uncle is BLACK and the current node is
				// left child
				//
				// At the very beginning of Case 3, current node and parent are
				// both RED, thus we violate 4).
				// Repaint parent to BLACK will fix it, but 5) does not allow
				// this because all paths that go through the parent will get
				// 1 more black node. Then repaint grandparent to RED (as we
				// discussed before, the grandparent is BLACK) and do a right
				// rotation will fix that.
				//
				z.parent.color = BLACK
				z.parent.parent.color = RED
				t.rightRotate(z.parent.parent)
			}
		} else { // same as then clause with "right" and "left" exchanged
			y := z.parent.parent.left
			if y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					t.rightRotate(z)
				}
				z.parent.color = BLACK
				z.parent.parent.color = RED
				t.leftRotate(z.parent.parent)
			}
		}
	}
	t.root.color = BLACK
}

// Just traverse the node from root to left recursively until left is nill.
// The node whose left is nill is the node with minimum value.
func (t *Rbtree) min(x *Rbnode) *Rbnode {
	if x == t.nill {
		return t.nill
	}

	for x.left != t.nill {
		x = x.left
	}

	return x
}

// Just traverse the node from root to right recursively until right is nill.
// The node whose right is nill is the node with maximum value.
func (t *Rbtree) max(x *Rbnode) *Rbnode {
	if x == t.nill {
		return t.nill
	}

	for x.right != t.nill {
		x = x.right
	}

	return x
}

func (t *Rbtree) search(x *Rbnode) *Rbnode {
	p := t.root

	for p != t.nill {

		if p.Item.Less(x.Item) {
			p = p.right
		} else if x.Item.Less(p.Item) {
			p = p.left
		} else {
			break
		}
	}

	return p
}

func (t *Rbtree) successor(x *Rbnode) *Rbnode {
	if x == t.nill {
		return t.nill
	}

	// Get the minimum from the right sub-tree if it existed.
	if x.right != t.nill {
		return t.min(x.right)
	}

	y := x.parent
	for y != t.nill && x == y.right {
		x = y
		y = y.parent
	}
	return y
}

func (t *Rbtree) transplant(u *Rbnode, v *Rbnode) {
	if u.parent == t.nill {
		t.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	v.parent = u.parent
}

func (t *Rbtree) remove(key *Rbnode) (*Rbnode, bool) {
	z := t.search(key)

	if z == t.nill {
		return nil, false
	}

	return t.remove_raw(z)
}

func (t *Rbtree) remove_raw(z *Rbnode) (*Rbnode, bool) {
	y := z
	yOriginalColor := y.color
	var x *Rbnode

	if z.left == t.nill {
		// one child (RIGHT)
		x = z.right
		t.transplant(z, z.right)

	} else if z.right == t.nill {
		// one child (LEFT)
		x = z.left
		t.transplant(z, z.left)

	} else {
		// two children
		y := z.right
		for y.left != t.nill {
			y = y.left
		}

		yOriginalColor = y.color
		x = y.right

		if y.parent == z {
			x.parent = y
		} else {
			t.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		t.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}

	if yOriginalColor == BLACK {
		t.deleteFixup(x)
	}

	t.count--

	if z.next != t.nill {
		z.next.prev = z.prev
	} else {
		t.last = z.prev
	}
	if z.prev != t.nill {
		z.prev.next = z.next
	} else {
		t.first = z.next
	}

	return z, true
}

func (t *Rbtree) deleteFixup(x *Rbnode) {
	for x != t.root && x.color == BLACK {
		if x == x.parent.left {
			w := x.parent.right
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				t.leftRotate(x.parent)
				w = x.parent.right
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				if w.right.color == BLACK {
					w.left.color = BLACK
					w.color = RED
					t.rightRotate(w)
					w = x.parent.right
				}
				w.color = x.parent.color
				x.parent.color = BLACK
				w.right.color = BLACK
				t.leftRotate(x.parent)
				x = t.root
			}
		} else {
			w := x.parent.left
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				t.rightRotate(x.parent)
				w = x.parent.left
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				if w.left.color == BLACK {
					w.right.color = BLACK
					w.color = RED
					t.leftRotate(w)
					w = x.parent.left
				}
				w.color = x.parent.color
				x.parent.color = BLACK
				w.left.color = BLACK
				t.rightRotate(x.parent)
				x = t.root
			}
		}
	}
	x.color = BLACK
}

// ========================== Tests Method ================================

/*
	preorder traversal recursive function
	print formated infomation to bytes.buffer
*/
func traverseP(node *Rbnode, nill *Rbnode, b *bytes.Buffer) {
	if node.left != nill {
		traverseP(node.left, nill, b)
	}

	fmt.Fprint(b, "[", node.prev.Item, ",", node.Item, ",", node.next.Item, ";")
	fmt.Fprintf(b, " %p,%p,%p]", node.prev, node, node.next)

	if node.right != nill {
		traverseP(node.right, nill, b)
	}
}

/*
	start the preorder traversal
	print all pointer address to a string
*/
func traversePrint(tree *Rbtree) string {
	b := bytes.NewBufferString("")
	if tree.root != tree.nill {
		traverseP(tree.root, tree.nill, b)
	}
	return b.String()
}

/*
	traverse by double-link
	print all pointer address to a string
*/
func linkedPrint(tree *Rbtree) string {
	b := bytes.NewBufferString("")
	for node := tree.First(); node != tree.nill; node = node.Next() {
		fmt.Fprint(b, "[", node.prev.Item, ",", node.Item, ",", node.next.Item, ";")
		fmt.Fprintf(b, " %p,%p,%p]", node.prev, node, node.next)
	}
	return b.String()
}

/*
	compare and check pointer between preorder traversal and double-link traversal
*/
func testPointer(t *Rbtree) {
	trastr := traversePrint(t)
	linstr := linkedPrint(t)
	if linstr != trastr {
		fmt.Println("tra:", trastr)
		fmt.Println("lin:", linstr)
		panic("test Pointer fail.")
	}
}

/*
	traverse the node to check whether the tree is a Binary-Search-Tree(BST)
	( node.left.Less(node) == true && node.Less(node.right) == true )
*/
func testBST(node *Rbnode, nill *Rbnode, count *int) {
	(*count)++
	if node.left != nill {
		// fmt.Println(node.left.Item)
		if !node.left.Item.Less(node.Item) {
			panic("rbtree BST error")
		}
		testBST(node.left, nill, count)
	}

	if node.right != nill {
		if !node.Item.Less(node.right.Item) {
			panic("rbtree BST error")
		}
		testBST(node.right, nill, count)
	}
}

/*
	traverse the node to check whether the tree has the same black node number from root to every single leaf node
*/
func testBlack(node *Rbnode, nill *Rbnode, blackDep int, total *int) {
	if node == nill {
		if *total == -1 {
			*total = blackDep
		} else {
			if *total != blackDep {
				panic("rbtree Black Dep Error")
			}
		}
		return
	}
	if node.color == BLACK {
		blackDep++
	}
	testBlack(node.left, nill, blackDep, total)
	testBlack(node.right, nill, blackDep, total)
}

/*
	This is a test method containing some testings below:
	1. whether is a BST
	2. whether tree.count is equal to preorder traversal count
	3. whether tree.count is equal to double-link(prev and next) count
	4. whether double-link satisfy key order with prev.Less(cur)
	5. whether preorder traversal pointer print is equal to double-link pointer print
	6. whether the red-black tree satisfies that every simple path from root to leaves contains the same number of black nodes
*/
func (tree *Rbtree) testStructure() {
	var count int = 0
	// root := tree.root
	if tree.root != tree.nill {
		testBST(tree.root, tree.nill, &count)
	}
	if count != tree.count {
		panic("testBST count error")
	}

	count = 0
	for p := tree.First(); p != tree.nill; p = p.Next() {
		//		fmt.Printf("%v ", p.Item)
		count++
		if p != tree.First() {
			if !p.Prev().Item.Less(p.Item) {
				panic("double link next error")
			}
		}
	}
	//	fmt.Printf("\n")
	if count != tree.count {
		fmt.Println("cnt: ", count, tree.count)
		panic("test double link next count error")
	}

	count = 0
	for p := tree.Last(); p != tree.nill; p = p.Prev() {
		count++
		if p != tree.First() {
			if !p.Prev().Item.Less(p.Item) {
				panic("double link prev error")
			}
		}
	}
	if count != tree.count {
		panic("test double link prev count error")
	}

	total := -1
	testBlack(tree.root, tree.nill, 0, &total)

	testPointer(tree)
}
