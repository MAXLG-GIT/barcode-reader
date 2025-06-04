package barcode

import (
	"bytes"
	"fmt"
	"image"
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
	}
	return "", fmt.Errorf("code128 not found")
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
