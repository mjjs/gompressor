package lzw

import (
	"reflect"
	"testing"
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
		actual, err := Compress(testCase.input)

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
	expectedOutput []byte
	shouldError    bool
}{
	{
		name:           "Hello world",
		input:          []uint16{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100},
		expectedOutput: []byte("Hello World"),
		shouldError:    false,
	},
	{
		name:           "Empty input",
		input:          []uint16{},
		expectedOutput: []byte{},
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
