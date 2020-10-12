package util

import (
	"log"
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
