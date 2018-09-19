package Import

import . "github.com/qlova/script"
import "github.com/qlova/script/compiler"

var Name = compiler.Translatable{
	compiler.English: "import",
}

type Import struct {
	Command string
}

//TODO merge this into Compiler
var Imports = make(map[string]Import)

var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		var Package = c.Scan()
		
		Imports[Package] = Import{
			Command: Package,
		}
	},
	
	Detect: func(c *compiler.Compiler) bool {
		
		if Import, ok := Imports[c.Token()]; ok {
			
			c.Expecting("(")
			
			var Arguments []Type
			
			//Scoop up the arguments and convert them to strings.
			for c.Peek() != ")" {
				Arguments = append(Arguments, c.ToString(c.ScanExpression()))
				if c.Peek() != ")" {
					c.Expecting(",")
				}
			}
			c.Expecting(")")

			c.Trace()
			c.External(c.String(Import.Command)).Run(Arguments...)
			return true
		}
		
		return false
	},
}

var Expression = compiler.Expression {
	Detect: func(c *compiler.Compiler) compiler.Type {
		if Import, ok := Imports[c.Token()]; ok {
			
			c.Expecting("(")
			
			var Arguments []Type
			
			//Scoop up the arguments and convert them to strings.
			for c.Peek() != ")" {
				Arguments = append(Arguments, c.ToString(c.ScanExpression()))
				if c.Peek() != ")" {
					c.Expecting(",")
				}
			}
			c.Expecting(")")
			
			c.Trace()
			return c.External(c.String(Import.Command)).Call(Arguments...)
		}
		
		return nil
	},
}
