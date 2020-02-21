package compiler

//ScanLoop scans a loop statement, such as for.
func (c *Compiler) ScanLoop() (err error) {
	c.Loop(nil, func() {
		err = c.CompileBlock()
	})
	return
}
