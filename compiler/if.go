package compiler

//ScanBranch scans an if statement.
func (c *Compiler) ScanBranch() error {
	condition, err := c.ScanExpression()
	if err != nil {
		return err
	}

	//Switch style branches.
	if condition.Type.Is(Switch) {
		var subtype = condition.Type.Subtype()
		var switching = condition.Literal.(Expression)

		//flag specifiying if this is a compiletime switch branch.
		var compiletime = false

		//var chain usm.ElseIf
		//var last usm.Block

		for {
			var next = c.Scan()
			if !next.Is("||") && !next.Is("|") {
				return c.NewError("expecting || or |")
			}

			test, err := c.ScanExpression()
			if err != nil {
				return err
			}

			if test.Type != subtype {
				return c.NewError("cannot compare ", test.Type, " with ", subtype)
			}

			test, err = switching.Equals(test, c)
			if err != nil {
				return err
			}

			if compiletime {
				_, err = c.CacheBlock()
				if err != nil {
					return err
				}
				continue
			}

			if test.Literal == true {
				err = c.CompileBlock()
				if err != nil {
					return err
				}
				compiletime = true
			} else {
				panic("todo")
			}
		}
	}

	if condition.Type != Bit {
		return c.NewError("if statement requires a 'bit' expression, not " + condition.Type)
	}

	var errors []error

	//Cache blocks.
	first, err := c.CacheBlock()
	if err != nil {
		return err
	}

	c.If(condition.Value, func() {
		c.GainScope()
		errors = append(errors, c.CompileCache(first))
		c.LoseScope()
	}, nil, nil)

	return nil
}
