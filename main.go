package main

import (
	"flag"
	"log"

	"github.com/mjjs/gompressor/fileio"
	"github.com/mjjs/gompressor/lzw"
	"github.com/mjjs/gompressor/vector"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	encodeFlag := flag.Bool("encode", false, "Encode the given file")
	decodeFlag := flag.Bool("decode", false, "Decode the given file")
	inFlag := flag.String("in", "", "input file")
	outFlag := flag.String("out", "", "output file")

	flag.Parse()

	if !*encodeFlag && !*decodeFlag {
		log.Fatal("No -encode or -decode flag given")
	} else if *encodeFlag && *decodeFlag {
		log.Fatal("Cannot use -encode and -decode flag at the same time")
	}

	if len(*inFlag) == 0 || len(*outFlag) == 0 {
		log.Fatal("Input or output file flag missing")
	}

	if *encodeFlag {
		bytes, err := fileio.ReadFile(*inFlag)
		if err != nil {
			log.Fatalf("Could not read input file: %s", err)
		}

		bv := vector.New(uint(len(bytes)))
		for i, b := range bytes {
			bv.MustSet(i, b)
		}

		encoded, err := lzw.Compress(bv)
		if err != nil {
			log.Fatalf("Data compression failed: %s", err)
		}

		err = fileio.WriteLZWFile(encoded, *outFlag)
		if err != nil {
			log.Fatalf("Could not write compressed file: %s", err)
		}
	}

	if *decodeFlag {
		codes, err := fileio.ReadLZWFile(*inFlag)
		if err != nil {
			log.Fatalf("Could not read compressed file: %s", err)
		}

		bytes, err := lzw.Decompress(codes)
		if err != nil {
			log.Fatalf("Data decompression failed: %s", err)
		}

		byts := []byte{}
		for i := 0; i < bytes.Size(); i++ {
			byts = append(byts, bytes.MustGet(i).(byte))
		}

		err = fileio.WriteFile(byts, *outFlag)
		if err != nil {
			log.Fatalf("Could not write decompressed file: %s", err)
		}
	}
}
