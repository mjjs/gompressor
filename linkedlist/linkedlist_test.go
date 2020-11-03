package linkedlist

import (
	"testing"
)

func TestNewStartsWithGivenNodeAndCorrectSize(t *testing.T) {
	node := &Node{Value: 1}
	ll := New(node)

	if ll.head != node {
		t.Error("Expected head to be set correctly")
	}

	if ll.size != 1 {
		t.Error("Expected size to be initialized to 1 for a Linkedlist created with New")
	}
}

func TestAppend(t *testing.T) {
	ll := &LinkedList{}

	ll.Append(1)
	if ll.size != 1 {
		t.Errorf("Expected size to be 1, got %d", ll.size)
	}

	ll.Append(2)
	if ll.size != 2 {
		t.Errorf("Expected size to be 2, got %d", ll.size)
	}

	ll.Append(3)

	if ll.size != 3 {
		t.Errorf("Expected size to be 3, got %d", ll.size)
	}

	if ll.head.Value != 1 {
		t.Errorf("Expected %v to be head, got %v", 1, ll.head.Value)
	}

	if ll.head.next.Value != 2 {
		t.Error("Expected elements to be appended into the list")
	}

	if ll.head.next.next.Value != 3 {
		t.Error("Expected elements to be appended into the list")
	}
}

func TestRemoveCanRemoveHead(t *testing.T) {
	ll := &LinkedList{}

	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	ll.Remove(1)

	if ll.head.Value != 2 {
		t.Errorf("Expected %v to become head, got %v", 2, ll.head)
	}
}

func TestRemoveCanRemoveTail(t *testing.T) {
	ll := &LinkedList{}

	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	ll.Remove(3)

	if ll.tail.Value != 2 {
		t.Errorf("Expected %v to become tail, got %v", 2, ll.tail.Value)
	}
}

func TestRemoveCanRemoveBetweenValues(t *testing.T) {
	ll := &LinkedList{}

	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	ll.Remove(2)

	if ll.head.next != ll.tail {
		t.Error("Expected middle value to be removed")
	}
}

func TestReturnsHeadCorrectly(t *testing.T) {
	ll := &LinkedList{}

	ll.Append(1)
	ll.Append(2)

	if ll.Head().Value != 1 {
		t.Errorf("Expected %v to be head, got %v", 1, ll.Head().Value)
	}
}

func TestReturnsTailCorrectly(t *testing.T) {
	ll := &LinkedList{}

	ll.Append(1)
	ll.Append(2)

	if ll.Tail().Value != 2 {
		t.Errorf("Expected %v to be tail, got %v", 2, ll.Tail().Value)
	}
}

func TestReturnsSizeCorrectly(t *testing.T) {
	ll := &LinkedList{}

	ll.Append(1)
	ll.Append(2)

	if ll.Size() != 2 {
		t.Errorf("Expected Size to return 2, got %d", ll.Size())
	}
}

func TestFind(t *testing.T) {
	ll := &LinkedList{}

	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	if res, exists := ll.Find(2); exists == false {
		t.Error("Expected Find to return true when value exists in the list")
	} else if res != 2 {
		t.Errorf("Expected %v to be returned, got %v", 2, res)
	}

	if _, exists := ll.Find(5); exists == true {
		t.Error("Expected Find to return false when value does not exist in the list")
	}
}
