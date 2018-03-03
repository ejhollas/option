package main

import (
	"container/list"
	"fmt"
)

type optionCB func(string) bool

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

type verb struct {
	main       *option
	suboptions *list.List
	callback   optionCB
}

func newVerb(text, description string, callback optionCB) *verb {
	v := verb{}
	v.main = newOption(text, description)
	v.callback = callback
	v.suboptions = list.New()
	return &v
}

func (v verb) onVerbFound(val string) {
	if v.callback != nil {
		v.callback(val)
	}
}

func (v verb) String() string {
	s := fmt.Sprintf("%-5s %s", v.main.text, v.main.description)
	for e := v.suboptions.Front(); e != nil; e = e.Next() {
		s = s + fmt.Sprintf("\n     -%s", e.Value)
	}
	return s
}

func onVerbGet(val string) bool {
	fmt.Println("Get verb here saw " + val)
	return true
}

func main() {
	opt := newOption("test", "Select the test option")
	verb := newVerb("get", "Retrieve information about the node", onVerbGet)
	verb.suboptions.PushBack(opt)
	fmt.Println(verb)
	verb.onVerbFound("Bitching")
}
