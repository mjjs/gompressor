package main

import (
	"flag"
	"log"

	"github.com/mjjs/gompressor/algorithm/huffman"
	"github.com/mjjs/gompressor/algorithm/lzw"
	"github.com/mjjs/gompressor/fileio"
	"github.com/mjjs/gompressor/ui"
)

func main() {
	tuiFlag := flag.Bool("tui", false, "use TUI")
	compressFlag := flag.Bool("compress", false, "compress the input file")
	decompressFlag := flag.Bool("decompress", false, "decompress the input file")
	huffmanFlag := flag.Bool("huffman", false, "use huffman algorithm")
	lzwFlag := flag.Bool("lzw", false, "use lzw algorithm")
	inputFileFlag := flag.String("in", "", "input file")
	outputFileFlag := flag.String("out", "", "output file")

	flag.Parse()

	if *tuiFlag {
		ui.New().Run()
	}

	if *compressFlag && *decompressFlag {
		log.Fatal("Only compress or decompress flag can be provided")
	} else if !*compressFlag && !*decompressFlag {
		log.Fatal("Compress or decompress flag must be provided")
	}

	if *inputFileFlag == "" || *outputFileFlag == "" {
		log.Fatal("Input and output files must be provided")
	}

	if *huffmanFlag && *lzwFlag {
		log.Fatal("Only supply one of the algorithm flags")
	} else if !*huffmanFlag && !*lzwFlag {
		log.Fatal("Supply one of the algorithm flags (-huffman, -lzw)")
	}

	if *compressFlag {
		if *huffmanFlag {
			compressHuffman(*inputFileFlag, *outputFileFlag)
		} else {
			compressLZW(*inputFileFlag, *outputFileFlag)
		}
	} else {
		if *huffmanFlag {
			decompressHuffman(*inputFileFlag, *outputFileFlag)
		} else {
			decompressLZW(*inputFileFlag, *outputFileFlag)
		}
	}
}

func compressHuffman(inputFilename string, outputFilename string) {
	bytes, err := fileio.ReadFile(inputFilename)
	if err != nil {
		log.Fatalf("Input file could not be read: %s", err)
	}

	compressed := huffman.Compress(bytes)

	err = fileio.WriteFile(compressed, outputFilename)
	if err != nil {
		log.Fatalf("Could not write compressed data: %s", err)
	}

	log.Printf("Wrote %d bytes to %s", compressed.Size(), outputFilename)
}

func decompressHuffman(inputFilename string, outputFilename string) {
	bytes, err := fileio.ReadFile(inputFilename)
	if err != nil {
		log.Fatalf("Input file could not be read: %s", err)
	}

	decompressed, err := huffman.Decompress(bytes)
	if err != nil {
		log.Fatalf("Could not decompress data: %s", err)
	}

	err = fileio.WriteFile(decompressed, outputFilename)
	if err != nil {
		log.Fatalf("Could not write decompressed data: %s", err)
	}

	log.Printf("Wrote %d bytes to %s", decompressed.Size(), outputFilename)
}

func compressLZW(inputFilename string, outputFilename string) {
	bytes, err := fileio.ReadFile(inputFilename)
	if err != nil {
		log.Fatalf("Input file could not be read: %s", err)
	}

	compressed, err := lzw.Compress(bytes)
	if err != nil {
		log.Fatalf("Could not compress data: %s", err)
	}

	err = fileio.WriteLZWFile(compressed, outputFilename)
	if err != nil {
		log.Fatalf("Could not write compressed data: %s", err)
	}

	log.Printf("Wrote %d bytes to %s", compressed.Size(), outputFilename)
}

func decompressLZW(inputFilename string, outputFilename string) {
	bytes, err := fileio.ReadLZWFile(inputFilename)
	if err != nil {
		log.Fatalf("Input file could not be read: %s", err)
	}

	decompressed, err := lzw.Decompress(bytes)
	if err != nil {
		log.Fatalf("Could not decompress data: %s", err)
	}

	err = fileio.WriteFile(decompressed, outputFilename)
	if err != nil {
		log.Fatalf("Could not write decompressed data: %s", err)
	}

	log.Printf("Wrote %d bytes to %s", decompressed.Size(), outputFilename)
}
