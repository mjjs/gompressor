package lzw

import (
	"reflect"
	"testing"

	"github.com/mjjs/gompressor/fileio"
	"github.com/mjjs/gompressor/vector"
)

var compressTestCases = []struct {
	name           string
	input          []byte
	expectedOutput *vector.Vector
	shouldError    bool
}{
	{
		name:           "Hello world",
		input:          []byte("Hello World"),
		expectedOutput: vector.New().AppendToCopy(uint16(72), uint16(101), uint16(108), uint16(108), uint16(111), uint16(32), uint16(87), uint16(111), uint16(114), uint16(108), uint16(100)),
		shouldError:    false,
	},
	{
		name:           "Empty input",
		input:          []byte{},
		expectedOutput: vector.New(),
		shouldError:    false,
	},
}

func TestCompress(t *testing.T) {
	for _, testCase := range compressTestCases {
		bv := vector.New(0, uint(len(testCase.input)))
		for _, b := range testCase.input {
			bv.Append(b)
		}
		actual, err := Compress(bv)

		if err != nil && !testCase.shouldError {
			t.Errorf("%s: unexpected error %s", testCase.name, err)
		} else if err == nil && testCase.shouldError {
			t.Errorf("%s: expected error, got nil", testCase.name)
		}

		if !reflect.DeepEqual(actual, testCase.expectedOutput) {
			t.Errorf("%s\nexpected %+v\ngot %+v", testCase.name, testCase.expectedOutput, actual)
		}
	}
}

var decompressTestCases = []struct {
	name           string
	input          *vector.Vector
	expectedOutput *vector.Vector
	shouldError    bool
}{
	{
		name:  "Hello world",
		input: vector.New().AppendToCopy(uint16(72), uint16(101), uint16(108), uint16(108), uint16(111), uint16(32), uint16(87), uint16(111), uint16(114), uint16(108), uint16(100)),
		expectedOutput: vector.New().AppendToCopy(byte('H'), byte('e'), byte('l'), byte('l'), byte('o'),
			byte(' '), byte('W'), byte('o'), byte('r'), byte('l'), byte('d')),
		shouldError: false,
	},
	{
		name:           "Empty input",
		input:          vector.New(),
		expectedOutput: vector.New(),
		shouldError:    false,
	},
}

func TestDecompress(t *testing.T) {
	for _, testCase := range decompressTestCases {
		actual, err := Decompress(testCase.input)

		if err != nil && !testCase.shouldError {
			t.Errorf("%s: unexpected error %s", testCase.name, err)
		} else if err == nil && testCase.shouldError {
			t.Errorf("%s: expected error, got nil", testCase.name)
		}

		if !reflect.DeepEqual(actual, testCase.expectedOutput) {
			t.Errorf("%s\nexpected %+v\ngot %+v", testCase.name, testCase.expectedOutput, actual)
		}
	}
}

func TestDecompressedEqualsOriginal(t *testing.T) {
	for _, filename := range []string{"../testdata/E.coli", "../testdata/world192.txt"} {
		t.Run(filename, func(t *testing.T) {
			byts, err := fileio.ReadFile(filename)
			if err != nil {
				t.Errorf("Expected no error, got %s", err)
			}

			original := vector.New(uint(len(byts)))
			for i, byt := range byts {
				original.MustSet(i, byt)
			}

			compressed, err := Compress(original)
			if err != nil {
				t.Errorf("Expected no error, got %s", err)
			}

			decompressed, err := Decompress(compressed)
			if err != nil {
				t.Errorf("Expected no error, got %s", err)
			}

			if reflect.DeepEqual(original, decompressed) {
				t.Error("NOT OK")
			}
		})
	}
}
