package tree

import "github.com/miekg/dns"

const (
	td234	= iota
	bu23
)
const mode = bu23

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if mode != td234 && mode != bu23 {
		panic("tree: unknown mode")
	}
}

type Color bool

const (
	red	Color	= false
	black	Color	= true
)

type Node struct {
	Elem		*Elem
	Left, Right	*Node
	Color		Color
}
type Tree struct {
	Root	*Node
	Count	int
}

func (n *Node) color() Color {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if n == nil {
		return black
	}
	return n.Color
}
func (n *Node) rotateLeft() (root *Node) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	root = n.Right
	n.Right = root.Left
	root.Left = n
	root.Color = n.Color
	n.Color = red
	return
}
func (n *Node) rotateRight() (root *Node) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	root = n.Left
	n.Left = root.Right
	root.Right = n
	root.Color = n.Color
	n.Color = red
	return
}
func (n *Node) flipColors() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n.Color = !n.Color
	n.Left.Color = !n.Left.Color
	n.Right.Color = !n.Right.Color
}
func (n *Node) fixUp() *Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if n.Right.color() == red {
		if mode == td234 && n.Right.Left.color() == red {
			n.Right = n.Right.rotateRight()
		}
		n = n.rotateLeft()
	}
	if n.Left.color() == red && n.Left.Left.color() == red {
		n = n.rotateRight()
	}
	if mode == bu23 && n.Left.color() == red && n.Right.color() == red {
		n.flipColors()
	}
	return n
}
func (n *Node) moveRedLeft() *Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n.flipColors()
	if n.Right.Left.color() == red {
		n.Right = n.Right.rotateRight()
		n = n.rotateLeft()
		n.flipColors()
		if mode == td234 && n.Right.Right.color() == red {
			n.Right = n.Right.rotateLeft()
		}
	}
	return n
}
func (n *Node) moveRedRight() *Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n.flipColors()
	if n.Left.Left.color() == red {
		n = n.rotateRight()
		n.flipColors()
	}
	return n
}
func (t *Tree) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.Count
}
func (t *Tree) Search(qname string) (*Elem, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.Root == nil {
		return nil, false
	}
	n, res := t.Root.search(qname)
	if n == nil {
		return nil, res
	}
	return n.Elem, res
}
func (n *Node) search(qname string) (*Node, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for n != nil {
		switch c := Less(n.Elem, qname); {
		case c == 0:
			return n, true
		case c < 0:
			n = n.Left
		default:
			n = n.Right
		}
	}
	return n, false
}
func (t *Tree) Insert(rr dns.RR) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var d int
	t.Root, d = t.Root.insert(rr)
	t.Count += d
	t.Root.Color = black
}
func (n *Node) insert(rr dns.RR) (root *Node, d int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if n == nil {
		return &Node{Elem: newElem(rr)}, 1
	} else if n.Elem == nil {
		n.Elem = newElem(rr)
		return n, 1
	}
	if mode == td234 {
		if n.Left.color() == red && n.Right.color() == red {
			n.flipColors()
		}
	}
	switch c := Less(n.Elem, rr.Header().Name); {
	case c == 0:
		n.Elem.Insert(rr)
	case c < 0:
		n.Left, d = n.Left.insert(rr)
	default:
		n.Right, d = n.Right.insert(rr)
	}
	if n.Right.color() == red && n.Left.color() == black {
		n = n.rotateLeft()
	}
	if n.Left.color() == red && n.Left.Left.color() == red {
		n = n.rotateRight()
	}
	if mode == bu23 {
		if n.Left.color() == red && n.Right.color() == red {
			n.flipColors()
		}
	}
	root = n
	return
}
func (t *Tree) DeleteMin() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.Root == nil {
		return
	}
	var d int
	t.Root, d = t.Root.deleteMin()
	t.Count += d
	if t.Root == nil {
		return
	}
	t.Root.Color = black
}
func (n *Node) deleteMin() (root *Node, d int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if n.Left == nil {
		return nil, -1
	}
	if n.Left.color() == black && n.Left.Left.color() == black {
		n = n.moveRedLeft()
	}
	n.Left, d = n.Left.deleteMin()
	root = n.fixUp()
	return
}
func (t *Tree) DeleteMax() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.Root == nil {
		return
	}
	var d int
	t.Root, d = t.Root.deleteMax()
	t.Count += d
	if t.Root == nil {
		return
	}
	t.Root.Color = black
}
func (n *Node) deleteMax() (root *Node, d int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if n.Left != nil && n.Left.color() == red {
		n = n.rotateRight()
	}
	if n.Right == nil {
		return nil, -1
	}
	if n.Right.color() == black && n.Right.Left.color() == black {
		n = n.moveRedRight()
	}
	n.Right, d = n.Right.deleteMax()
	root = n.fixUp()
	return
}
func (t *Tree) Delete(rr dns.RR) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.Root == nil {
		return
	}
	el, _ := t.Search(rr.Header().Name)
	if el == nil {
		t.deleteNode(rr)
		return
	}
	empty := el.Delete(rr)
	if empty {
		t.deleteNode(rr)
		return
	}
}
func (t *Tree) deleteNode(rr dns.RR) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.Root == nil {
		return
	}
	var d int
	t.Root, d = t.Root.delete(rr)
	t.Count += d
	if t.Root == nil {
		return
	}
	t.Root.Color = black
}
func (n *Node) delete(rr dns.RR) (root *Node, d int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if Less(n.Elem, rr.Header().Name) < 0 {
		if n.Left != nil {
			if n.Left.color() == black && n.Left.Left.color() == black {
				n = n.moveRedLeft()
			}
			n.Left, d = n.Left.delete(rr)
		}
	} else {
		if n.Left.color() == red {
			n = n.rotateRight()
		}
		if n.Right == nil && Less(n.Elem, rr.Header().Name) == 0 {
			return nil, -1
		}
		if n.Right != nil {
			if n.Right.color() == black && n.Right.Left.color() == black {
				n = n.moveRedRight()
			}
			if Less(n.Elem, rr.Header().Name) == 0 {
				n.Elem = n.Right.min().Elem
				n.Right, d = n.Right.deleteMin()
			} else {
				n.Right, d = n.Right.delete(rr)
			}
		}
	}
	root = n.fixUp()
	return
}
func (t *Tree) Min() *Elem {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.Root == nil {
		return nil
	}
	return t.Root.min().Elem
}
func (n *Node) min() *Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for ; n.Left != nil; n = n.Left {
	}
	return n
}
func (t *Tree) Max() *Elem {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.Root == nil {
		return nil
	}
	return t.Root.max().Elem
}
func (n *Node) max() *Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for ; n.Right != nil; n = n.Right {
	}
	return n
}
func (t *Tree) Prev(qname string) (*Elem, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.Root == nil {
		return nil, false
	}
	n := t.Root.floor(qname)
	if n == nil {
		return nil, false
	}
	return n.Elem, true
}
func (n *Node) floor(qname string) *Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if n == nil {
		return nil
	}
	switch c := Less(n.Elem, qname); {
	case c == 0:
		return n
	case c <= 0:
		return n.Left.floor(qname)
	default:
		if r := n.Right.floor(qname); r != nil {
			return r
		}
	}
	return n
}
func (t *Tree) Next(qname string) (*Elem, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.Root == nil {
		return nil, false
	}
	n := t.Root.ceil(qname)
	if n == nil {
		return nil, false
	}
	return n.Elem, true
}
func (n *Node) ceil(qname string) *Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if n == nil {
		return nil
	}
	switch c := Less(n.Elem, qname); {
	case c == 0:
		return n
	case c > 0:
		return n.Right.ceil(qname)
	default:
		if l := n.Left.ceil(qname); l != nil {
			return l
		}
	}
	return n
}
