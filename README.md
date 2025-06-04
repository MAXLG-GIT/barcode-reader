# Barcode Reader

This repository contains a small Go module for decoding CODE128 barcodes from images.
It relies on [gocv](https://gocv.io) for basic image operations and uses the
[ZBar](https://github.com/ZBar/ZBar) utility `zbarimg` to decode the barcode data.

## Usage

```
go run ./cmd/barcode-reader -img path/to/photo.jpg
```

The decoder will attempt several preprocessing steps and handle images rotated at
0°, 90°, 180° and 270°.


## Development

Run `go mod tidy` to download dependencies. Ensure OpenCV is installed for gocv
and the `zbar-tools` package is available so that the `zbarimg` command can be
invoked.

