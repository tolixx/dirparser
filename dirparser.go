package dirparser

import (
	"io"
	"os"
	"strings"
)

// Reader interface compilant with csv.Reader to use it with csv
type Reader interface {
	Read() ([]string, error)
}

// Parser for implement real parser
type Parser interface {
	Init(io.Reader, string) (Reader, error)
	Parse([]string) error
	Close() error
}

// ParsePath for parse path
func ParsePath(path string, parser Parser) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return parseDir(path, parser)
	}
	return ParseFile(path, parser)
}

// ParseFile parse single file
func ParseFile(path string, parser Parser) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer reader.Close()
	return parseReader(reader, parser, path)
}

func parseDir(path string, parser Parser) error {
	path = strings.TrimRight(path, "/")
	files, err := readDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		name := file.Name()
		if name[0] == '.' {
			continue
		}

		filename := path + "/" + name
		ParsePath(filename, parser)
	}
	return nil
}

func parseReader(r io.Reader, parser Parser, filename string) error {
	reader, err := parser.Init(r, filename)
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		parser.Parse(record)
	}

	return parser.Close()
}

func readDir(name string) ([]os.DirEntry, error) {
	d, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	return d.ReadDir(0)
}
