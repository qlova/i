package If

import "github.com/qlova/script/compiler"

var Statement = compiler.Statement{
	Name: compiler.Translatable{
		compiler.English: "if",
	},
	
	OnScan: func(c *compiler.Compiler) {
		var condition = c.ScanExpression()
		
		var body compiler.Cache
		//var chain []compiler.Cache
		var end compiler.Cache
		
		var inline bool
		
		if c.ScanIf(":") {
			inline = true
			body = c.NewCache("", "|", "\n")
		} else {
			body = c.NewCache("", "|", "}")
		}
		
		if body.Match == "|" {
			if inline {
				end = c.NewCache("", "|", "\n")
			} else {
				end = c.NewCache("", "}")
			}
		}
		
		c.If(condition.Value().Bool(), func() {
			c.GainScope()
			c.CompileCache(body)
			c.LoseScope()
		}, c.Else(func() {
			c.GainScope()
			c.CompileCache(end)
			c.LoseScope()
		}))
	},
}
