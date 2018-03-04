package option

import (
	"container/list"
	"fmt"
	"strings"
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

// NewVerb returns an initilized Verb
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

// Parser structure
type Parser struct {
	programName          string
	verbs                *list.List
	options              *list.List
	activeVerb           *Verb
	activePreVerbOptions *list.List
	activeOptions        *list.List
}

// NewParser returns an initialized Parser
func NewParser(programName string) *Parser {
	p := Parser{}
	p.programName = programName
	p.verbs = list.New()
	p.options = list.New()
	return &p
}

// Parse reviews command line array and creates list of verb and options to run
func (p *Parser) Parse(args []string) bool {
	argHandled := 0
	p.activeVerb = nil

	for _, arg := range args {
		if nil == p.activeVerb && arg[1] == '-' {
			// Find options before verbs
			argOption := strings.Split(arg, "=")
			for e := p.options.Front(); e != nil; e = e.Next() {
				option := e.Value.(*Option)
				if argOption[0] == option.text {
					if len(argOption) > 1 {
						option.data = argOption[1]
						p.activePreVerbOptions.PushBack(option)
						argHandled++
					}
				}
			}
		} else {
			if nil == p.activeVerb {
				// Find Verb
				for e := p.verbs.Front(); e != nil; e = e.Next() {
					verb := e.Value.(*Verb)
					if verb.main.text == arg {
						p.activeVerb = verb
						argHandled++
					}
				}
			} else {
				// Assume all other options are verb options
				if len(args) < 2 {
					continue
				}
				// When we split, ignore the first character
				argOption := strings.Split(arg[2:], "=")
				for e := p.activeVerb.suboptions.Front(); e != nil; e = e.Next() {
					option := e.Value.(*Option)
					if argOption[0] == option.text {
						if len(argOption) > 1 {
							option.data = argOption[1]
							p.activeOptions.PushBack(option)
							argHandled++
						}
					}
				}
			}
		}
	}
	if len(args) == 0 {
		p.ShowHelp()
	} else {
		if len(args) != argHandled {
			p.ShowHelp()
		} else {
			return true
		}
	}
	if p.activeVerb == nil {
		// Missing Verb
		return false
	} else {
		// Invalid Option
		return false
	}
}

// ShowHelp prints the options, verbs, and suboptions dynamically
func (p Parser) ShowHelp() {
	fmt.Printf("usage: %s  [-options] [command] [--command_option=value]\n", p.programName)
	fmt.Println("options:")
	for e := p.options.Front(); e != nil; e = e.Next() {
		fmt.Printf("   %s\n", e.Value.(*Option))
	}
	fmt.Println("commands:")
	for e := p.verbs.Front(); e != nil; e = e.Next() {
		fmt.Printf("   %s\n", e.Value.(*Verb))
	}
}

// AddOption adds the option to the parser
func (p Parser) AddOption(o *Option) {
	p.options.PushBack(o)
}

// AddVerb adds the verb to the parser
func (p Parser) AddVerb(v *Verb) {
	p.verbs.PushBack(v)
}
