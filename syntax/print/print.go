package Print

import . "github.com/qlova/script"
import "github.com/qlova/script/compiler"
	

var Statement = compiler.Statement{
	Name: compiler.Translatable{
		compiler.English: "print",
	},
	
	OnScan: func(c *compiler.Compiler) {
		c.Expecting("(")
		
		if c.Peek() == ")" {
			c.Scan()
			c.Print()
			return
		}
		
		//Deal with the first argument.
		var Arguments = []Type{
			c.ScanExpression().Value().String(),
		}
		
		//Deal with the subsequent arguments.
		for {
			if c.Peek() == "," {
				c.Scan()
				Arguments = append(Arguments, c.ScanExpression().Value().String())
			} else {
				c.Expecting(")")
				break
			}
		}
		
		c.Print(Arguments...)
	},
}
