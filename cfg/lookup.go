package cfg

import (
	"fmt"
	"os"
	"path/filepath"
)

func NewLookup(name string, depth int) *Lookup {
	return &Lookup{
		depth: depth,
		name:  name,
	}
}

type Lookup struct {
	depth int    // Number of directories up where the file is needed to be looked for
	name  string // Name of the file to search for
}

type pathParameters struct {
	currentDir string
	depth      int
	name       string
}

func (l *Lookup) FindFile() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get current working directory: %w", err)
	}
	for i := 0; i <= l.depth; i++ {
		path := generatePath(pathParameters{currentDir, i, l.name})
		err = l.checkExists(path)
		if err != nil {
			continue
		}

		return path, nil
	}

	return "", fmt.Errorf("unable to find file: %s", l.name)
}

func generatePath(params pathParameters) string {
	var path string
	for i := 0; i < params.depth; i++ {
		path = filepath.Join("..", path)
	}

	return filepath.Join(params.currentDir, path, params.name)
}

func (l *Lookup) checkExists(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file error: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()
	fileInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("retrive file info error: %w", err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("file %s is a directory", path)
	}

	return nil
}
