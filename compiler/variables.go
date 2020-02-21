package compiler

import "github.com/qlova/usm"

//Variable is a variable.
type Variable struct {
	usm.Register
	Type
	Owned bool
}

//DefineVariable defines the variable to be in scope.
//This does not generate any usm.
func (c *Compiler) DefineVariable(name string, v Variable) {
	c.Scope.Variables[name] = v
}

//GetVariable returns the variable in scope with the given name.
func (c *Compiler) GetVariable(name string) (v Variable, ok bool) {
	v, ok = c.Scope.Variables[name]
	if !ok {
		for _, scope := range c.Scopes {
			v, ok = scope.Variables[name]
			if ok {
				return
			}
		}
	}
	return
}
