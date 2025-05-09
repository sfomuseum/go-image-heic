package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sfomuseum/go-image-heic"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Convert one or more HEIC files to JPEG files, preserving EXIF data.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s path(N) path(N)\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	for _, path := range flag.Args() {

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatalf("Failed to derive absolute path for '%s', %v", path, err)
		}

		r, err := os.Open(abs_path)

		if err != nil {
			log.Fatal(err)
		}

		defer r.Close()

		root := filepath.Dir(abs_path)
		fname := filepath.Base(abs_path)
		ext := filepath.Ext(fname)

		jpeg_fname := strings.Replace(fname, ext, ".jpg", 1)
		jpeg_path := filepath.Join(root, jpeg_fname)

		wr, err := os.OpenFile(jpeg_path, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Fatalf("Failed to open '%s' for writing, %v", jpeg_path, err)
		}

		err = heic.ToJPEG(r, wr)

		if err != nil {
			log.Fatalf("Failed to convert '%s' to JPEG, %v", abs_path, err)
		}

		err = wr.Close()

		if err != nil {
			log.Fatalf("Failed to close '%s' after writing, %v", jpeg_path, err)
		}
	}
}
