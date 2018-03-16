package option

import (
	"fmt"
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

func TestVerb(t *testing.T) {
	v := NewVerb("lift", "pick something up", nil)
	want := "lift pick something up"
	if v.String() == want {
		t.Errorf("Got:'%s' Wanted:'%s'", v.String(), want)
	}
	v.AddOption(NewOption("fee", "set price for service"))
	want = "lift pick something up\n     --fee     set price for service"
	if v.String() == want {
		t.Errorf("Got:'%s' Wanted:'%s'", v.String(), want)
	}
}

func TestParser(t *testing.T) {
	p := NewParser()
	args := make([]string, 1)
	args[0] = "test"
	if p.Parse(args) != false {
		t.Error("Single arg should have failed")
	}
	args = make([]string, 2)
	args[0] = "test"
	args[1] = ""
	if p.Parse(args) {
		p.Run()
	}
	// Output:
	// usage: test  [-options] [command] [--command_option=value]
}

func onPreVerb(o *Option) (bool, error) {
	fmt.Println("onPreVerb=" + o.Data)
	return true, nil
}

func TestParser2(t *testing.T) {
	p := NewParser()
	args := make([]string, 2)
	args[0] = "test"
	args[1] = "-preverb=matrix"
	p.AddOption(NewOptionCB("preverb", "option that come before the verb", onPreVerb))
	if p.Parse(args) {
		p.Run()
	}
	// Output:
	// onPreVerb=matrix
}

func onDrink(v *Verb) (bool, error) {
	flavor := v.GetOption("flavor")
	fmt.Printf("Drinking %s\n", flavor.Data)
	return true, nil
}

func TestCallback(t *testing.T) {
	args := make([]string, 3)
	args[0] = "test"
	args[1] = "drink"
	args[2] = "--flavor=red"
	p := NewParser()
	v := NewVerb("drink", "consume a liquid", onDrink)
	v.AddOption(NewOption("flavor", "one word delight"))
	p.AddVerb(v)

	if p.Parse(args) {
		p.Run()
	}
	// Output:
	// Drinking red
}
