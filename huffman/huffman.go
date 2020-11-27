// Package huffman implements the huffman coding algorithm which can be used
// for lossless data compression.
package huffman

import (
	"errors"
	"fmt"

	"github.com/mjjs/gompressor/dictionary"
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
	byteFrequencies := createFrequencyTable(uncompressed)
	prefixTree := buildPrefixTree(byteFrequencies)

	codes := dictionary.New()
	buildCodes(prefixTree, vector.New(), codes)

	compressedPrefixTree := vector.New()
	compressPrefixTree(prefixTree, compressedPrefixTree)

	encodedBytes := encodeToHuffmanCodes(uncompressed, codes)

	compressedCodes, lastByteInBits := compressHuffmanCodes(encodedBytes)

	// Reserve space for the last byte size, prefix tree and huffman codes
	compressed := vector.New(0, 1+uint(compressedPrefixTree.Size()+compressedCodes.Size()))
	compressed.Append(byte(lastByteInBits))

	for i := 0; i < compressedPrefixTree.Size(); i++ {
		compressed.Append(compressedPrefixTree.MustGet(i))
	}

	for i := 0; i < compressedCodes.Size(); i++ {
		compressed.Append(compressedCodes.MustGet(i))
	}

	return compressed
}

// Decompress takes in a vector of huffman compressed bytes and outputs a vector
// of uncompressed bytes. Returns a non-nil error if the decompression fails.
func Decompress(compressed *vector.Vector) (*vector.Vector, error) {
	if compressed.Size() == 0 {
		return compressed, nil
	}

	lastByteInBits := int(compressed.MustGet(0).(byte))

	prefixTree, nextIndex := decompressPrefixTree(compressed, 1)
	decompressed := vector.New()

	var err error

	codes := decompressHuffmanCodes(compressed, nextIndex, lastByteInBits)

	nextIndex = 0

	for nextIndex < codes.Size() {
		nextIndex, err = decodeHuffmanCode(codes, nextIndex, prefixTree, decompressed)

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

func buildCodes(root *huffmanTreeNode, code *vector.Vector, result *dictionary.Dictionary) {
	if root == nil {
		return
	}

	if isLeafNode(root) {
		result.Set(root.value, code)
	}

	buildCodes(root.left, code.AppendToCopy(byte(0)), result)
	buildCodes(root.right, code.AppendToCopy(byte(1)), result)
}

func compressPrefixTree(root *huffmanTreeNode, to *vector.Vector) {
	switch isLeafNode(root) {
	case true:
		to.Append(byte(1))
		to.Append(root.value)
	case false:
		to.Append(byte(0))
		compressPrefixTree(root.left, to)
		compressPrefixTree(root.right, to)
	}
}

func decompressPrefixTree(compressed *vector.Vector, index int) (*huffmanTreeNode, int) {
	byt := compressed.MustGet(index).(byte)
	switch byt {
	case byte(0):
		left, nextIndex := decompressPrefixTree(compressed, index+1)
		right, nextIndex := decompressPrefixTree(compressed, nextIndex)
		return &huffmanTreeNode{left: left, right: right}, nextIndex

	case byte(1):
		return &huffmanTreeNode{value: compressed.MustGet(index + 1).(byte)}, index + 2

	default:
		return nil, index + 1
	}
}

func decompressHuffmanCodes(compressed *vector.Vector, index int, lastByteInBits int) *vector.Vector {
	huffmanCodes := vector.New(0, uint(compressed.Size()-index))

	for i := index; i < compressed.Size(); i++ {
		codeByte := compressed.MustGet(i).(byte)

		totalBits := 7
		if i == compressed.Size()-1 {
			totalBits = lastByteInBits - 1
		}

		for j := totalBits; j >= 0; j-- {
			huffmanCodes.Append((codeByte >> j) & 1)
		}
	}

	return huffmanCodes
}

func decodeHuffmanCode(codes *vector.Vector, index int, root *huffmanTreeNode, to *vector.Vector) (int, error) {
	if root == nil {
		return 0, errors.New("No prefix tree supplied")
	}

	if isLeafNode(root) {
		to.Append(root.value)
		return index, nil
	}

	next := codes.MustGet(index)
	switch next {
	case byte(0):
		return decodeHuffmanCode(codes, index+1, root.left, to)
	case byte(1):
		return decodeHuffmanCode(codes, index+1, root.right, to)
	default:
		return 0, fmt.Errorf("An unexpected symbol %x found in the compressed data", next)
	}
}

func encodeToHuffmanCodes(uncompressed *vector.Vector, codes *dictionary.Dictionary) *vector.Vector {
	encodedHuffmanCodes := vector.New()

	for i := 0; i < uncompressed.Size(); i++ {
		byt := uncompressed.MustGet(i)

		iCode, _ := codes.Get(byt)
		code := iCode.(*vector.Vector)

		for j := 0; j < code.Size(); j++ {
			encodedHuffmanCodes.Append(code.MustGet(j))
		}
	}

	return encodedHuffmanCodes
}

func compressHuffmanCodes(codes *vector.Vector) (compressedCodes *vector.Vector, lastByteInBits int) {
	currentCode := vector.New(0, 8)
	encodedCode := byte(0)
	totalBits := 0

	compressedCodes = vector.New()

	for i := 0; i < codes.Size(); i++ {
		currentCode.Append(codes.MustGet(i))

		if currentCode.Size() == 8 || i == codes.Size()-1 {
			for j := 0; j < currentCode.Size(); j++ {
				totalBits++

				if currentCode.MustGet(j) == byte(0) {
					encodedCode <<= 1
				} else {
					encodedCode <<= 1
					encodedCode |= 1
				}
			}

			compressedCodes.Append(encodedCode)
			currentCode = vector.New(0, 8)
			encodedCode = byte(0)
		}
	}

	lastByteInBits = totalBits % 8
	if lastByteInBits == 0 {
		lastByteInBits = 8
	}

	return compressedCodes, lastByteInBits
}

func isLeafNode(n *huffmanTreeNode) bool {
	return n != nil && n.left == nil && n.right == nil
}
