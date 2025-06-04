# Barcode Reader

This repository contains a small Go module for decoding CODE128 barcodes from images.
It uses the [imaging](https://github.com/disintegration/imaging) library for image
rotation and the [ZBar](https://github.com/ZBar/ZBar) utility `zbarimg` to decode
the barcode data.

## Usage

```
go run ./cmd/barcode-reader
```

The program reads all `.jpg` files from the `tmp` directory in the project root
and prints each filename with the decoded barcode text. The decoder attempts
several rotations (0°, 90°, 180°, 270°) to improve recognition.

## Development

Run `go mod tidy` to download dependencies. Install the `zbar-tools` package so
that the `zbarimg` command can be invoked.

