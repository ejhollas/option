package main

import (
	"eric/option/option"
	"fmt"
	"os"
)

type Program struct {
}

func onOptionGet(o *option.Option) (result bool, err error) {
	fmt.Printf("Get option here arg='%s'\n", o.Data)
	return false, nil
}

func onVerbGetAlone(v *option.Verb) (result bool, err error) {
	fmt.Printf("Get verbAlone here arg='%s'\n", v.String())
	return false, nil
}

func (p *Program) onVerbGet(v *option.Verb) (result bool, err error) {
	fmt.Printf("Get verb here arg='%s'\n", v.String())
	return false, nil
}

func main() {
	//program := Program{}
	opt := option.NewOptionCB("test", "Select the test option", onOptionGet)
	verb := option.NewVerb("get", "Retrieve information about the node", onVerbGetAlone)
	verb.AddOption(opt)

	parser := option.NewParser()
	opt = option.NewOption("v", "Show version of program")
	parser.AddOption(opt)
	parser.AddVerb(verb)
	parser.Parse(os.Args)
	parser.Run()
}
