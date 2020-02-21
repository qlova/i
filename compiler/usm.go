package compiler

//ScanInlineAssembly scans an inline assembly statement.
func (c *Compiler) ScanInlineAssembly() error {
	if !c.ScanIf('.') {
		return c.NewError("Expecting inline assembly")
	}

	switch string(c.Scan()) {
	case "read":
		args, err := c.ScanCall()
		if err != nil {
			return err
		}

		if len(args) != 1 || args[0].Type != Data {
			return c.NewError("usm.read takes 1 data argument.")
		}
		c.Discard(c.Read(nil, args[0].Value))
		return nil
	default:
		return c.NewError("Unknown instruction: ", c.Token())

	}
}
