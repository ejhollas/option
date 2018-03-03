package main

import (
	"fmt"
)

type optionCB func(option) bool

type option struct {
	text        string
	description string
	callback    optionCB
	data        string
}

func (o option) String() string {
	return fmt.Sprintf("-%-10s %s", o.text, o.description)
}

func newOption(text, description string) *option {
	return &option{text, description, nil, ""}
}

func newOptionCB(text, description string, callback optionCB) *option {
	return &option{text, description, callback, ""}
}

func main() {
	opt := newOption("test", "Verb to run a test")
	fmt.Println(opt)
}
