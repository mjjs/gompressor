package fileio

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mjjs/gompressor/vector"
)

func WriteLZWFile(codes *vector.Vector, filename string) error {
	absolutePath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	file, err := os.Create(absolutePath)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	for i := 0; i < codes.Size(); i++ {
		code := codes.MustGet(i)

		codeBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(codeBytes, code.(uint16))

		_, err := writer.Write(codeBytes)

		if err != nil {
			return err
		}
	}

	err = file.Sync()
	if err != nil {
		panic(err)
	}

	return writer.Flush()
}

func ReadLZWFile(filename string) (*vector.Vector, error) {
	absolutePath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absolutePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	codes := vector.New()

	for {
		buf := make([]byte, 2)
		_, err := io.ReadAtLeast(file, buf, 2)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		code := binary.BigEndian.Uint16(buf)

		codes.Append(code)
	}

	return codes, nil
}

func ReadFile(filename string) (*vector.Vector, error) {
	absolutePath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absolutePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	byteVector := vector.New(0, uint(len(bytes)))
	for _, b := range bytes {
		byteVector.Append(b)
	}

	return byteVector, nil
}

func WriteFile(byteVector *vector.Vector, filename string) error {
	absolutePath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	file, err := os.Create(absolutePath)
	if err != nil {
		return err
	}

	defer file.Close()

	bytes := make([]byte, 0, byteVector.Size())
	for i := 0; i < byteVector.Size(); i++ {
		bytes = append(bytes, byteVector.MustGet(i).(byte))
	}

	_, err = file.Write(bytes)
	if err != nil {
		return nil
	}

	return file.Sync()
}
