package If

import "github.com/qlova/script/compiler"
import ."github.com/qlova/script"

var Statement = compiler.Statement{
	Name: compiler.Translatable{
		compiler.English: "if",
	},
	
	OnScan: func(c *compiler.Compiler) {
		var condition = c.ScanType(Boolean{}).(Boolean)
		
		c.If(condition, func(q *Script) {
			c.GainScope()
			if c.ScanIf(":") {
				c.ScanStatement()
			} else {
				c.CompileBlock("", "}")
			}
			c.LoseScope()
		})
	},
}
