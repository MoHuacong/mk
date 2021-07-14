package tools

type Node struct {
	data interface{}
	next *Node
}

type Stack struct {
	topNode *Node
}

func NewStack() *Stack {
	return &Stack{nil}
}

func (t *Stack) Push(value interface{}) {
	t.topNode = &Node{data: value, next: t.topNode}
}

func (t *Stack) isEmpty() bool {
	if t.topNode == nil {
		return true
	}
	return false
}
func (t *Stack) Pop() interface{} {
	if t.isEmpty() { return nil }
	value := t.topNode.data
	t.topNode = t.topNode.next
	return value
}

func (t *Stack) Top() interface{} {
	if t.isEmpty() { return nil }
	return t.topNode.data
}