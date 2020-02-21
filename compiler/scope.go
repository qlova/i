package compiler

//Scope can contain variables and Cleanup routines.
type Scope struct {
	Variables map[string]Variable
	//Cleanups []func()
	//Afters   []func()
}

//NewScope returns a new scope.
func NewScope() Scope {
	return Scope{
		Variables: make(map[string]Variable),
	}
}

//GainScope pushes a new scopeset to the compiler.
func (c *Compiler) GainScope() {
	c.Scopes = append(c.Scopes, NewScope())
}

//LoseScope pops the last scopeset from the compiler.
func (c *Compiler) LoseScope() {
	if len(c.Scopes) == 0 {
		return
	}
	c.Scope = c.Scopes[len(c.Scopes)-1]
	c.Scopes = c.Scopes[:len(c.Scopes)-1]
}
