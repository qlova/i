package String

//import "fmt"
//import "reflect"

import "strconv"
import "github.com/qlova/script/compiler"

var Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) compiler.Type {
	
	//Concatenate strings.
	
	//fmt.Println(reflect.TypeOf(a.LanguageType()).Name(), reflect.TypeOf(b.LanguageType()).Name())
	
	//TODO use IsString methods
	if a.Value().Is(c.String()) && symbol == "+" {
		
		if b.Value().Is(c.String()) {
			return a.Value().String().Add(b.Value().String())
		}
		
	}
	
	
	return nil
}

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) compiler.Type {
		
		//Cast to string.
		if c.Token() == "string" {
			c.Expecting("(")
			var expression = c.ScanExpression()
			c.Expecting(")")
			
			return expression.Value().String()
		}
		
		
		//String Literal.
		if c.Token()[0] == '"' {
			
			str, err := strconv.Unquote(c.Token())
			if err != nil {
				c.RaiseError(compiler.Translatable{
					compiler.English: "Invalid String!",
				})
			}
			
			return c.String(str)
		}
		
		return nil
	},
}
