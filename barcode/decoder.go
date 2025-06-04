package barcode

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"os/exec"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

// DecodeCode128FromFile attempts to decode a CODE128 barcode from the image at the given path.
// It rotates the image in several orientations to improve detection.
func DecodeCode128FromFile(path string) (string, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return "", fmt.Errorf("cannot read image %s: %w", path, err)
	}

	img = preprocess(img)

	rotations := []func(image.Image) *image.NRGBA{
		imaging.Clone,
		imaging.Rotate90,
		imaging.Rotate180,
		imaging.Rotate270,
	}

	for _, rot := range rotations {
		rimg := rot(img)
		text, err := decodeWithZbar(rimg)
		if err == nil && text != "" {
			return text, nil
		}
		text, err = decodeWithZXing(rimg)
		if err == nil && text != "" {
			return text, nil
		}
	}
	return "", fmt.Errorf("code128 not found")
}

// preprocess scales large images down and converts them to grayscale
// to improve barcode recognition speed and accuracy.
func preprocess(img image.Image) *image.NRGBA {
	// Resize if the image is very large
	const maxDim = 1600
	b := img.Bounds()
	if b.Dx() > maxDim {
		ratio := float64(maxDim) / float64(b.Dx())
		img = imaging.Resize(img, maxDim, int(float64(b.Dy())*ratio), imaging.Lanczos)
	}

	// Convert to grayscale and slightly increase contrast
	gray := imaging.Grayscale(img)
	gray = imaging.AdjustContrast(gray, 20)

	return imaging.Clone(gray)
}

// decodeWithZbar saves the Mat to a temporary PNG and invokes zbarimg to decode it.
func decodeWithZbar(img image.Image) (string, error) {
	tmp, err := os.CreateTemp("", "barcode-*.png")
	if err != nil {
		return "", err
	}
	fname := tmp.Name()
	tmp.Close()
	if err := imaging.Save(img, fname); err != nil {
		os.Remove(fname)
		return "", fmt.Errorf("failed to write temp image: %w", err)
	}
	defer os.Remove(fname)

	cmd := exec.Command("zbarimg", "--raw", fname)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("zbarimg: %v: %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// decodeWithZXing uses the pure Go gozxing library to decode CODE128 barcodes.
func decodeWithZXing(img image.Image) (string, error) {
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}
	reader := oned.NewCode128Reader()
	res, err := reader.Decode(bmp, map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	})
	if err != nil {
		return "", err
	}
	return res.GetText(), nil
}
