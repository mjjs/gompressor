package lzw

import (
	"errors"
	"fmt"

	"github.com/mjjs/gompressor/bytevector"
	"github.com/mjjs/gompressor/lzw/dictionary"
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
	createInitialDict := func() *dictionary.Dictionary {
		dict := dictionary.NewWithSize(uint(initialDictSize))

		for i := uint16(0); i < initialDictSize; i++ {
			dict.Set(string([]byte{byte(i)}), i)
		}

		return dict
	}

	dict := createInitialDict()

	compressed := []uint16{}
	word := bytevector.New()

	for i := 0; i < uncompressed.Size(); i++ {
		if dict.Size() == int(maxDictSize) {
			dict = createInitialDict()
		}

		byt := uncompressed.MustGet(i)
		newWord := word.AppendToCopy(byt)

		if _, ok := dict.Get(newWord.String()); ok {
			word = bytevector.New(uint(newWord.Size()))

			for j := 0; j < newWord.Size(); j++ {
				word.MustSet(j, newWord.MustGet(j))
			}
		} else {
			code, _ := dict.Get(word.String())
			compressed = append(compressed, code.(uint16))

			dict.Set(newWord.String(), uint16(dict.Size()))
			word = bytevector.New().AppendToCopy(byt)
		}
	}

	if word.Size() > 0 {
		code, _ := dict.Get(word.String())
		compressed = append(compressed, code.(uint16))
	}

	return compressed, nil
}

// Decompress takes in a slice of LZW codes representing some compressed data
// and outputs the decompressed data as a slice of bytes.
// An error is returned if the decompression algorithm finds a bad LZW code.
func Decompress(compressed []uint16) (*bytevector.Bytevector, error) {
	createInitialDict := func() *dictionary.Dictionary {
		dict := dictionary.NewWithSize(uint(initialDictSize))

		for i := uint16(0); i < initialDictSize; i++ {
			bv := bytevector.New(1)
			bv.MustSet(0, byte(i))
			dict.Set(i, bv)
		}

		return dict
	}

	dictionary := createInitialDict()

	result := bytevector.New()
	word := bytevector.New()

	for _, code := range compressed {
		if dictionary.Size() == int(maxDictSize) {
			dictionary = createInitialDict()
		}

		entry := bytevector.New()

		if c, ok := dictionary.Get(code); ok {
			vector := c.(*bytevector.Bytevector)

			entry = bytevector.New(uint(vector.Size()))
			for i := 0; i < vector.Size(); i++ {
				entry.MustSet(i, vector.MustGet(i))
			}
		} else if int(code) == dictionary.Size() && word.Size() > 0 {
			entry = word.AppendToCopy(word.MustGet(0))
		} else {
			return nil, fmt.Errorf("%w: %d", ErrBadCompressedCode, code)
		}

		for i := 0; i < entry.Size(); i++ {
			result.Append(entry.MustGet(i))
		}

		if word.Size() > 0 {
			word = word.AppendToCopy(entry.MustGet(0))
			dictionary.Set(uint16(dictionary.Size()), word)
		}

		word = entry
	}

	return result, nil
}
