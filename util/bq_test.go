package util

import (
	"testing"
)

func TestBQNullString(t *testing.T) {
	b1 := BQNullString("s")
	if got, want := b1.StringVal, "s"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := b1.Valid, true; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	b2 := BQNullString(nil)
	if got, want := b2.Valid, false; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
