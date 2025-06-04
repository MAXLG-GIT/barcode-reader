package main

import (
	"flag"
	"fmt"
	"log"

	"barcode-reader/barcode"
)

func main() {
	path := flag.String("img", "", "path to image file")
	flag.Parse()
	if *path == "" {
		log.Fatal("no image provided")
	}
	text, err := barcode.DecodeCode128FromFile(*path)
	if err != nil {
		log.Fatalf("failed to decode: %v", err)
	}
	fmt.Println(text)
}
