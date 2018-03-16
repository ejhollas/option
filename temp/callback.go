package main

import (
	"fmt"
)

type Context struct {
}

type Handler interface {
	OnEvent(Context)
}

type HandlerFunc func(Context)

type Program struct {
	callback Handler
}

func (p *Program) OnEvent(c Context) {
	fmt.Println("OnEvent called")
}

func (f HandlerFunc) OnEvent(c Context) {
	f(c)
}

func (p Program) Run() {
	c := Context{}
	if p.callback != nil {
		p.callback.OnEvent(c)
	}
}

func main() {
	p := Program{}
	p.callback = HandlerFunc(func(Context) {
		fmt.Println("Anon handler")
	})

	fmt.Println(p)
	p.Run()

	p.callback = HandlerFunc(p.OnEvent)
	p.Run()
}
