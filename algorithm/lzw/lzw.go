// Package lzw implements the Lempel-Ziv-Welch lossless data compression algorithm.
package lzw

import (
	"errors"
	"fmt"

	"github.com/mjjs/gompressor/datastructure/dictionary"
	"github.com/mjjs/gompressor/datastructure/vector"
)

// DictionarySize determines how large the dictionary used in compression can
// grow before needing to be reset. Larger values result in more efficient
// compression.
type DictionarySize uint16

// Dictionary sizes
const (
	XS DictionarySize = 512
	S                 = 1023
	M                 = 4095
	L                 = 32767
	XL                = 65535
)

const initialDictSize uint16 = 255

// ErrBadCompressedCode represents an error that occurs when the LZW decompression
// algorithm finds a code that is not valid for the assumed compression algorithm.
var ErrBadCompressedCode = errors.New("bad compression code")

// ErrInvalidDictionarySize represents an error indicating usage of an invalid
// dictionary size. This can happen when attempting to compress or decompress
// data.
var ErrInvalidDictionarySize = errors.New("invalid dictionary size")

// CompressWithDictSize takes a slice of uncompressed bytes and a dictionary size
// as input and returns a slice of LZW codes that represent the compressed data.
// This is mostly a utility function for testing how the dictionary size changes
// the compression level.
func CompressWithDictSize(uncompressed *vector.Vector, size DictionarySize) (*vector.Vector, error) {
	if !isValidDictionarySize(size) {
		return nil, fmt.Errorf("%w: %d", ErrInvalidDictionarySize, int(size))
	}

	if uncompressed.Size() == 0 {
		return uncompressed, nil
	}

	dict := createInitialCompressDictionary()

	compressed := vector.New(0, uint(uncompressed.Size()))
	compressed.Append(uint16(size))

	word := vector.New()

	for i := 0; i < uncompressed.Size(); i++ {
		if dict.Size() == int(size) {
			dict = createInitialCompressDictionary()
		}

		byt := uncompressed.MustGet(i)
		newWord := word.AppendToCopy(byt)

		if _, ok := dict.Get(newWord.String()); ok {
			word = newWord
		} else {
			code, _ := dict.Get(word.String())
			compressed.Append(code.(uint16))

			dict.Set(newWord.String(), uint16(dict.Size()))
			word = vector.New().AppendToCopy(byt)
		}
	}

	if word.Size() > 0 {
		code, _ := dict.Get(word.String())
		compressed.Append(code.(uint16))
	}

	return compressed, nil
}

// Compress is a shortcut for compressing with the largest dictionary size.
func Compress(uncompressed *vector.Vector) (*vector.Vector, error) {
	return CompressWithDictSize(uncompressed, XL)
}

// Decompress takes in a slice of LZW codes representing some compressed data
// and outputs the decompressed data as a slice of bytes.
// An error is returned if the decompression algorithm finds a bad LZW code.
func Decompress(compressed *vector.Vector) (*vector.Vector, error) {
	if compressed.Size() == 0 {
		return compressed, nil
	}

	size := compressed.MustGet(0).(uint16)
	if !isValidDictionarySize(DictionarySize(size)) {
		return nil, fmt.Errorf("the data is compressed with an invalid dictionary size %d", size)
	}

	dict := createInitialDecompressDictionary()

	result := vector.New()
	word := vector.New()

	for i := 1; i < compressed.Size(); i++ {
		if dict.Size() == int(size) {
			dict = createInitialDecompressDictionary()
		}

		code := compressed.MustGet(i)

		entry := vector.New()

		if c, ok := dict.Get(code); ok {
			byteVector := c.(*vector.Vector)

			entry = vector.New(uint(byteVector.Size()))
			for i := 0; i < byteVector.Size(); i++ {
				entry.MustSet(i, byteVector.MustGet(i))
			}
		} else if int(code.(uint16)) == dict.Size() && word.Size() > 0 {
			entry = word.AppendToCopy(word.MustGet(0))
		} else {
			return nil, fmt.Errorf("%w: %d", ErrBadCompressedCode, code)
		}

		for i := 0; i < entry.Size(); i++ {
			result.Append(entry.MustGet(i))
		}

		if word.Size() > 0 {
			word = word.AppendToCopy(entry.MustGet(0))
			dict.Set(uint16(dict.Size()), word)
		}

		word = entry
	}

	return result, nil
}

func createInitialCompressDictionary() *dictionary.Dictionary {
	dict := dictionary.NewWithSize(uint(initialDictSize))

	for i := uint16(0); i <= initialDictSize; i++ {
		dict.Set(string([]byte{byte(i)}), i)
	}

	return dict
}

func createInitialDecompressDictionary() *dictionary.Dictionary {
	dict := dictionary.NewWithSize(uint(initialDictSize))

	for i := uint16(0); i <= initialDictSize; i++ {
		bv := vector.New(1)
		bv.MustSet(0, byte(i))
		dict.Set(i, bv)
	}

	return dict
}

func isValidDictionarySize(size DictionarySize) bool {
	validSizes := []DictionarySize{
		XS, S, M, L, XL,
	}

	for _, dictionarySize := range validSizes {
		if dictionarySize == size {
			return true
		}
	}

	return false
}
