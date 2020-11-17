package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mjjs/gompressor/fileio"
	"github.com/mjjs/gompressor/huffman"
	"github.com/mjjs/gompressor/lzw"
	"github.com/mjjs/gompressor/vector"
)

// Load test file
//
// Run LZW tests:
// * Test files with different dictionary sizes
// * Max dictionary size, different file sizes
//
// Run Huffman tests:
// * Different file sizes
//
// For all tests, log:
// * Time to compress
// * Compression ratio
// * Time to decompress

const testFileNameA string = "../testdata/E.coli"
const testFileNameB string = "../testdata/world192.txt"
const testFileNameC string = "../testdata/Randomdata"

var fileNames = []string{
	testFileNameA, testFileNameB, testFileNameC,
}

var lzwDictSizes = []lzw.DictionarySize{
	lzw.XS, lzw.S, lzw.M, lzw.L, lzw.XL,
}

func main() {
	for _, filename := range fileNames {
		log.Printf("RUNNING TESTS FOR FILE %s", strings.Split(filename, "/")[2])
		byteVector, err := readTestFile(filename)
		if err != nil {
			panic(err)
		}

		if err := testLZW(byteVector); err != nil {
			panic(err)
		}

		if err := testHuffman(byteVector); err != nil {
			panic(err)
		}
	}
}

func testLZW(uncompressed *vector.Vector) error {
	originalSize := uncompressed.Size() * 8

	for _, dictSize := range lzwDictSizes {
		log.Printf("Testing LZW compression with dictionary size of %d bytes", dictSize)
		compressStart := time.Now()

		compressed, err := lzw.CompressWithDictSize(uncompressed, dictSize)
		if err != nil {
			return fmt.Errorf("lzw compression failed: %s", err)
		}

		compressDuration := time.Since(compressStart).Microseconds()
		log.Printf(
			"Compression took %d µs (%.2f s)",
			compressDuration,
			float64(compressDuration)/float64(1_000_000),
		)

		compressedSize := compressed.Size() * 8
		ratio := float64(compressedSize) / float64(originalSize)

		log.Printf(
			"Compressed size is: %d bytes, which is %.2f%% of original",
			compressedSize,
			ratio,
		)

		decompressStart := time.Now()
		decompressed, err := lzw.Decompress(compressed)
		if err != nil {
			return fmt.Errorf("lzw decompression failed: %s", err)
		}

		decompressDuration := time.Since(decompressStart).Microseconds()
		log.Printf(
			"Decompression took %d µs (%.2f s)",
			decompressDuration,
			float64(decompressDuration)/float64(1_000_000),
		)

		if !compare(uncompressed, decompressed) {
			log.Printf("FAILURE: Original and decompressed bytes do not match!")
		} else {
			log.Printf("SUCCESS: Original and compressed bytes match.")
		}
		fmt.Println()
	}
	return nil
}

func testHuffman(uncompressed *vector.Vector) error {
	originalSize := uncompressed.Size() * 8

	log.Println("Testing Huffman compression")

	compressStart := time.Now()

	compressed := huffman.Compress(uncompressed)

	compressDuration := time.Since(compressStart).Microseconds()
	log.Printf(
		"Compression took %d µs (%.2f s)",
		compressDuration,
		float64(compressDuration)/float64(1_000_000),
	)

	compressedSize := compressed.Size() * 8
	ratio := float64(compressedSize) / float64(originalSize)

	log.Printf(
		"Compressed size is: %d bytes, which is %.2f%% of original",
		compressedSize,
		ratio,
	)

	decompressStart := time.Now()
	decompressed, err := huffman.Decompress(compressed)
	if err != nil {
		return fmt.Errorf("lzw decompression failed: %s", err)
	}

	decompressDuration := time.Since(decompressStart).Microseconds()
	log.Printf(
		"Decompression took %d µs (%.2f s)",
		decompressDuration,
		float64(decompressDuration)/float64(1_000_000),
	)

	if !compare(uncompressed, decompressed) {
		log.Printf("FAILURE: Original and decompressed bytes do not match!")
	} else {
		log.Printf("SUCCESS: Original and compressed bytes match.")
	}
	fmt.Println()
	return nil
}

func readTestFile(fn string) (*vector.Vector, error) {
	return fileio.ReadFile(fn)
}

func getNBytes(from *vector.Vector, n int) *vector.Vector {
	newVector := vector.New(uint(n))
	for i := 0; i < n; i++ {
		newVector.MustSet(i, from.MustGet(i))
	}

	return newVector
}

func compare(a *vector.Vector, b *vector.Vector) bool {
	if a == nil || b == nil {
		return false
	}

	if a.Size() != b.Size() {
		return false
	}

	for i := 0; i < a.Size(); i++ {
		if a.MustGet(i) != b.MustGet(i) {
			return false
		}
	}

	return true
}
