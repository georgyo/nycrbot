package math

import (
	"errors"
	"fmt"
)

type NodeType uint8

const (
	NT_Operator NodeType = iota
	NT_Number
)

type Node struct {
	Left   *Node
	Right  *Node
	Parent *Node
	typ    NodeType
	opval  rune
	numval *Numnum
}

func (n *Node) IsOperator() bool {
	return n.typ == NT_Operator
}

func (n *Node) SetType(nt NodeType) {
	n.typ = nt
}

type Tree struct {
	root          *Node
	nextNumParent *Node
	rootStack     []*Node
}

func (t *Tree) StackRoot() {
	t.rootStack = append(t.rootStack, t.root)
	t.root = nil
}

func (t *Tree) PopRoot() {
	oldRoot := t.root
	t.root = t.rootStack[len(t.rootStack)-1]
	t.rootStack = t.rootStack[0 : len(t.rootStack)-1]
	if t.root == nil {
		t.root = oldRoot
		return
	}
	if oldRoot == nil {
		return
	}
	oldRoot.Parent = t.root
	if t.root.Left == nil {
		t.root.Left = oldRoot
	} else {
		t.root.Right = oldRoot
	}
}

func (t *Tree) AddOperator(op rune) {
	if t.root == nil {
		t.root = &Node{
			typ:   NT_Operator,
			opval: op,
		}
		return
	}
	if (op == MUL || op == DIV) && t.root.Left != nil {
		nx := &Node{
			typ:    NT_Operator,
			opval:  op,
			Parent: t.root,
		}
		nx.Left = t.root.Right
		t.root.Right = nx
		t.nextNumParent = nx
		return
	}
	t.root.Parent = &Node{
		typ:   NT_Operator,
		opval: op,
		Left:  t.root,
	}
	t.root = t.root.Parent
}

func (t *Tree) AddNumber(num *Numnum) {
	n := &Node{
		typ:    NT_Number,
		numval: num,
	}
	if t.nextNumParent != nil {
		t.nextNumParent.Right = n
		t.nextNumParent = nil
		return
	}
	if t.root == nil {
		t.root = n
		return
	}
	if t.root.Left == nil {
		t.root.Left = n
	} else {
		t.root.Right = n
	}
}

func ReduceNode(n *Node) *Node {
	if n == nil {
		return nil
	}
	if n.Left == nil || n.Right == nil {
		fmt.Printf("Cant reduce this bro %s\n", n)
		return nil
	}
	fmt.Printf("ReduceNode: %s, Left: %s, Right: %s\n", n.numval, n.Left.numval, n.Right.numval)
	if n.Left.IsOperator() {
		n.Left = ReduceNode(n.Left)
	}
	if n.Right.IsOperator() {
		n.Right = ReduceNode(n.Right)
	}
	if n.Left == nil || n.Right == nil {
		fmt.Printf("Cant reduce this bro, previous reduce turn node to nil\n")
		return nil
	}
	result, err := (*n.Left.numval).ExecOp(n.opval, *n.Right.numval)
	if err != nil {
		return nil
	}
	return &Node{
		typ:    NT_Number,
		numval: &result,
		Parent: n.Parent,
	}
}

func (t *Tree) Parse() (Numnum, error) {
	resultNode := ReduceNode(t.root)
	if resultNode == nil {
		return 0, errors.New("Something failed")
	}
	return *resultNode.numval, nil
}

func NewTree() Tree {
	t := Tree{}
	t.rootStack = make([]*Node, 0)
	return t
}
