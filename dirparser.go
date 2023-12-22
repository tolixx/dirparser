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

// Processor just recursive walk on file tree
type Processor interface {
	Process(filename string) error
}

type proxyProcessor struct {
	parser Parser
}

func (p *proxyProcessor) Process(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader, err := p.parser.Init(file, filename)
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
		p.parser.Parse(record)
	}

	return p.parser.Close()
}

func parseProcessor(parser Parser) Processor {
	return &proxyProcessor{parser: parser}
}

// ParsePath parse using Parse Interface
func ParsePath(path string, p Parser) error {
	return ProcessPath(path, parseProcessor(p))
}

// ProcessPath for parse path
func ProcessPath(path string, p Processor) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return processDir(path, p)
	}
	return ProcessFile(path, p)
}

// ProcessFile process a single file
func ProcessFile(path string, p Processor) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer reader.Close()
	return p.Process(path)
}

func processDir(path string, p Processor) error {
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
		ProcessPath(filename, p)
	}
	return nil
}

func readDir(name string) ([]os.DirEntry, error) {
	d, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	return d.ReadDir(0)
}
