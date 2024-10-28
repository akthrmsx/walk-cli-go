package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
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

func deleteFile(path string, logger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	logger.Panicln(path)
	return nil
}

func archiveFile(destDir string, root string, path string) error {
	info, err := os.Stat(destDir)

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not directory", destDir)
	}

	relDir, err := filepath.Rel(root, filepath.Dir(path))

	if err != nil {
		return err
	}

	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	target := filepath.Join(destDir, relDir, dest)

	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(path)

	if err != nil {
		return err
	}

	defer in.Close()

	writer := gzip.NewWriter(out)
	writer.Name = filepath.Base(path)

	if _, err = io.Copy(writer, in); err != nil {
		return err
	}

	if err = writer.Close(); err != nil {
		return err
	}

	return out.Close()
}
