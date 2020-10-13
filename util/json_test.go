package util

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExportJSON(t *testing.T) {
	f := filepath.Join(os.TempDir(), "TestExportJSON.json")
	err := ExportJSON(time.Now(), f)
	if err != nil {
		t.Error(err)
	}
}
