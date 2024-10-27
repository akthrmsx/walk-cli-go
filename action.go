package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func filterOut(path string, extension string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if extension != "" && filepath.Ext(path) != extension {
		return true
	}

	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}
