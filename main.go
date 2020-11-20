package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mjjs/gompressor/fileio"
	"github.com/mjjs/gompressor/huffman"
	"github.com/mjjs/gompressor/lzw"
	"github.com/mjjs/gompressor/vector"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	compressFlag := flag.Bool("compress", false, "Compress the given file")
	decompressFlag := flag.Bool("decompress", false, "Decompress the given file")
	inFlag := flag.String("in", "", "input file")
	outFlag := flag.String("out", "", "output file")
	algorithm := flag.String("algorithm", "", "lzw or huffman")

	flag.Parse()

	if !*compressFlag && !*decompressFlag {
		log.Fatal("No -compress or -decompress flag given")
	} else if *compressFlag && *decompressFlag {
		log.Fatal("Cannot use -compress and -decompress flag at the same time")
	} else if *algorithm == "" {
		log.Fatal("No algorithm specified")
	} else if *algorithm != "huffman" && *algorithm != "lzw" {
		log.Fatal("Algorithm can be only \"lzw\" or \"huffman\"")
	} else if len(*inFlag) == 0 || len(*outFlag) == 0 {
		log.Fatal("Input or output file flag missing")
	}

	if *compressFlag {
		compress(*algorithm, *inFlag, *outFlag)
	}

	if *decompressFlag {
		decompress(*algorithm, *inFlag, *outFlag)
	}
}

func compress(algorithm string, inputFilename string, outputFilename string) error {
	bytes, err := fileio.ReadFile(inputFilename)
	if err != nil {
		return fmt.Errorf("Could not read input file: %s", err)
	}

	var compressed *vector.Vector

	switch algorithm {
	case "huffman":
		compressed = huffman.Compress(bytes)

		err = fileio.WriteFile(compressed, outputFilename)
		if err != nil {
			return fmt.Errorf("Could not write compressed file: %s", err)
		}

	case "lzw":
		compressed, err = lzw.Compress(bytes)
		if err != nil {
			return fmt.Errorf("Could not compress: %s", err)
		}

		err = fileio.WriteLZWFile(compressed, outputFilename)
		if err != nil {
			return fmt.Errorf("Could not write compressed file: %s", err)
		}

	case "default":
		return fmt.Errorf("Invalid compression algorithm provided: %s", algorithm)
	}

	return nil
}

func decompress(algorithm string, inputFilename string, outputFilename string) error {
	switch algorithm {
	case "huffman":
		bytes, err := fileio.ReadFile(inputFilename)
		if err != nil {
			return fmt.Errorf("Could not read input file: %s", err)
		}

		decompressed, err := huffman.Decompress(bytes)
		if err != nil {
			return fmt.Errorf("Data decompression failed: %s", err)
		}

		err = fileio.WriteFile(decompressed, outputFilename)
		if err != nil {
			return fmt.Errorf("Could not write decompressed file: %s", err)
		}

	case "lzw":
		codes, err := fileio.ReadLZWFile(inputFilename)
		if err != nil {
			return fmt.Errorf("Could not read compressed file: %s", err)
		}

		decompressed, err := lzw.Decompress(codes)
		if err != nil {
			return fmt.Errorf("Data decompression failed: %s", err)
		}

		err = fileio.WriteFile(decompressed, outputFilename)
		if err != nil {
			return fmt.Errorf("Could not write decompressed file: %s", err)
		}

	case "default":
		return fmt.Errorf("Invalid compression algorithm provided: %s", algorithm)
	}

	return nil
}
