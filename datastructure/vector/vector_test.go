package vector

import (
	"testing"
)

var newTestCases = []struct {
	name             string
	input            []uint
	expectedSize     int
	expectedCapacity int
}{
	{name: "No parameters", input: nil, expectedSize: 0, expectedCapacity: 0},
	{name: "Size given", input: []uint{5}, expectedSize: 5, expectedCapacity: 5},
	{name: "Size and capacity given", input: []uint{5, 15}, expectedSize: 5, expectedCapacity: 15},
	{name: "Capacity less than size", input: []uint{5, 3}, expectedSize: 5, expectedCapacity: 5},
}

func TestNew(t *testing.T) {
	for _, testCase := range newTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			bv := New(testCase.input...)
			if actualCapacity := bv.Capacity(); actualCapacity != testCase.expectedCapacity {
				t.Errorf("Expected capacity to be %d, got %d", testCase.expectedCapacity, actualCapacity)
			}

			if actualSize := bv.Size(); actualSize != testCase.expectedSize {
				t.Errorf("Expected size to be %d, got %d", testCase.expectedSize, actualSize)
			}
		})
	}
}

func TestAppendAddsValuesToEnd(t *testing.T) {
	bv := New()

	for i := 0; i < 5; i++ {
		bv.Append(byte(i + 1))
		val, err := bv.Get(i)

		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		if val != byte(i+1) {
			t.Errorf("Expected %x got %x", byte(i+1), val)
		}
	}
}

func TestCanAppendPastGivenSize(t *testing.T) {
	bv := New()

	for i := 0; i < 20; i++ {
		bv.Append(byte(i))

		if elem, err := bv.Get(i); err != nil {
			t.Errorf("Expected nil error, got %s", err)
		} else if elem != byte(i) {
			t.Errorf("Expected to find appended element %x, got %x", byte(i), elem)
		}
	}
}

func TestAppendToCopyCreatesNewCopy(t *testing.T) {
	bv := New()
	bv.Append(byte(1), byte(2), byte(3))
	newBV := bv.AppendToCopy(byte(123))

	if &bv == &newBV {
		t.Errorf("Expected AppendToCopy to create a copy of the original vector")
	}

	if newBV.Size() != bv.Size()+1 {
		t.Errorf("Expected newBV to include old elements as well as the new one")
	}
}

var setAndGetTestCases = []struct {
	name        string
	vector      *Vector
	index       int
	shouldError bool
}{
	{name: "Sets valid index without error", vector: New(5, 5), index: 2, shouldError: false},
	{name: "Panics on negative index", vector: New(5, 5), index: -1, shouldError: true},
	{name: "Panics on too large index", vector: New(5, 5), index: 50, shouldError: true},
}

func TestSetAndGet(t *testing.T) {
	for _, testCase := range setAndGetTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.vector.Set(testCase.index, byte(123))

			if err == nil && testCase.shouldError {
				t.Error("Expected an error, got nil")
			} else if err != nil && !testCase.shouldError {
				t.Errorf("Expected no error, got %s", err)
			}

			val, err := testCase.vector.Get(testCase.index)
			if err == nil && testCase.shouldError {
				t.Error("Expected an error, got nil")
			} else if err != nil && !testCase.shouldError {
				t.Errorf("Expected no error, got %s", err)
			}

			if !testCase.shouldError && val != byte(123) {
				t.Errorf("Expected %x got %x", byte(123), val)
			}
		})
	}
}

var mustSetAndGetTestCases = []struct {
	name        string
	vector      *Vector
	index       int
	input       byte
	shouldPanic bool
}{
	{name: "Sets valid index without panic", vector: New(5, 5), index: 2, shouldPanic: false},
	{name: "Panics on negative index", vector: New(5, 5), index: -1, shouldPanic: true},
	{name: "Panics on too large index", vector: New(5, 5), index: 50, shouldPanic: true},
}

func TestMustSetAndGet(t *testing.T) {
	for _, testCase := range mustSetAndGetTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer func() {
				r := recover()

				if testCase.shouldPanic && r == nil {
					t.Errorf("Expected MustSet to panic when index is out of range")
				} else if !testCase.shouldPanic && r != nil {
					t.Errorf("Expected MustSet not to panic on valid index")
				}
			}()

			testCase.vector.MustSet(testCase.index, byte(1))
		})
	}
}

func TestMustGet(t *testing.T) {
	for _, testCase := range mustSetAndGetTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer func() {
				r := recover()

				if testCase.shouldPanic && r == nil {
					t.Errorf("Expected MustGet to panic when index is out of range")
				} else if !testCase.shouldPanic && r != nil {
					t.Errorf("Expected MustGet not to panic on valid index")
				}
			}()

			testCase.vector.MustGet(testCase.index)
		})
	}
}

func TestString(t *testing.T) {
	bv := New()
	input := "Hello World"

	for _, c := range input {
		bv.Append(byte(c))
	}

	if actual := bv.String(); actual != input {
		t.Errorf("Expected %+v, got %+v", input, actual)
	}
}

func TestStringPanicsOnInvalidType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected String to panic on types that are out of this project's scope")
		}
	}()

	vec := New()
	vec.Append(1, 2, 3)

	_ = vec.String()
}

func TestPop(t *testing.T) {
	vec := New()
	vec.Append(1, 2, 3, 4)

	for i := 4; i > 0; i-- {
		if val := vec.Pop(); val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}

		if vec.Size() != i-1 {
			t.Errorf("Expected size to be %d, got %d", i-1, vec.Size())
		}
	}

	if val := vec.Pop(); val != nil {
		t.Errorf("Expected nil to be returned from popping empty vector, got %v", val)
	}
}
