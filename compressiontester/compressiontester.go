package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mjjs/gompressor/algorithm/huffman"
	"github.com/mjjs/gompressor/algorithm/lzw"
	"github.com/mjjs/gompressor/datastructure/vector"
	"github.com/mjjs/gompressor/fileio"
)

type testResult struct {
	filename                   string
	algorithm                  string
	originalSizeBytes          int
	compressedSizeBytes        int
	compressRatio              float64
	compressTimeMicroseconds   int64
	decompressTimeMicroseconds int64
	success                    bool
	dictionarySize             uint16
}

const testFileNameA string = "../testdata/E.coli"
const testFileNameB string = "../testdata/world192.txt"
const testFileNameC string = "../testdata/Randomdata"

var fileNames = []string{
	testFileNameA, testFileNameB, testFileNameC,
}

var lzwDictSizes = []lzw.DictionarySize{
	lzw.XS, lzw.S, lzw.M, lzw.L, lzw.XL,
}

var testSizes = []int{
	1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152,
}

func main() {
	parallel := flag.Bool("parallel", false, "Run tests in parallel")
	flag.Parse()

	lzwResults := []testResult{}
	huffmanResults := []testResult{}

	wg := sync.WaitGroup{}

	for _, filepath := range fileNames {
		byteVector, err := readTestFile(filepath)
		if err != nil {
			panic(err)
		}

		for i := 0; i < len(testSizes)+1; i++ {
			var n int

			if i == len(testSizes) {
				n = byteVector.Size()
			} else {
				n = testSizes[i]
			}

			bytes := getNBytes(byteVector, n)

			filename := strings.Split(filepath, "/")[2]
			log.Printf("Running tests for %d bytes of file %s", bytes.Size(), filename)

			for _, dictSize := range lzwDictSizes {
				if *parallel {
					wg.Add(1)
					go func(ds lzw.DictionarySize) {
						result := testLZW(filename, ds, bytes)
						lzwResults = append(lzwResults, result)
						wg.Done()
					}(dictSize)
				} else {
					result := testLZW(filename, dictSize, bytes)
					lzwResults = append(lzwResults, result)
				}
			}

			if *parallel {
				wg.Add(1)
				go func() {
					huffmanResult := testHuffman(filename, bytes)
					huffmanResults = append(huffmanResults, huffmanResult)
					wg.Done()
				}()
			} else {
				huffmanResult := testHuffman(filename, bytes)
				huffmanResults = append(huffmanResults, huffmanResult)
			}
		}
	}

	if *parallel {
		wg.Wait()
	}

	writeCSV(lzwResults, "lzw.csv")
	writeCSV(huffmanResults, "huffman.csv")
}

func testLZW(filename string, dictSize lzw.DictionarySize, uncompressed *vector.Vector) testResult {
	originalSize := uncompressed.Size()

	result := testResult{
		algorithm:         "LZW",
		filename:          filename,
		dictionarySize:    uint16(dictSize),
		originalSizeBytes: originalSize,
	}

	log.Printf("Testing LZW compression with dictionary size of %d bytes", dictSize)
	compressStart := time.Now()

	compressed, err := lzw.CompressWithDictSize(uncompressed, dictSize)
	if err != nil {
		panic(fmt.Sprintf("lzw compression failed: %s", err))
	}

	result.compressTimeMicroseconds = time.Since(compressStart).Microseconds()
	result.compressedSizeBytes = compressed.Size() * 2 // Codes are 16 bits
	result.compressRatio = float64(result.compressedSizeBytes) / float64(originalSize) * 100

	decompressStart := time.Now()
	decompressed, err := lzw.Decompress(compressed)
	if err != nil {
		panic(fmt.Sprintf("lzw decompression failed: %s", err))
	}

	result.decompressTimeMicroseconds = time.Since(decompressStart).Microseconds()
	result.success = compare(uncompressed, decompressed)

	return result
}

func testHuffman(filename string, uncompressed *vector.Vector) testResult {
	log.Println("Testing Huffman compression")
	originalSize := uncompressed.Size()

	result := testResult{
		algorithm:         "Huffman",
		filename:          filename,
		originalSizeBytes: originalSize,
	}

	compressStart := time.Now()
	compressed := huffman.Compress(uncompressed)

	result.compressTimeMicroseconds = time.Since(compressStart).Microseconds()
	result.compressedSizeBytes = compressed.Size()
	result.compressRatio = float64(result.compressedSizeBytes) / float64(originalSize) * 100

	decompressStart := time.Now()
	decompressed, err := huffman.Decompress(compressed)
	if err != nil {
		panic(fmt.Sprintf("lzw decompression failed: %s", err))
	}

	result.decompressTimeMicroseconds = time.Since(decompressStart).Microseconds()
	result.success = compare(uncompressed, decompressed)

	return result
}

func readTestFile(fn string) (*vector.Vector, error) {
	return fileio.ReadFile(fn)
}

func getNBytes(from *vector.Vector, n int) *vector.Vector {
	if n > from.Size() {
		return from
	}

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

func writeCSV(results []testResult, name string) {
	headers := []string{"filename", "algorithm", "original size", "compressed size", "compression ratio", "compress time", "decompress time", "dictionary size"}
	records := [][]string{headers}

	for _, result := range results {
		record := []string{
			result.filename,
			result.algorithm,
			fmt.Sprintf("%d", result.originalSizeBytes),
			fmt.Sprintf("%d", result.compressedSizeBytes),
			fmt.Sprintf("%.2f", result.compressRatio),
			fmt.Sprintf("%d", result.compressTimeMicroseconds),
			fmt.Sprintf("%d", result.decompressTimeMicroseconds),
			fmt.Sprintf("%d", result.dictionarySize),
		}

		records = append(records, record)
	}

	fp, err := os.Create(name)
	if err != nil {
		panic(err)
	}

	defer fp.Close()

	writer := csv.NewWriter(fp)
	writer.WriteAll(records)

	if err := writer.Error(); err != nil {
		log.Fatalf("csv write failed: %s", err)
	}
}
