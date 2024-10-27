package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type config struct {
	extension string
	size      int64
	list      bool
}

func main() {
	root := flag.String("root", ".", "Root directory to start")
	extension := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	list := flag.Bool("list", false, "List files only")
	flag.Parse()

	config := config{
		extension: *extension,
		size:      *size,
		list:      *list,
	}

	if err := run(*root, os.Stdout, config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, config config) error {
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

		return listFile(path, out)
	})
}
