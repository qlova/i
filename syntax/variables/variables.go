package Variables

import "github.com/qlova/script/compiler"

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) compiler.Type {
		return c.GetVariable(c.Token())
	},
}

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		if !c.GetVariable(c.Token()).Defined {
			
			var name = c.Token()
			
			if c.Peek() != "=" {
				return false
			}
			
			c.Expecting("=")
			c.SetVariable(name, c.ScanExpression().Value().Var(name))
			
			return true
			
		} else {
			
			var name = c.Token()
			var variable = c.Variable(name)
			
			//Functions
			/*if c.Peek() == "(" && variable.Equals(qlova.FunctionType{}) {
				c.Scan()
				c.Expecting(")")
				
				c.Script.RunFunctionType(name, nil)
				
				return true
			}*/
			
			if c.ScanIf("=") {
				variable.Value().Set(c.ScanType(variable))
				return true
			}
		}
		
		return false
	},
}
