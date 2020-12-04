// Package huffman implements the huffman coding algorithm which can be used
// for lossless data compression.
package huffman

import (
	"errors"
	"fmt"

	"github.com/mjjs/gompressor/datastructure/dictionary"
	"github.com/mjjs/gompressor/datastructure/priorityqueue"
	"github.com/mjjs/gompressor/datastructure/vector"
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

	codes := decompressHuffmanCodes(compressed, nextIndex, lastByteInBits)

	var err error
	nextIndex = 0

	for nextIndex < codes.Size() {
		nextIndex, err = decodeHuffmanCode(codes, nextIndex, prefixTree, decompressed)

		if err != nil {
			return nil, err
		}
	}

	return decompressed, nil
}

// createFrequencyTable takes in a vector of bytes and makes a frequency table
// indicating how often each byte appears in the vector.
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

// buildPrefixTree builds a tree-style data structure from a frequency table.
// The table priorizes bytes that have a higher frequency. This way the most
// often used bytes in the data get the shortest codeword when encoding.
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

// buildCodes traverses the prefix tree and builds a code for each unique byte found in
// the original, uncompressed data. Each code consists of a vector of 0s and 1s.
// A 0 indicates taking the left child of the current node and 1 the right one.
// The codes are stored in the result dictionary, which maps each byte to the code.
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

// compressPrefixTree takes in the root of a prefix tree and the output vector to.
// Leaf nodes are encoded as byte 1 and the value it represents, other nodes
// are encoded as byte 0.
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

// decompressPrefixTree goes through the compressed vector starting from index
// and recreates the prefix tree from the encoded data.
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

// decompressHuffmanCodes takes in compressed bytes and an index where to start
// decompressing. As all 0/1 bytes have been encoded as bits, lastByteInBits
// indicates how many bits to read from the last byte.
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

// decodeHuffmanCode reads a huffman code from codes at index and writes it into to.
// Returns the index where to start reading the next code.
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

// encodeToHuffmanCodes goes through each byte in the uncompressed data and returns
// a vector where each byte has been replaced with the huffman code it represents.
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

// compressHuffmanCodes takes in huffman codes represented as a vector of 1 and 0
// bytes. Since the bytes can only be 0 and 1, they can be compressed as bits
// so that 8 bytes can be encoded as one byte.
func compressHuffmanCodes(codes *vector.Vector) (compressedCodes *vector.Vector, lastByteInBits int) {
	currentCode := vector.New(0, 8)
	encodedCode := byte(0)
	totalBits := 0

	compressedCodes = vector.New()

	for i := 0; i < codes.Size(); i++ {
		currentCode.Append(codes.MustGet(i))

		if currentCode.Size() != 8 && i != codes.Size()-1 {
			continue
		}

		for j := 0; j < currentCode.Size(); j++ {
			totalBits++

			encodedCode <<= 1

			if currentCode.MustGet(j) != byte(0) {
				encodedCode |= 1
			}
		}

		compressedCodes.Append(encodedCode)
		currentCode = vector.New(0, 8)
		encodedCode = byte(0)
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
