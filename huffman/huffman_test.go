package huffman

import (
	"reflect"
	"testing"

	"github.com/mjjs/gompressor/vector"
)

func TestCreateFrequencyTable(t *testing.T) {
	bytes := vector.New()

	inputs := []struct {
		val  byte
		freq int
	}{
		{val: byte(1), freq: 2},
		{val: byte(4), freq: 6},
		{val: byte(2), freq: 12},
		{val: byte(3), freq: 1},
	}

	for _, input := range inputs {
		for i := 0; i < input.freq; i++ {
			bytes.Append(input.val)
		}
	}

	frequencies := createFrequencyTable(bytes)

	if n := frequencies.Size(); n != len(inputs) {
		t.Errorf("Expected %d, got %d", 4, n)
	}

	for _, input := range inputs {
		if frequency, ok := frequencies.Get(input.val); !ok {
			t.Errorf("Expected %x to be found in frequency table", input.val)
		} else if frequency != input.freq {
			t.Errorf("Expected %d, got %d", input.freq, frequency)
		}
	}
}

func TestIsLeafNode(t *testing.T) {
	leafA := &huffmanTreeNode{}
	leafB := &huffmanTreeNode{}
	innerNodeA := &huffmanTreeNode{left: leafA, right: leafB}
	innerNodeB := &huffmanTreeNode{left: leafA}

	if !isLeafNode(leafA) {
		t.Error("Expected a node with no children to be a leaf node")
	}

	if !isLeafNode(leafB) {
		t.Error("Expected a node with no children to be a leaf node")
	}

	if isLeafNode(innerNodeA) {
		t.Error("Expected a node with two children to be an inner node")
	}

	if isLeafNode(innerNodeB) {
		t.Error("Expected a node with one child to be an inner node")
	}
}

func TestBuildPrefixTree(t *testing.T) {
	original := "AABCDEF"
	originalBytes := vector.New()
	for _, c := range original {
		originalBytes.Append(byte(c))
	}

	frequencies := createFrequencyTable(originalBytes)
	prefixTree := buildPrefixTree(frequencies)

	if prefixTree.frequency != 7 {
		t.Errorf("Expected %d, got %d", 7, prefixTree.frequency)
	}

	if prefixTree.left.frequency != 3 {
		t.Errorf("Expected %d, got %d", 3, prefixTree.left.frequency)
	}

	if prefixTree.left.left.frequency != 1 {
		t.Errorf("Expected %d, got %d", 1, prefixTree.left.left.frequency)
	} else if prefixTree.left.left.value != byte('C') {
		t.Errorf("Expected %v, got %v", 'C', prefixTree.left.left.value)
	}

	if prefixTree.left.right.frequency != 2 {
		t.Errorf("Expected %d, got %d", 2, prefixTree.left.right.frequency)
	}

	if prefixTree.left.right.left.frequency != 1 {
		t.Errorf("Expected %d, got %d", 1, prefixTree.left.right.left.frequency)
	} else if prefixTree.left.right.left.value != byte('E') {
		t.Errorf("Expected %v, got %v", 'E', prefixTree.left.right.left.value)
	}

	if prefixTree.left.right.right.frequency != 1 {
		t.Errorf("Expected %d, got %d", 1, prefixTree.left.right.right.frequency)
	} else if prefixTree.left.right.right.value != byte('D') {
		t.Errorf("Expected %v, got %v", 'D', prefixTree.left.right.right.value)
	}

	if prefixTree.right.frequency != 4 {
		t.Errorf("Expected %d, got %d", 4, prefixTree.right.frequency)
	}

	if prefixTree.right.left.frequency != 2 {
		t.Errorf("Expected %d, got %d", 2, prefixTree.right.left.frequency)
	} else if prefixTree.right.left.value != byte('A') {
		t.Errorf("Expected %v, got %v", 'A', prefixTree.right.left.value)
	}

	if prefixTree.right.right.frequency != 2 {
		t.Errorf("Expected %d, got %d", 2, prefixTree.right.left.frequency)
	}

	if prefixTree.right.right.left.frequency != 1 {
		t.Errorf("Expected %d, got %d", 1, prefixTree.right.right.left.frequency)
	} else if prefixTree.right.right.left.value != byte('B') {
		t.Errorf("Expected %v, got %v", 'B', prefixTree.right.right.left.value)
	}

	if prefixTree.right.right.right.frequency != 1 {
		t.Errorf("Expected %d, got %d", 1, prefixTree.right.right.right.frequency)
	} else if prefixTree.right.right.right.value != byte('F') {
		t.Errorf("Expected %v, got %v", 'F', prefixTree.right.right.right.value)
	}
}

func TestDecompressedEqualsOriginal(t *testing.T) {
	original := []byte("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")
	origVector := vector.New()
	for _, c := range original {
		origVector.Append(c)
	}

	compressed := Compress(origVector)
	decompressed, err := Decompress(compressed)

	if err != nil {
		t.Errorf("Expected a nil error, got %s", err)
	}

	if !reflect.DeepEqual(origVector, decompressed) {
		t.Errorf("Expected '%s', got '%s'", origVector, decompressed)
	}
}

func TestDecompressReturnsErrorForInvalidData(t *testing.T) {
	_, err := Decompress(vector.New().AppendToCopy(byte(6), byte(5), byte(4), byte(3), byte(2), byte(1)))
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	_, err = Decompress(vector.New().AppendToCopy(
		// Bits to use from last byte
		byte(4),
		// Encoded prefix tree
		byte(0), byte(0), byte(1), byte('E'),
		byte(0), byte(1), byte(1), byte('s'),
		byte(1), byte(0), byte(1), byte('k'),
		byte(1), byte(1), byte(1), byte('o'),
		// Badly encoded huffman codes
		byte(6),
	))

	if err == nil {
		t.Error("Expected an error, got nil")
	}
}
