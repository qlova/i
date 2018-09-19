package Concept

import . "github.com/qlova/script"
import "github.com/qlova/script/compiler"


var Name = compiler.Translatable{
	compiler.English: "concept",
}

//TODO merge this into Compiler
type Concept struct {
	Arguments []string
	Cache compiler.Cache
}

var Concepts = make(map[string]Concept)

var Functions = make(map[string]Function)

var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		
		var Name = c.Scan()
		c.Expecting("(")

		var Arguments []string
		for c.Peek() != ")" {
			Arguments = append(Arguments, c.Scan())
			
			if c.Peek() != ")" {
				c.Expecting(",")
			}
		}
		c.Expecting(")")
		
		var cache compiler.Cache
		if c.ScanIf(":") {
			cache = c.NewCache("", "\n")
		} else {
			cache = c.NewCache("", "}")
		}
		
		Concepts[Name] = Concept{
			Arguments: Arguments,
			Cache: cache,
		}
	},
	
	Detect: func(c *compiler.Compiler) bool {
		
		var Name = c.Token()
		
		if function, ok := c.Variable(Name).(Function); ok {
			
			if len(function.Arguments()) == 0 {
				c.Expecting("(")
				c.Expecting(")")
				
				function.Call()
				return true
			}
			
			c.Unimplemented()
			
		}
		
		if concept, ok := Concepts[Name]; ok {

			CreateAndCall(c, Name, concept)
			return true
		}
		return false
	},
}

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) compiler.Type {
		var Name = c.Token()
		
		if concept, ok := Concepts[Name]; ok {
			
			if c.Peek() == "(" {
				return CreateAndCall(c, Name, concept)
			} else {
				return Create(c, Name, concept)
			}
		}
		return nil
	},
}

var Return = compiler.Statement {
	Name: compiler.Translatable{
		compiler.English: "return",
	},
	
	OnScan: func(c *compiler.Compiler) {
		c.Return(c.ScanExpression())
	},
}
