package barcode

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"
	"os/exec"
	"strings"

	"github.com/disintegration/imaging"
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

		up := imaging.Resize(rimg, rimg.Bounds().Dx()*2, 0, imaging.Lanczos)
		text, err = decodeWithZXing(up)
		if err == nil && text != "" {
			return text, nil
		}

		timg := threshold(rimg, 0.5)
		text, err = decodeWithZbar(timg)
		if err == nil && text != "" {
			return text, nil
		}
	}
	return "", fmt.Errorf("code128 not found")
}

// preprocess scales large images down and converts them to grayscale
// to improve barcode recognition speed and accuracy.
func preprocess(img image.Image) *image.NRGBA {
	b := img.Bounds()
	const target = 1000
	ratio := float64(target) / float64(b.Dx())
	img = imaging.Resize(img, target, int(float64(b.Dy())*ratio), imaging.Lanczos)

	gray := imaging.Grayscale(img)
	return imaging.Clone(gray)
}

func threshold(img image.Image, t float64) *image.NRGBA {
	b := img.Bounds()
	dst := image.NewNRGBA(b)
	limit := uint8(t * 255)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			if c.Y > limit {
				dst.Set(x, y, color.White)
			} else {
				dst.Set(x, y, color.Black)
			}
		}
	}
	return dst
}

// decodeWithZbar saves the image to a temporary JPEG and invokes zbarimg to decode it.
func decodeWithZbar(img image.Image) (string, error) {
	tmp, err := os.CreateTemp("", "barcode-*.jpg")
	if err != nil {
		return "", err
	}
	fname := tmp.Name()
	tmp.Close()
	if err := imaging.Save(img, fname, imaging.JPEGQuality(95)); err != nil {
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

// decodeWithZXing saves the image to a temporary PNG and invokes ZXingReader
// to decode CODE128 barcodes.
func decodeWithZXing(img image.Image) (string, error) {
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

	cmd := exec.Command("ZXingReader", "-format", "CODE_128", fname)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("zxing: %v: %s", err, stderr.String())
	}

	for _, line := range strings.Split(out.String(), "\n") {
		if strings.HasPrefix(line, "Text:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Text:")), nil
		}
	}
	return "", fmt.Errorf("no result")
}
