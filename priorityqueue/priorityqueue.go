// Package priorityqueue implements implements a priority queue which considers
// lower values to be of higher priority.
package priorityqueue

import (
	"github.com/mjjs/gompressor/vector"
)

type node struct {
	priority int
	value    interface{}
}

// PriorityQueue is the type which implements the priority queue.
type PriorityQueue struct {
	nodes vector.Vector
	size  int
}

// Enqueue adds value to the queue with the given priority. Elements are re-ordered
// if needed so that the heap property is satisfied.
func (pq *PriorityQueue) Enqueue(priority int, value interface{}) {
	pq.nodes.Append(&node{priority: priority, value: value})
	pq.siftUp(pq.nodes.Size() - 1)
}

// Dequeue removes the element with the highest priority and returns it to the caller.
// The rest of the tree is re-ordered to satisfy the heap property.
func (pq *PriorityQueue) Dequeue() interface{} {
	if pq.nodes.Size() == 0 {
		return nil
	}

	head := pq.nodes.MustGet(0).(*node).value

	pq.nodes.MustSet(0, pq.nodes.MustGet(pq.nodes.Size()-1))
	pq.nodes.Pop()

	if pq.nodes.Size() > 1 {
		pq.siftDown(0)
	}

	return head
}

// Peek returns the highest priority element to the caller without removing it
// from the queue.
func (pq *PriorityQueue) Peek() interface{} {
	if pq.nodes.Size() == 0 {
		return nil
	}

	return pq.nodes.MustGet(0).(*node).value
}

func (pq *PriorityQueue) siftUp(i int) {
	parent := (i - 1) / 2

	currentNode := pq.nodes.MustGet(i).(*node)
	parentNode := pq.nodes.MustGet(parent).(*node)

	if currentNode.priority < parentNode.priority {
		pq.swap(i, parent)
		pq.siftUp(parent)
	}
}

func (pq *PriorityQueue) siftDown(i int) {
	if i >= pq.nodes.Size()/2 && i <= pq.nodes.Size() {
		return
	}

	smallest := i
	left := i*2 + 1
	right := i*2 + 2

	if left < pq.nodes.Size() && pq.nodes.MustGet(left).(*node).priority < pq.nodes.MustGet(smallest).(*node).priority {
		smallest = left
	}

	if right < pq.nodes.Size() && pq.nodes.MustGet(right).(*node).priority < pq.nodes.MustGet(smallest).(*node).priority {
		smallest = right
	}

	if smallest != i {
		pq.swap(i, smallest)
		pq.siftDown(smallest)
	}
}

func (pq *PriorityQueue) swap(i, j int) {
	temp := pq.nodes.MustGet(i)
	pq.nodes.MustSet(i, pq.nodes.MustGet(j))
	pq.nodes.MustSet(j, temp)
}
