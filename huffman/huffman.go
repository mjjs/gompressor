// Package huffman implements the huffman coding algorithm which can be used
// for compressing data.
package huffman

import (
	"errors"
	"fmt"

	"github.com/mjjs/gompressor/lzw/dictionary"
	"github.com/mjjs/gompressor/priorityqueue"
	"github.com/mjjs/gompressor/vector"
)

type huffmanTreeNode struct {
	frequency int
	value     byte
	left      *huffmanTreeNode
	right     *huffmanTreeNode
}

// Compress takes in a vector of uncompressed bytes and outputs a vector of
// compressed bytes.
func Compress(uncompressed *vector.Vector) *vector.Vector {
	frequencies := createFrequencyTable(uncompressed)
	prefixTree := buildPrefixTree(frequencies)

	codes := dictionary.New()
	buildCodes(prefixTree, vector.New(), codes)

	compressed := vector.New()
	compressPrefixTree(prefixTree, compressed)
	compressHuffmanCodes(uncompressed, codes, compressed)

	return compressed
}

// Decompress takes in a vector of huffman compressed bytes and outputs a vector
// of uncompressed bytes. Returns a non-nil error if the decompression fails.
func Decompress(compressed *vector.Vector) (*vector.Vector, error) {
	prefixTree, nextIndex := decompressPrefixTree(compressed, 0)
	decompressed := vector.New()

	var err error

	for nextIndex < compressed.Size() {
		nextIndex, err = decodeHuffmanCode(compressed, nextIndex, prefixTree, decompressed)

		if err != nil {
			return nil, err
		}
	}

	return decompressed, nil
}

func createFrequencyTable(bytes *vector.Vector) *dictionary.Dictionary {
	dict := dictionary.New()
	for i := 0; i < bytes.Size(); i++ {
		byt := bytes.MustGet(i)

		if frequency, exists := dict.Get(byt); !exists {
			dict.Set(byt, 1)
		} else {
			dict.Set(byt, frequency.(int)+1)
		}
	}

	return dict
}

func buildPrefixTree(byteFrequencies *dictionary.Dictionary) *huffmanTreeNode {
	tree := new(priorityqueue.PriorityQueue)
	keys := byteFrequencies.Keys()

	for i := 0; i < keys.Size(); i++ {
		byt := keys.MustGet(i)
		frequency, _ := byteFrequencies.Get(byt)

		tree.Enqueue(frequency.(int), &huffmanTreeNode{frequency: frequency.(int), value: byt.(byte)})
	}

	for tree.Size() > 1 {
		aPrio, a := tree.Dequeue()
		bPrio, b := tree.Dequeue()

		newPrio := aPrio + bPrio

		node := &huffmanTreeNode{frequency: newPrio, left: a.(*huffmanTreeNode), right: b.(*huffmanTreeNode)}

		tree.Enqueue(newPrio, node)
	}

	_, root := tree.Dequeue()

	return root.(*huffmanTreeNode)
}

func buildCodes(root *huffmanTreeNode, str *vector.Vector, result *dictionary.Dictionary) {
	if root == nil {
		return
	}

	if isLeafNode(root) {
		result.Set(root.value, str)
	}

	buildCodes(root.left, str.AppendToCopy(byte(0)), result)
	buildCodes(root.right, str.AppendToCopy(byte(1)), result)
}

func compressPrefixTree(root *huffmanTreeNode, to *vector.Vector) {
	switch isLeafNode(root) {
	case true:
		// This can be optimzied to 1 bit if needed
		to.Append(byte(1))
		to.Append(root.value)
	case false:
		// This can be optimzied to 1 bit if needed
		to.Append(byte(0))
		compressPrefixTree(root.left, to)
		compressPrefixTree(root.right, to)
	}
}

func decompressPrefixTree(compressed *vector.Vector, index int) (*huffmanTreeNode, int) {
	byt := compressed.MustGet(index).(byte)
	switch byt {
	case 0x00:
		left, nextIndex := decompressPrefixTree(compressed, index+1)
		right, nextIndex := decompressPrefixTree(compressed, nextIndex)
		return &huffmanTreeNode{left: left, right: right}, nextIndex

	case 0x01:
		return &huffmanTreeNode{value: compressed.MustGet(index + 1).(byte)}, index + 2

	default:
		return nil, index + 1
	}
}

func decodeHuffmanCode(compressed *vector.Vector, index int, root *huffmanTreeNode, to *vector.Vector) (int, error) {
	if root == nil {
		return 0, errors.New("No prefix tree supplied")
	}

	if isLeafNode(root) {
		to.Append(root.value)
		return index, nil
	}

	next := compressed.MustGet(index).(byte)
	switch next {
	case 0x00:
		return decodeHuffmanCode(compressed, index+1, root.left, to)
	case 0x01:
		return decodeHuffmanCode(compressed, index+1, root.right, to)
	default:
		return 0, fmt.Errorf("An unexpected symbol %x found in the compressed data", next)
	}
}

func compressHuffmanCodes(uncompressed *vector.Vector, codes *dictionary.Dictionary, to *vector.Vector) {
	for i := 0; i < uncompressed.Size(); i++ {
		value := uncompressed.MustGet(i)
		iCode, _ := codes.Get(value)
		code := iCode.(*vector.Vector)

		for j := 0; j < code.Size(); j++ {
			to.Append(code.MustGet(j))
		}
	}
}

func isLeafNode(n *huffmanTreeNode) bool {
	return n != nil && n.left == nil && n.right == nil
}
