package option

import (
	"testing"
)

func TestOption(t *testing.T) {
	opt := NewOption("Test", "Description")
	if opt.callback != nil {
		t.Error("Callback set unexpectantly")
	}
	if opt.Data != "" {
		t.Error("Data has a non-empty value")
	}
	want := ("-Test       Description")
	if opt.String() != want {
		t.Errorf("Got:'%s' Wanted:'%s'", opt.String(), want)
	}
}
