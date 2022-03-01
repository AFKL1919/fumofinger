package utils

import (
	"bufio"
	"os"
)

type File struct {
	Fp   *os.File
	Scan *bufio.Scanner
}

func Open(filename string) (*File, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	file := &File{
		Fp:   fp,
		Scan: bufio.NewScanner(fp),
	}
	return file, nil
}

func (file File) ReadLine() (string, bool) {
	isNotEOF := file.Scan.Scan()
	if isNotEOF {
		return file.Scan.Text(), true
	}
	return "", false
}

func (file File) Close() error {
	if err := file.Scan.Err(); err != nil {
		return err
	}
	return file.Fp.Close()
}
