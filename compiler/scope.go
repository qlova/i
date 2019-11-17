package compiler

//Scope can contain variables and Cleanup routines.
type Scope struct {
	//Table    map[string]Type
	//Cleanups []func()
	//Afters   []func()
}

//NewScope returns a new scope.
func NewScope() Scope {
	return Scope{}
}

//GainScope pushes a new scopeset to the compiler.
func (c *Compiler) GainScope() {
	c.Scopes = append(c.Scopes, NewScope())
}

//LoseScope pops the last scopeset from the compiler.
func (c *Compiler) LoseScope() {
	c.Scope = c.Scopes[len(c.Scopes)-1]
	c.Scopes = c.Scopes[:len(c.Scopes)-1]
}
