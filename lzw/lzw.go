package lzw

import (
	"errors"
	"fmt"

	"github.com/mjjs/gompressor/lzw/dictionary"
	"github.com/mjjs/gompressor/vector"
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
func Compress(uncompressed *vector.Vector) ([]uint16, error) {
	createInitialDict := func() *dictionary.Dictionary {
		dict := dictionary.NewWithSize(uint(initialDictSize))

		for i := uint16(0); i < initialDictSize; i++ {
			dict.Set(string([]byte{byte(i)}), i)
		}

		return dict
	}

	dict := createInitialDict()

	compressed := []uint16{}
	word := vector.New()

	for i := 0; i < uncompressed.Size(); i++ {
		if dict.Size() == int(maxDictSize) {
			dict = createInitialDict()
		}

		byt := uncompressed.MustGet(i)
		newWord := word.AppendToCopy(byt)

		if _, ok := dict.Get(newWord.String()); ok {
			word = vector.New(uint(newWord.Size()))

			for j := 0; j < newWord.Size(); j++ {
				word.MustSet(j, newWord.MustGet(j))
			}
		} else {
			code, _ := dict.Get(word.String())
			compressed = append(compressed, code.(uint16))

			dict.Set(newWord.String(), uint16(dict.Size()))
			word = vector.New().AppendToCopy(byt)
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
func Decompress(compressed []uint16) (*vector.Vector, error) {
	createInitialDict := func() *dictionary.Dictionary {
		dict := dictionary.NewWithSize(uint(initialDictSize))

		for i := uint16(0); i < initialDictSize; i++ {
			bv := vector.New(1)
			bv.MustSet(0, byte(i))
			dict.Set(i, bv)
		}

		return dict
	}

	dictionary := createInitialDict()

	result := vector.New()
	word := vector.New()

	for _, code := range compressed {
		if dictionary.Size() == int(maxDictSize) {
			dictionary = createInitialDict()
		}

		entry := vector.New()

		if c, ok := dictionary.Get(code); ok {
			byteVector := c.(*vector.Vector)

			entry = vector.New(uint(byteVector.Size()))
			for i := 0; i < byteVector.Size(); i++ {
				entry.MustSet(i, byteVector.MustGet(i))
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
