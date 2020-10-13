package util

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestExportJSON(t *testing.T) {
	f, _ := ioutil.TempFile(os.TempDir(), "TestExportJSON")
	t.Log(f.Name())
	err := ExportJSON(time.Now(), f.Name())
	if err != nil {
		t.Error(err)
	}
}
