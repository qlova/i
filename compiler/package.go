package compiler

import "github.com/qlova/i/compiler/scanner"

//PackageCtx contains package level context, such as defined symbols, and the package scanner.
type PackageCtx struct {
	scanner.Scanner

	//Scope
	Scope
	Scopes []Scope

	Directory string

	Concepts map[string]Concept

	Throws bool
}

//NewPackageCtx returns a new PackageCtx.
func NewPackageCtx() PackageCtx {
	return PackageCtx{
		Scope:    NewScope(),
		Concepts: make(map[string]Concept),
	}
}

//NewCtx returns a new transparent PackageCtx.
func (c *Compiler) NewCtx() PackageCtx {
	return c.PackageCtx
}

//Pop pops the current PackageCtx of the compiler and returns to the PackageCtx.
func (c *Compiler) Pop() {
	var context = c.Stack[len(c.Stack)-1]
	c.PackageCtx = context
	c.Stack = c.Stack[:len(c.Stack)-1]
}

//Push pushes the specified PackageCtx onto the compiler.
func (c *Compiler) Push(p PackageCtx) {
	c.Stack = append(c.Stack, c.PackageCtx)
	c.PackageCtx = p
}
