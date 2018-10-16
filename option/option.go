package option

import (
	"container/list"
	"fmt"
	"path/filepath"
	"strings"
)

type OptionCB func(option *Option) (result bool, err error)
type VerbCB func(verb *Verb) (result bool, err error)

var debug = false

type Optioner interface {
	onOption(option *Option) (result bool, err error)
}

type Verber interface {
	OnVerb() (result bool, err error)
}

type VerberFunc func() (bool, error)

type Option struct {
	text        string
	description string
	callback    OptionCB
	Data        string
}

func (o *Option) String() string {
	if o.Data != "" {
		return fmt.Sprintf("-%-10s %s Data='%s'", o.text, o.description, o.Data)
	}
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
func (o *Option) OnOptionFound() (result bool, err error) {
	if o.callback != nil {
		return o.callback(o)
	}
	return true, nil
}

// Verb contains an action callback, text description and optional options
type Verb struct {
	main             *Option
	suboptions       *list.List
	callback         VerbCB
	opts             map[string]*Option
	isOptionRequired bool
}

// NewVerb returns an initilized Verb
func NewVerb(text, description string, callback VerbCB) *Verb {
	v := Verb{}
	v.main = NewOption(text, description)
	v.callback = callback
	v.suboptions = list.New()
	v.opts = make(map[string]*Option)
	v.isOptionRequired = false
	return &v
}

// OnVerbFound calls the callback attached to the verb if it exits
func (v *Verb) OnVerbFound() {
	if v.callback != nil {
		v.callback(v)
	}
}

func (v *Verb) String() string {
	s := fmt.Sprintf("%-5s %s", v.main.text, v.main.description)
	for e := v.suboptions.Front(); e != nil; e = e.Next() {
		s = s + fmt.Sprintf("\n     -%s", e.Value)
	}
	return s
}

// AddOption adds an option to the verb
func (v *Verb) AddOption(o *Option) {
	v.suboptions.PushBack(o)
	v.opts[o.text] = o
}

// SetRequiresOption marka this verb as requiring an option
func (v *Verb) SetRequiresOption() {
	v.isOptionRequired = true
}

// IsOptionRequired returns true if the verbs requires an option
func (v *Verb) IsOptionRequired() bool {
	return v.isOptionRequired
}

// GetOption returns a pointer to the option by name
func (v *Verb) GetOption(name string) (p *Option) {
	return v.opts[name]
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
		if debug {
			fmt.Println("Debug: Considering: " + arg)
		}
		if nil == p.activeVerb && len(arg) > 1 && arg[0] == '-' {
			// Find options before verbs
			// When we split, ignore the first character
			argOption := strings.Split(arg[1:], "=")
			for e := p.options.Front(); e != nil; e = e.Next() {
				option := e.Value.(*Option)
				if argOption[0] == option.text {
					if len(argOption) > 1 {
						option.Data = argOption[1]
						p.activePreVerbOptions.PushBack(option)
						argHandled++
						if debug {
							fmt.Println("Debug: Found PreOption: " + option.String())
						}
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
				// When we split, ignore the first two characters
				argOption := strings.Split(arg[2:], "=")
				for e := p.activeVerb.suboptions.Front(); e != nil; e = e.Next() {
					option := e.Value.(*Option)
					if argOption[0] == option.text {
						if len(argOption) > 1 {
							option.Data = argOption[1]
							p.activeOptions.PushBack(option)
							argHandled++
						} else {
							p.activeOptions.PushBack(option)
							argHandled++
						}
					}
				}
			}
		}
	}

	// Is not all of the args were handled, show help
	if argsToConsider != argHandled {
		p.ShowHelp()
	} else {
		return true
	}

	if debug {
		fmt.Printf("Debug: args=%d handled=%d\n", argsToConsider, argHandled)
	}
	if p.activeVerb == nil {
		// Missing Verb
		if debug {
			fmt.Println("Debug: Missing Verb")
		}
	} else {
		// Invalid Option
		if debug {
			fmt.Println("Debug: Invalid Option")
		}
	}
	return false
}

// ShowHelp prints the options, verbs, and suboptions dynamically
func (p Parser) ShowHelp() {
	fmt.Printf("usage: %s  [-options] [command] [--command_option=value]\n", p.programName)
	if p.options.Len() > 0 {
		fmt.Println("options:")
		for e := p.options.Front(); e != nil; e = e.Next() {
			fmt.Printf("   %s\n", e.Value.(*Option))
		}
	}
	if p.verbs.Len() > 0 {
		fmt.Println("commands:")
		for e := p.verbs.Front(); e != nil; e = e.Next() {
			fmt.Printf("   %s\n", e.Value.(*Verb))
		}
	}
}

// GetActiveVerb returns the active verb or nil
func (p Parser) GetActiveVerb() (v *Verb) {
	return p.activeVerb
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
func (p Parser) Run() (bool, error) {
	if nil != p.activePreVerbOptions {
		// Process the pre-verb options
		for e := p.activePreVerbOptions.Front(); e != nil; e = e.Next() {
			result, err := e.Value.(*Option).OnOptionFound()
			if err != nil {
				return result, err
			}
		}
	}
	// Process the verb
	if nil != p.activeVerb {
		// Process the post-verb options
		for e := p.activeOptions.Front(); e != nil; e = e.Next() {
			e.Value.(*Option).OnOptionFound()
		}
		// Abort the run, if the verb requires an option
		if p.activeVerb.IsOptionRequired() && p.activeOptions.Len() == 0 {
			fmt.Printf("%s\n", p.activeVerb)
			return false, nil
		}
		// Finally process the verb
		p.activeVerb.OnVerbFound()
	}
	return true, nil
}
