package Number

import . "github.com/qlova/script"
import "github.com/qlova/script/compiler"

import (
	"math/big"
	"strings"
)

var Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) compiler.Type {
	if A, ok := a.(Number); ok  {
		
		if B, ok := b.(Number); ok {
			
			switch symbol {
				case "+":
					return c.Add(A, B)

				case "-":
					return c.Sub(A, B)

				case "*":
					return c.Mul(A, B)

				case "/":
					return c.Div(A, B)
					
				case "%":
					return c.Mod(A, B)
					
				case "^":
					return c.Pow(A, B)
			}
			
			
		}
		
	}
	return nil
}

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) compiler.Type {
		
		if c.Token() == "number" {
			c.Expecting("(")
			var expression = c.ScanExpression()
			c.Expecting(")")
			
			return c.Script.ToNumber(expression)
		}
		
		switch c.Token()[0] {
			
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				
				if strings.Contains(c.Token(), ".") {
					return nil
				}
				
				var b big.Int
				var worked bool
				
				if c.Token()[0] == '0' {
					_, worked = b.SetString(c.Token(), 2)
					
				} else {
				
					_, worked = b.SetString(c.Token(), 10)
				
				}
				
				if !worked {
					if len(c.Token()) > 2 {
						_, worked = b.SetString(c.Token()[2:], 16)
					}
				}
				
				if worked {
					
					/*if c.ScanIf(symbols.Factorial) {
						c.Call(&Factorial)
					}*/
					
					return c.Script.BigNumber(&b)
				} else {
					return nil
				}
				
			default:
				return nil
		}
		
		return nil
	},
}
