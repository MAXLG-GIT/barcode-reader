package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"barcode-reader/barcode"
)

func main() {
	dir := filepath.Join(".", "tmp")
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read tmp directory: %v", err)
	}

	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".jpg" {
			continue
		}
		path := filepath.Join(dir, e.Name())
		text, err := barcode.DecodeCode128FromFile(path)
		if err != nil {
			fmt.Printf("%s: error %v\n", e.Name(), err)
		} else {
			fmt.Printf("%s: %s\n", e.Name(), text)
		}
	}
}
