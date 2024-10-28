package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	extension string
	size      int64
	list      bool
	delete    bool
	writer    io.Writer
	archive   string
}

func main() {
	root := flag.String("root", ".", "Root directory to start")
	extension := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	list := flag.Bool("list", false, "List files only")
	delete := flag.Bool("delete", false, "Delete files")
	filename := flag.String("log", "", "Log deletes to this file")
	archive := flag.String("archive", "", "Archive files")
	flag.Parse()

	var (
		writer = os.Stdout
		err    error
	)

	if *filename != "" {
		writer, err = os.OpenFile(*filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer writer.Close()
	}

	config := config{
		extension: *extension,
		size:      *size,
		list:      *list,
		delete:    *delete,
		writer:    writer,
		archive:   *archive,
	}

	if err := run(*root, os.Stdout, config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, config config) error {
	logger := log.New(config.writer, "DELETE FILE: ", log.LstdFlags)

	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filterOut(path, config.extension, config.size, info) {
			return nil
		}

		if config.list {
			return listFile(path, out)
		}

		if config.archive != "" {
			if err := archiveFile(config.archive, root, path); err != nil {
				return err
			}
		}

		if config.delete {
			return deleteFile(path, logger)
		}

		return listFile(path, out)
	})
}
