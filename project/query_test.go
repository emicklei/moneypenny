package project

import (
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	q := QueryPastDays("Test", time.Now(), time.Now().Add(-30*time.Hour*24))
	t.Log(q)
}
