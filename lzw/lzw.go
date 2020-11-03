package lzw

import (
	"errors"
	"fmt"

	"github.com/mjjs/gompressor/bytevector"
)

const (
	initialDictSize uint16 = 256
	maxDictSize     uint16 = 65535
)

// ErrBadCompressedCode represents an error that occurs when the LZW decompression
// algorithm finds a code that is not valid for the assumed compression algorithm.
var ErrBadCompressedCode = errors.New("bad compression code")

// Compress takes a slice of uncompressed bytes as input and returns a slice of
// LZW codes that represent the compressed data.
func Compress(uncompressed *bytevector.Bytevector) ([]uint16, error) {
	createInitialDict := func() map[string]uint16 {
		dictionary := make(map[string]uint16, initialDictSize)

		for i := uint16(0); i < initialDictSize; i++ {
			dictionary[string([]byte{byte(i)})] = i
		}

		return dictionary
	}

	dictionary := createInitialDict()

	compressed := []uint16{}
	word := bytevector.New()

	for i := 0; i < uncompressed.Size(); i++ {
		if len(dictionary) == int(maxDictSize) {
			dictionary = createInitialDict()
		}

		byt := uncompressed.MustGet(i)
		newWord := word.AppendToCopy(byt)

		if _, ok := dictionary[newWord.String()]; ok {
			word = bytevector.New(uint(newWord.Size()))

			for j := 0; j < newWord.Size(); j++ {
				word.MustSet(j, newWord.MustGet(j))
			}
		} else {
			compressed = append(compressed, dictionary[word.String()])
			dictionary[newWord.String()] = uint16(len(dictionary))
			word = bytevector.New().AppendToCopy(byt)
		}
	}

	if word.Size() > 0 {
		compressed = append(compressed, dictionary[word.String()])
	}

	return compressed, nil
}

// Decompress takes in a slice of LZW codes representing some compressed data
// and outputs the decompressed data as a slice of bytes.
// An error is returned if the decompression algorithm finds a bad LZW code.
func Decompress(compressed []uint16) (*bytevector.Bytevector, error) {
	createInitialDict := func() map[uint16]*bytevector.Bytevector {
		dictionary := make(map[uint16]*bytevector.Bytevector, initialDictSize)

		for i := uint16(0); i < initialDictSize; i++ {
			bv := bytevector.New(1)
			bv.MustSet(0, byte(i))
			dictionary[i] = bv
		}

		return dictionary
	}

	dictionary := createInitialDict()

	result := bytevector.New()
	word := bytevector.New()

	for _, code := range compressed {
		if len(dictionary) == int(maxDictSize) {
			dictionary = createInitialDict()
		}

		entry := bytevector.New()

		if c, ok := dictionary[code]; ok {
			entry = bytevector.New(uint(c.Size()))
			for i := 0; i < c.Size(); i++ {
				entry.MustSet(i, c.MustGet(i))
			}
		} else if int(code) == len(dictionary) && word.Size() > 0 {
			entry = word.AppendToCopy(word.MustGet(0))
		} else {
			return nil, fmt.Errorf("%w: %d", ErrBadCompressedCode, code)
		}

		for i := 0; i < entry.Size(); i++ {
			result.Append(entry.MustGet(i))
		}

		if word.Size() > 0 {
			word = word.AppendToCopy(entry.MustGet(0))
			dictionary[uint16(len(dictionary))] = word
		}

		word = entry
	}

	return result, nil
}
