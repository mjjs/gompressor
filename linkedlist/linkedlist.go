package linkedlist

type node struct {
	value interface{}
	next  *node
}

// LinkedList is the main structure of the linked list.
type LinkedList struct {
	head *node
	tail *node
	size int
}

// New returns a pointer to a new LinkedList and sets root as the head and tail nodes.
func New(root interface{}) *LinkedList {
	n := &node{value: root}

	return &LinkedList{
		head: n,
		tail: n,
		size: 1,
	}
}

// Append appends the given value to the end of the list.
func (ll *LinkedList) Append(value interface{}) {
	newNode := &node{value: value}

	if ll.size == 0 {
		ll.head = newNode
		ll.tail = newNode
	} else {
		ll.tail.next = newNode
		ll.tail = newNode
	}

	ll.size++
}

// Remove removes value from the linked list.
func (ll *LinkedList) Remove(value interface{}) {
	if ll.head.value == value {
		ll.head = ll.head.next
	}

	removed := ll.remove(ll.head, value)
	if removed {
		ll.size--
	}
}

// Find finds the given value from the linked list and returns it if found.
// A boolean value is returned indicating whether or not the value actually
// exists in the linked list.
func (ll *LinkedList) Find(value interface{}) (interface{}, bool) {
	result := find(ll.head, value)
	if result == nil {
		return nil, false
	}

	return result.value, true
}

// Head returns the head of the linked list.
func (ll *LinkedList) Head() interface{} {
	if ll.head == nil {
		return nil
	}

	return ll.head.value
}

// Tail returns the tail of the linked list.
func (ll *LinkedList) Tail() interface{} {
	if ll.tail == nil {
		return nil
	}

	return ll.tail.value
}

// Size returns the number of elements in the linked list.
func (ll *LinkedList) Size() int {
	return ll.size
}

// ForEach executes f for each element in the list.
func (ll *LinkedList) ForEach(f func(val interface{})) {
	execFunc(ll.head, f)
}

func (ll *LinkedList) remove(n *node, val interface{}) bool {
	if n == nil || n.next == nil {
		return false
	}

	if n.next.value == val {
		if n.next == ll.tail {
			ll.tail = n
		}

		n.next = n.next.next
		return true
	}

	return ll.remove(n.next, val)
}

func execFunc(n *node, f func(val interface{})) {
	if n == nil {
		return
	}

	f(n.value)

	execFunc(n.next, f)
}

func find(n *node, val interface{}) *node {
	if n == nil {
		return nil
	}

	if n.value == val {
		return n
	}

	return find(n.next, val)
}
