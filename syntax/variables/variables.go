package Variables

import "github.com/qlova/script/compiler"
import qlova "github.com/qlova/script"
import "github.com/qlova/i/syntax/errors"

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) *compiler.Type {
		if c.GetVariable(c.Token()).Defined {
			return compiler.ScriptType(c.Script.Raw(c.GetVariable(c.Token()).Type.Type, c.Token()))
		}
		return nil
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
			
			var expression = c.ScanExpression()
			
			c.Script.New(expression.Type, name)
			
			c.SetVariable(name, expression)
			
			return true
		} else {
			
			var name = c.Token()
			var variable = c.GetVariable(name).Type.Type
			
			//Functions
			if c.Peek() == "(" && variable.Equals(qlova.FunctionType{}) {
				c.Scan()
				c.Expecting(")")
				
				c.Script.RunFunctionType(name, nil)
				
				return true
			}
			
			//Arrays
			if c.Peek() == "[" && variable.Equals(qlova.Array{}) {
				c.Scan()
				
				var index = c.ScanExpression()
				if !index.Type.Equals(qlova.Number("")) {
					c.RaiseError(errors.IndexError())
				}
				c.Expecting("]")
				
				name = c.Script.IndexArrayRaw(name, index.Type.String())
				variable = variable.(qlova.Array).Subtype()
			}
			
			c.Expecting("=")
			
			var expression = c.ScanExpression()
			
			if !expression.Type.Equals(variable) {
				c.ExpectingTypeName(variable.Name(), expression.Type.Name())
			}
			
			c.Script.Set(expression.Type, name)
			
			return true
		}
		
		return false
	},
}
