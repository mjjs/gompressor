package lzw

import (
	"errors"
	"fmt"
	"strings"
)

const (
	initialDictSize uint16 = 256
	maxDictSize     uint16 = 65535
)

var ErrBadCompressedCode = errors.New("bad compression code")

func Compress(uncompressed []byte) ([]uint16, error) {
	createInitialDict := func() map[string]uint16 {
		dictionary := make(map[string]uint16, initialDictSize)

		for i := uint16(0); i < initialDictSize; i++ {
			dictionary[string([]byte{byte(i)})] = i
		}

		return dictionary
	}

	dictionary := createInitialDict()

	compressed := []uint16{}
	word := []byte{}

	for _, byt := range uncompressed {
		if len(dictionary) == int(maxDictSize) {
			dictionary = createInitialDict()
		}

		newWord := append(word, byt)

		if _, ok := dictionary[string(newWord)]; ok {
			word = newWord
		} else {
			compressed = append(compressed, dictionary[string(word)])
			dictionary[string(newWord)] = uint16(len(dictionary))
			word = []byte{byt}
		}
	}

	if len(word) > 0 {
		compressed = append(compressed, dictionary[string(word)])
	}

	return compressed, nil
}

func Decompress(compressed []uint16) ([]byte, error) {
	createInitialDict := func() map[uint16][]byte {
		dictionary := make(map[uint16][]byte, initialDictSize)

		for i := uint16(0); i < initialDictSize; i++ {
			dictionary[i] = []byte{byte(i)}
		}

		return dictionary
	}

	dictionary := createInitialDict()

	result := strings.Builder{}
	word := []byte{}

	for _, code := range compressed {
		if len(dictionary) == int(maxDictSize) {
			dictionary = createInitialDict()
		}

		entry := []byte{}

		if c, ok := dictionary[code]; ok {
			entry = c[:len(c):len(c)]
		} else if int(code) == len(dictionary) && len(word) > 0 {
			entry = append(word, word[0])
		} else {
			return nil, fmt.Errorf("%w: %d", ErrBadCompressedCode, code)
		}

		result.Write(entry)

		if len(word) > 0 {
			word = append(word, entry[0])
			dictionary[uint16(len(dictionary))] = word
		}

		word = entry
	}

	return []byte(result.String()), nil
}
