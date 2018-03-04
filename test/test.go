package main

import (
	"eric/option/option"
	"fmt"
	"os"
)

func onVerbGet(val string) bool {
	fmt.Println("Get verb here saw " + val)
	return true
}

func main() {
	opt := option.NewOption("test", "Select the test option")
	verb := option.NewVerb("get", "Retrieve information about the node", onVerbGet)
	verb.AddOption(opt)

	parser := option.NewParser()
	opt = option.NewOption("v", "Show version of program")
	parser.AddOption(opt)
	parser.AddVerb(verb)
	parser.Parse(os.Args)
	parser.Run()
}
