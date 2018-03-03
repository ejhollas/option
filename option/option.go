package option

import (
	"container/list"
	"fmt"
)

type OptionCB func(string) bool

type Option struct {
	text        string
	description string
	callback    OptionCB
	data        string
}

func (o Option) String() string {
	return fmt.Sprintf("-%-10s %s", o.text, o.description)
}

func NewOption(text, description string) *Option {
	return &Option{text, description, nil, ""}
}

func NewOptionCB(text, description string, callback OptionCB) *Option {
	return &Option{text, description, callback, ""}
}

// Verb contains an action callback, text description and optional options
type Verb struct {
	main       *Option
	suboptions *list.List
	callback   OptionCB
}

// NewVerb creates a new verb with an empty list of options
func NewVerb(text, description string, callback OptionCB) *Verb {
	v := Verb{}
	v.main = NewOption(text, description)
	v.callback = callback
	v.suboptions = list.New()
	return &v
}

// OnVerbFound calls the callback attached to the verb if it exits
func (v Verb) OnVerbFound(val string) {
	if v.callback != nil {
		v.callback(val)
	}
}

func (v Verb) String() string {
	s := fmt.Sprintf("%-5s %s", v.main.text, v.main.description)
	for e := v.suboptions.Front(); e != nil; e = e.Next() {
		s = s + fmt.Sprintf("\n     -%s", e.Value)
	}
	return s
}

// AddOption adds an option to the verb
func (v Verb) AddOption(o *Option) {
	v.suboptions.PushBack(o)
}
