package dirparser

import (
	"bufio"
	"io"
	"strings"
)

type deepReader struct {
	scanner   *bufio.Scanner
	Separator string
}

// NewDeepReader returns multibyteCSV-like reader
func NewDeepReader(r io.Reader, sep string) Reader {
	dr := &deepReader{}
	dr.scanner = bufio.NewScanner(r)
	dr.Separator = sep
	return dr
}

func (dr *deepReader) Read() ([]string, error) {
	if !dr.scanner.Scan() {
		return nil, io.EOF
	}

	if err := dr.scanner.Err(); err != nil {
		return nil, err
	}

	parts := strings.Split(strings.TrimSpace(dr.scanner.Text()), dr.Separator)
	return parts, nil
}
