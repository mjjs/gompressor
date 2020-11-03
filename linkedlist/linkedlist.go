package linkedlist

type Node struct {
	Value interface{}
	next  *Node
}

type LinkedList struct {
	head *Node
	tail *Node
	size int
}

func New(root *Node) *LinkedList {
	return &LinkedList{
		head: root,
		size: 1,
	}
}

func (ll *LinkedList) Append(value interface{}) {
	newNode := &Node{Value: value}

	if ll.size == 0 {
		ll.head = newNode
		ll.tail = newNode
	} else {
		ll.tail.next = newNode
		ll.tail = newNode
	}

	ll.size++
}

func (ll *LinkedList) Remove(val interface{}) {
	if ll.head.Value == val {
		ll.head = ll.head.next
	}

	ll.remove(ll.head, val)

	ll.size--
}

func (ll *LinkedList) Find(val interface{}) (interface{}, bool) {
	result := find(ll.head, val)
	if result == nil {
		return nil, false
	}

	return result.Value, true
}

func (ll *LinkedList) Head() *Node {
	return ll.head
}

func (ll *LinkedList) Tail() *Node {
	return ll.tail
}

func (ll *LinkedList) Size() int {
	return ll.size
}

func find(node *Node, val interface{}) *Node {
	if node == nil {
		return nil
	}

	if node.Value == val {
		return node
	}

	return find(node.next, val)
}

func (ll *LinkedList) remove(node *Node, val interface{}) {
	if node == nil || node.next == nil {
		return
	}

	if node.next.Value == val {
		if node.next == ll.tail {
			ll.tail = node
		}

		node.next = node.next.next
		return
	}

	ll.remove(node.next, val)
}
