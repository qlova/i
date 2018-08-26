package Main 

import "github.com/qlova/script/compiler"
import qlova "github.com/qlova/script"

var Statement = compiler.Statement{
	Name: compiler.Translatable{
		compiler.English: "main",
	},
	
	OnScan: func(c *compiler.Compiler) {
		c.Script.Main(func(q *qlova.Script) {
			c.GainScope()
			if c.ScanIf(":") {
				c.CompileBlock("", "\n")
			} else {
				c.CompileBlock("", "}")
			}
			c.LoseScope()
		})
	},
}
