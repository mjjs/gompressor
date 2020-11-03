package bytevector

import (
	"errors"
	"fmt"
)

const defaultCapacity int = 10

// ErrIndexOutOfRange is an error that is returned when accessing the vector
// using an index that is negative or greater than the capacity of the given
// vector.
var ErrIndexOutOfRange = errors.New("index out of range")

// Bytevector is a dynamic array that holds bytes.
type Bytevector struct {
	size     int
	capacity int
	elements []byte
}

// New returns a pointer to a Bytevector. The size specifies the length of the
// initial vector. The vector will hold size amount of zero valued bytes.
//
// An additional uint can be passed to New, which indicates how many values
// the underlying array can hold before resizing. If capacity is smaller than length,
// capacity will be set to equal length.
func New(size ...uint) *Bytevector {
	switch len(size) {
	case 1:
		return &Bytevector{
			size:     int(size[0]),
			capacity: int(size[0]),
			elements: make([]byte, size[0]),
		}

	case 2:
		var sz = int(size[0])
		var capacity = int(size[1])

		if capacity < sz {
			capacity = sz
		}

		return &Bytevector{
			capacity: capacity,
			size:     sz,
			elements: make([]byte, capacity),
		}

	default:
		return &Bytevector{
			elements: make([]byte, 0),
		}
	}
}

// Append adds the values to the end of the vector, growing it if necessary.
func (bv *Bytevector) Append(values ...byte) {
	for _, value := range values {
		if bv.size == bv.capacity {
			bv.grow()
		}

		bv.elements[bv.size] = value
		bv.size++
	}
}

// AppendToCopy creates a copy of bv and appends the values to the end of the
// copy and returns the copy. The copy is grown if necessary.
func (bv *Bytevector) AppendToCopy(values ...byte) *Bytevector {
	newBV := New(uint(bv.size), uint(bv.capacity))

	for i := 0; i < bv.size; i++ {
		newBV.Set(i, bv.MustGet(i))
	}

	newBV.Append(values...)

	return newBV
}

// Get returns the byte value at the given index. Returns an error if the index
// is out or range.
func (bv *Bytevector) Get(index int) (byte, error) {
	if index < 0 || index >= bv.size {
		return byte(0), fmt.Errorf("%w: %d", ErrIndexOutOfRange, index)
	}

	return bv.elements[index], nil
}

// MustGet returns the byte value at the given index without performing any bounds checking.
func (bv *Bytevector) MustGet(index int) byte {
	return bv.elements[index]
}

// Set sets value into the given index. Returns an error if the index is out of range.
func (bv *Bytevector) Set(index int, value byte) error {
	if index < 0 || index >= bv.size {
		return fmt.Errorf("%w: %d", ErrIndexOutOfRange, index)
	}

	bv.elements[index] = value

	return nil
}

// MustSet sets value into the given index without performing any bounds checking.
func (bv *Bytevector) MustSet(index int, value byte) {
	bv.elements[index] = value
}

// Size returns the number of elements in the vector.
func (bv *Bytevector) Size() int {
	return bv.size
}

// Capacity returns the number of elements the vector can hold before resizing.
func (bv *Bytevector) Capacity() int {
	return bv.capacity
}

// String returns the string representation of the bytes in the vector.
func (bv *Bytevector) String() string {
	s := ""

	for i := 0; i < bv.size; i++ {
		s += string([]byte{bv.elements[i]})
	}

	return s
}

func (bv *Bytevector) grow() {
	newCapacity := bv.capacity * 2

	if bv.capacity == 0 {
		newCapacity = 1
	}

	newElements := make([]byte, newCapacity)

	for i := 0; i < bv.size; i++ {
		newElements[i] = bv.elements[i]
	}

	bv.capacity = newCapacity
	bv.elements = newElements
}
