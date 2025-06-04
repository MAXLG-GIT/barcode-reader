package barcode

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gocv.io/x/gocv"
)

// DecodeCode128FromFile attempts to decode a CODE128 barcode from the image at the given path.
// It rotates the image in several orientations to improve detection.
func DecodeCode128FromFile(path string) (string, error) {
	img := gocv.IMRead(path, gocv.IMReadColor)
	if img.Empty() {
		return "", fmt.Errorf("cannot read image %s", path)
	}
	defer img.Close()

	rotations := []gocv.RotateFlag{
		gocv.RotateNone,
		gocv.Rotate90Clockwise,
		gocv.Rotate180Clockwise,
		gocv.Rotate90CounterClockwise,
	}

	for _, rot := range rotations {
		rimg := gocv.NewMat()
		gocv.Rotate(img, &rimg, rot)
		text, err := decodeWithZbar(rimg)
		rimg.Close()
		if err == nil && text != "" {
			return text, nil
		}
	}
	return "", fmt.Errorf("code128 not found")
}

// decodeWithZbar saves the Mat to a temporary PNG and invokes zbarimg to decode it.
func decodeWithZbar(m gocv.Mat) (string, error) {
	tmp, err := os.CreateTemp("", "barcode-*.png")
	if err != nil {
		return "", err
	}
	fname := tmp.Name()
	tmp.Close()
	if ok := gocv.IMWrite(fname, m); !ok {
		os.Remove(fname)
		return "", fmt.Errorf("failed to write temp image")
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
