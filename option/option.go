package option

import (
	"container/list"
	"fmt"
	"path/filepath"
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

// NewOption returns an initialized Option
func NewOption(text, description string) *Option {
	return &Option{text, description, nil, ""}
}

// NewOptionCB sets the option's callback. Set nil to disable a previous setting.
func NewOptionCB(text, description string, callback OptionCB) *Option {
	return &Option{text, description, callback, ""}
}

// OnOptionFound calls the callback attached to the verb if it exits
func (o Option) OnOptionFound(val string) {
	if o.callback != nil {
		o.callback(val)
	}
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
func NewParser() *Parser {
	p := Parser{}
	p.programName = ""
	p.verbs = list.New()
	p.options = list.New()
	return &p
}

// Parse reviews command line array and creates list of verb and options to run
func (p *Parser) Parse(args []string) bool {
	argHandled := 0
	argsToConsider := len(args) - 1
	p.activeVerb = nil
	p.activePreVerbOptions = list.New()
	p.activeOptions = list.New()
	p.programName = filepath.Base(args[0])

	// Bail if no args presented
	if argsToConsider == 0 {
		p.ShowHelp()
		return false
	}

	// Begin parsing args after the program name
	for _, arg := range args[1:] {
		fmt.Println("Debug: Considering: " + arg)
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
						fmt.Println("Debug: Found PreOption: " + option.String())
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
						fmt.Println("Debug: Found Verb: " + verb.String())
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
							fmt.Println("Debug: Found Option: " + option.String())
						}
					}
				}
			}
		}
	}
	fmt.Printf("Debug: args=%d handled=%d\n", argsToConsider, argHandled)

	// Is not all of the args were handled, show help
	if argsToConsider != argHandled {
		p.ShowHelp()
	} else {
		return true
	}

	if p.activeVerb == nil {
		// Missing Verb
		fmt.Println("Debug: Missing Verb")
	} else {
		// Invalid Option
		fmt.Println("Debug: Invalid Option")
	}
	return false
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

// Run excutes parsed verbs and options
func (p Parser) Run() {
	if nil != p.activePreVerbOptions {
		// Process the pre-verb options
		for e := p.activeOptions.Front(); e != nil; e = e.Next() {
			e.Value.(*Option).OnOptionFound("")
		}
	}
}
