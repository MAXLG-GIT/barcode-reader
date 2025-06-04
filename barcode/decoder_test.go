package barcode

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDecodeTmpImages decodes all jpg images in the ../tmp directory and prints
// the detected barcode text. It logs decode errors instead of failing so that
// all barcodes are attempted.
func TestDecodeTmpImages(t *testing.T) {
	dir := filepath.Join("..", "tmp")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read tmp directory: %v", err)
	}

	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".jpg" {
			continue
		}
		name := e.Name()
		if name == "barimage6.jpg" || name == "barimage7.jpg" || name == "barimage10.jpg" {
			t.Logf("skipping %s due to size", name)
			continue
		}
		path := filepath.Join(dir, name)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			text, err := DecodeCode128FromFile(path)
			if err != nil {
				t.Logf("decode error: %v", err)
				return
			}
			t.Logf("%s", text)
		})
	}
}
