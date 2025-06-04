module barcode-reader

go 1.23.8

require github.com/disintegration/imaging v1.6.2

require (
	github.com/makiuchi-d/gozxing v0.1.1
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace github.com/makiuchi-d/gozxing => ./third_party/github.com/makiuchi-d/gozxing

replace golang.org/x/xerrors => ./third_party/golang.org/x/xerrors

replace golang.org/x/text => ./third_party/golang.org/x/text

replace golang.org/x/tools => ./third_party/golang.org/x/tools

replace golang.org/x/mod => ./third_party/golang.org/x/mod

replace golang.org/x/sync => ./third_party/golang.org/x/sync

replace golang.org/x/net => ./third_party/golang.org/x/net

replace golang.org/x/telemetry => ./third_party/golang.org/x/telemetry

replace golang.org/x/sys => ./third_party/golang.org/x/sys

replace golang.org/x/crypto => ./third_party/golang.org/x/crypto

replace golang.org/x/term => ./third_party/golang.org/x/term
