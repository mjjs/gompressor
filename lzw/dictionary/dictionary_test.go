package dictionary

import (
	"testing"
)

func TestNewStartsWithInitialSize(t *testing.T) {
	dictionary := New()
	if len(dictionary.buckets) != int(defaultSize) {
		t.Errorf("Expected initial size to be %d, got %d", defaultSize, len(dictionary.buckets))
	}
}

func TestNewStartsWithNonNilBuckets(t *testing.T) {
	dictionary := New()

	for _, bucket := range dictionary.buckets {
		if bucket == nil {
			t.Error("Expected buckets to be initialized, got nil bucket")
		}
	}
}

func TestCanSetAndGet(t *testing.T) {
	dict := New()
	key := "key"
	var expected uint16 = 123

	if _, exists := dict.Get(key); exists {
		t.Errorf("Expected %t, got %t", false, exists)
	}

	dict.Set(key, expected)

	if actual, exists := dict.Get(key); !exists {
		t.Errorf("Expected %t, got %t", true, exists)
	} else if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSetUpdatesExistingValue(t *testing.T) {
	dict := New()
	key := "key"
	var expected uint16 = 123

	dict.Set(key, expected)

	if actual, exists := dict.Get(key); !exists {
		t.Errorf("Expected %t, got %t", true, exists)
	} else if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	var newExpected uint16 = 212
	dict.Set(key, newExpected)

	if actual, exists := dict.Get(key); !exists {
		t.Errorf("Expected %t, got %t", true, exists)
	} else if actual != newExpected {
		t.Errorf("Expected %v, got %v", newExpected, actual)
	}
}

func TestCanRemoveValues(t *testing.T) {
	dict := New()
	key := "key"
	var expected uint16 = 123

	dict.Set(key, expected)

	if actual, exists := dict.Get(key); !exists {
		t.Errorf("Expected %t, got %t", true, exists)
	} else if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	dict.Remove(key)

	if _, exists := dict.Get(key); exists {
		t.Errorf("Expected %t, got %t", false, exists)
	}
}

func TestCanGrowPastSize(t *testing.T) {
	dict := NewWithSize(1)
	var (
		key1 string = "key1"
		key2 string = "key2"
		val  uint16 = 123
		val2 uint16 = 212
	)

	dict.Set(key1, val)
	dict.Set(key2, val2)
}

func TestElementIsFoundAfterGrow(t *testing.T) {
	dict := NewWithSize(1)
	key := "key"
	val := uint16(15)

	dict.Set(key, val)
	dict.Set("a", 123)
	dict.Set("b", 121)

	if v, ok := dict.Get(key); !ok {
		t.Errorf("Expected %t got %t", true, ok)
	} else if v != val {
		t.Errorf("Expected %d got %d", val, v)
	}
}

func TestReturnsSizeCorrectly(t *testing.T) {
	dict := New()
	if n := dict.Size(); n != 0 {
		t.Errorf("Expected size to be %d, got %d", 0, n)
	}

	dict.Set("key", 123)
	if n := dict.Size(); n != 1 {
		t.Errorf("Expected size to be %d, got %d", 1, n)
	}
}
