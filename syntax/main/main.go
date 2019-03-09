package Main 

import "github.com/qlova/script/compiler"

var Statement = compiler.Statement{
	Name: compiler.Translatable{
		compiler.English: "main",
	},
	
	OnScan: func(c *compiler.Compiler) {
		c.Main(func() {
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
