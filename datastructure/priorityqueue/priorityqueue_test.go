package priorityqueue

import (
	"testing"
)

func TestDequeueReturnsSmallestPriorityFirst(t *testing.T) {
	pq := &PriorityQueue{}
	pq.Enqueue(1, 1)
	pq.Enqueue(2, 2)
	pq.Enqueue(1, 3)
	pq.Enqueue(5, 4)

	if _, val := pq.Dequeue(); val != 1 {
		t.Errorf("Expected %d, got %v", 1, val)
	}

	if _, val := pq.Dequeue(); val != 3 {
		t.Errorf("Expected %d, got %v", 3, val)
	}

	if _, val := pq.Dequeue(); val != 2 {
		t.Errorf("Expected %d, got %v", 4, val)
	}

	if _, val := pq.Dequeue(); val != 4 {
		t.Errorf("Expected %d, got %v", 5, val)
	}
}

func TestDequeueReturnsFirstQueuedValueForEqualPriorities(t *testing.T) {
	pq := &PriorityQueue{}
	pq.Enqueue(1, 1)
	pq.Enqueue(1, 2)
	pq.Enqueue(2, 4)
	pq.Enqueue(2, 5)

	if _, val := pq.Dequeue(); val != 1 {
		t.Errorf("Expected %d, got %v", 1, val)
	}

	if _, val := pq.Dequeue(); val != 2 {
		t.Errorf("Expected %d, got %v", 2, val)
	}

	if _, val := pq.Dequeue(); val != 4 {
		t.Errorf("Expected %d, got %v", 4, val)
	}

	if _, val := pq.Dequeue(); val != 5 {
		t.Errorf("Expected %d, got %v", 5, val)
	}
}

func TestPeekReturnsSmallestPriorityWithoutRemovingIt(t *testing.T) {
	pq := &PriorityQueue{}
	pq.Enqueue(9, "g")
	pq.Enqueue(5, "a")
	pq.Enqueue(8, "b")

	for i := 0; i < 5; i++ {
		if prio, val := pq.Peek(); val != "a" {
			t.Errorf("Expected peek to return first element on all calls, got %v", val)
		} else if prio != 5 {
			t.Errorf("Expected priority to be %d, got %d", 5, prio)
		}
	}
}

func TestPeekReturnsNilOnEmptyQueue(t *testing.T) {
	pq := &PriorityQueue{}
	if _, val := pq.Peek(); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestDequeueReturnsNilOnEmptyQueue(t *testing.T) {
	pq := &PriorityQueue{}
	if _, val := pq.Dequeue(); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestDequeueRemovesNodes(t *testing.T) {
	pq := &PriorityQueue{}
	pq.Enqueue(3, "a")

	if pq.nodes.Size() != 1 {
		t.Errorf("Expected size to be %d got %d", pq.nodes.Size(), 1)
	}

	pq.Dequeue()

	if pq.nodes.Size() != 0 {
		t.Errorf("Expected size to be %d got %d", pq.nodes.Size(), 0)
	}
}

func TestSize(t *testing.T) {
	pq := &PriorityQueue{}

	if n := pq.Size(); n != 0 {
		t.Errorf("Expected size to be 0, got %d", n)
	}

	pq.Enqueue(1, 1)

	if n := pq.Size(); n != 1 {
		t.Errorf("Expected size to be 1, got %d", n)
	}

	pq.Dequeue()

	if n := pq.Size(); n != 0 {
		t.Errorf("Expected size to be 0, got %d", n)
	}
}
