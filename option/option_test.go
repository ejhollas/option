package option

import (
	"testing"
)

// See https://golang.org/pkg/testing/

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

func ExampleParser() {
	args := make([]string, 2)
	args[0] = "test"
	args[1] = ""
	p := NewParser()
	if p.Parse(args) {
		p.Run()
	}
	// Output:
	// usage: test  [-options] [command] [--command_option=value]
}
