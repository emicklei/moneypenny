package util

import (
	"log"
	"os"
	"strings"
)

func CheckOpex(opex string) {
	if len(opex) == 0 {
		log.Fatalln("-opex cannot be empty")
	}
}

func CheckMonth(index int) {
	if index < 1 || index > 12 {
		log.Fatalln("-month must be in [1..12], got:", index)
	}
}

func CheckBigQueryTable(id string) {
	if p := strings.Split(id, "."); len(p) != 3 {
		log.Fatalln("full qualified bigquery table must be PROJECT.DATASET.TABLE, got:", id)
	}
}

func CheckNonEmpty(parameter, value string) {
	if len(value) == 0 {
		log.Fatalf("parameter [%s] cannot be empty", parameter)
	}
}

var checkedGCP = false

func CheckGCPCredentials() {
	if checkedGCP {
		return
	}
	fileRef := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(fileRef) == 0 {
		log.Fatalln("no value for environment variable GOOGLE_APPLICATION_CREDENTIALS")
	}
	_, err := os.Stat(fileRef)
	if os.IsNotExist(err) {
		log.Fatalln("file referenced by GOOGLE_APPLICATION_CREDENTIALS does not exist", err)
	}
	checkedGCP = true
}
