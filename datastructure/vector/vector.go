// Package vector implements a dynamically growing array list.
package vector

import (
	"errors"
	"fmt"
)

const defaultCapacity int = 10

// ErrIndexOutOfRange is an error that is returned when accessing the vector
// using an index that is negative or greater than the capacity of the given
// vector.
var ErrIndexOutOfRange = errors.New("index out of range")

// Vector is a dynamic array that holds bytes.
type Vector struct {
	size     int
	capacity int
	elements []interface{}
}

// New returns a pointer to a Vector. The size specifies the length of the
// initial vector. The vector will hold size amount of zero valued bytes.
//
// An additional uint can be passed to New, which indicates how many values
// the underlying array can hold before resizing. If capacity is smaller than length,
// capacity will be set to equal length.
func New(size ...uint) *Vector {
	switch len(size) {
	case 1:
		return &Vector{
			size:     int(size[0]),
			capacity: int(size[0]),
			elements: make([]interface{}, size[0]),
		}

	case 2:
		var sz = int(size[0])
		var capacity = int(size[1])

		if capacity < sz {
			capacity = sz
		}

		return &Vector{
			capacity: capacity,
			size:     sz,
			elements: make([]interface{}, capacity),
		}

	default:
		return &Vector{
			elements: make([]interface{}, 0),
		}
	}
}

// Append adds the values to the end of the vector, growing it if necessary.
func (v *Vector) Append(values ...interface{}) {
	for _, value := range values {
		if v.size == v.capacity {
			v.grow()
		}

		v.elements[v.size] = value
		v.size++
	}
}

// Pop removes the last element from the vector.
func (v *Vector) Pop() interface{} {
	if v.size == 0 {
		return nil
	}

	tail := v.elements[v.size-1]
	v.size--
	return tail
}

// AppendToCopy creates a copy of v and appends the values to the end of the
// copy and returns the copy. The copy is grown if necessary.
func (v *Vector) AppendToCopy(values ...interface{}) *Vector {
	newVector := New(uint(v.size), uint(v.capacity))

	for i := 0; i < v.size; i++ {
		newVector.Set(i, v.MustGet(i))
	}

	newVector.Append(values...)

	return newVector
}

// Get returns the value at the given index. Returns an error if the index
// is out or range.
func (v *Vector) Get(index int) (interface{}, error) {
	if index < 0 || index >= v.size {
		return nil, fmt.Errorf("%w: %d", ErrIndexOutOfRange, index)
	}

	return v.elements[index], nil
}

// MustGet returns the value at the given index without performing any bounds checking.
func (v *Vector) MustGet(index int) interface{} {
	return v.elements[index]
}

// Set sets value into the given index. Returns an error if the index is out of range.
func (v *Vector) Set(index int, value interface{}) error {
	if index < 0 || index >= v.size {
		return fmt.Errorf("%w: %d", ErrIndexOutOfRange, index)
	}

	v.elements[index] = value

	return nil
}

// MustSet sets value into the given index without performing any bounds checking.
func (v *Vector) MustSet(index int, value interface{}) {
	v.elements[index] = value
}

// Size returns the number of elements in the vector.
func (v *Vector) Size() int {
	return v.size
}

// Capacity returns the number of elements the vector can hold before resizing.
func (v *Vector) Capacity() int {
	return v.capacity
}

// String returns the string representation of the bytes in the vector.
func (v *Vector) String() string {
	s := ""

	for i := 0; i < v.size; i++ {
		switch val := v.elements[i].(type) {
		case byte:
			s += string([]byte{val})
		default:
			panic(fmt.Sprintf("String not implemented for %T", v.elements[i]))
		}
	}

	return s
}

func (v *Vector) grow() {
	newCapacity := v.capacity * 2

	if v.capacity == 0 {
		newCapacity = 1
	}

	newElements := make([]interface{}, newCapacity)

	for i := 0; i < v.size; i++ {
		newElements[i] = v.elements[i]
	}

	v.capacity = newCapacity
	v.elements = newElements
}
