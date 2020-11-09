package huffman

import (
	"fmt"

	"github.com/mjjs/gompressor/lzw/dictionary"
	"github.com/mjjs/gompressor/priorityqueue"
	"github.com/mjjs/gompressor/vector"
)

type node struct {
	frequency int
	value     byte
	left      *node
	right     *node
}

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

func Decompress(compressed *vector.Vector) *vector.Vector {
	prefixTree, nextIndex := decompressPrefixTree(compressed, 0)
	decompressed := vector.New()

	for nextIndex < compressed.Size() {
		nextIndex = decodeHuffmanCode(compressed, nextIndex, prefixTree, decompressed)
	}

	return decompressed
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

func buildPrefixTree(byteFrequencies *dictionary.Dictionary) *node {
	tree := new(priorityqueue.PriorityQueue)
	keys := byteFrequencies.Keys()

	for i := 0; i < keys.Size(); i++ {
		byt := keys.MustGet(i)
		frequency, _ := byteFrequencies.Get(byt)

		tree.Enqueue(frequency.(int), &node{frequency: frequency.(int), value: byt.(byte)})
	}

	for tree.Size() > 1 {
		aPrio, a := tree.Dequeue()
		bPrio, b := tree.Dequeue()

		newPrio := aPrio + bPrio

		node := &node{frequency: newPrio, left: a.(*node), right: b.(*node)}

		tree.Enqueue(newPrio, node)
	}

	_, root := tree.Dequeue()

	return root.(*node)
}

func buildCodes(root *node, str *vector.Vector, result *dictionary.Dictionary) {
	if root == nil {
		return
	}

	if isLeafNode(root) {
		result.Set(root.value, str)
	}

	buildCodes(root.left, str.AppendToCopy(byte(0)), result)
	buildCodes(root.right, str.AppendToCopy(byte(1)), result)
}

func compressPrefixTree(root *node, to *vector.Vector) {
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

func decompressPrefixTree(compressed *vector.Vector, index int) (*node, int) {
	byt := compressed.MustGet(index).(byte)
	switch byt {
	case 0x00:
		left, nextIndex := decompressPrefixTree(compressed, index+1)
		right, nextIndex := decompressPrefixTree(compressed, nextIndex)
		return &node{left: left, right: right}, nextIndex

	case 0x01:
		return &node{value: compressed.MustGet(index + 1).(byte)}, index + 2

	default:
		return nil, index + 1
	}
}

func decodeHuffmanCode(compressed *vector.Vector, index int, root *node, to *vector.Vector) int {
	if root == nil {
		return -1
	}

	if isLeafNode(root) {
		to.Append(root.value)
		return index
	}

	next := compressed.MustGet(index).(byte)
	switch next {
	case 0x00:
		return decodeHuffmanCode(compressed, index+1, root.left, to)
	case 0x01:
		return decodeHuffmanCode(compressed, index+1, root.right, to)
	default:
		panic(fmt.Sprintf(
			"Got %x when decoding a %d-length vector at index %d",
			next, compressed.Size(), index),
		)
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

func isLeafNode(n *node) bool {
	return n != nil && n.left == nil && n.right == nil
}
