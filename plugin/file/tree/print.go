package tree

import "fmt"

func (t *Tree) Print() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.Root == nil {
		fmt.Println("<nil>")
	}
	t.Root.print()
}
func (n *Node) print() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	q := newQueue()
	q.push(n)
	nodesInCurrentLevel := 1
	nodesInNextLevel := 0
	for !q.empty() {
		do := q.pop()
		nodesInCurrentLevel--
		if do != nil {
			fmt.Print(do.Elem.Name(), " ")
			q.push(do.Left)
			q.push(do.Right)
			nodesInNextLevel += 2
		}
		if nodesInCurrentLevel == 0 {
			fmt.Println()
		}
		nodesInCurrentLevel = nodesInNextLevel
		nodesInNextLevel = 0
	}
	fmt.Println()
}

type queue []*Node

func newQueue() queue {
	_logClusterCodePath()
	defer _logClusterCodePath()
	q := queue([]*Node{})
	return q
}
func (q *queue) push(n *Node) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*q = append(*q, n)
}
func (q *queue) pop() *Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n := (*q)[0]
	*q = (*q)[1:]
	return n
}
func (q *queue) empty() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(*q) == 0
}
