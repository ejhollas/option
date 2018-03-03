package main

import (
	"eric/option/option"
	"fmt"
)

func onVerbGet(val string) bool {
	fmt.Println("Get verb here saw " + val)
	return true
}

func main() {
	opt := option.NewOption("test", "Select the test option")
	verb := option.NewVerb("get", "Retrieve information about the node", onVerbGet)
	verb.AddOption(opt)
	fmt.Println(verb)
	verb.OnVerbFound("Bitching")
}
