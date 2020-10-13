package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func ExportJSON(data interface{}, filename string) error {
	log.Println("exporting JSON to", filename)
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "\t")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buf.Bytes(), os.ModeAppend)
}
