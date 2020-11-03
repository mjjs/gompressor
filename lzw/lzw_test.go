package lzw

import (
	"reflect"
	"testing"

	"github.com/mjjs/gompressor/bytevector"
	"github.com/mjjs/gompressor/fileio"
)

var compressTestCases = []struct {
	name           string
	input          []byte
	expectedOutput []uint16
	shouldError    bool
}{
	{
		name:           "Hello world",
		input:          []byte("Hello World"),
		expectedOutput: []uint16{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100},
		shouldError:    false,
	},
	{
		name:           "Empty input",
		input:          []byte{},
		expectedOutput: []uint16{},
		shouldError:    false,
	},
}

func TestCompress(t *testing.T) {
	for _, testCase := range compressTestCases {
		bv := bytevector.New(0, uint(len(testCase.input)))
		bv.Append(testCase.input...)
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
	input          []uint16
	expectedOutput *bytevector.Bytevector
	shouldError    bool
}{
	{
		name:           "Hello world",
		input:          []uint16{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100},
		expectedOutput: bytevector.New().AppendToCopy([]byte("Hello World")...),
		shouldError:    false,
	},
	{
		name:           "Empty input",
		input:          []uint16{},
		expectedOutput: bytevector.New(),
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

			original := bytevector.New(uint(len(byts)))
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
