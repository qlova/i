package Concept

import "reflect"
import . "github.com/qlova/script"
import "github.com/qlova/script/compiler"

import (
	"fmt"
	"os"
	"bytes"
	"runtime/debug"
)

func Create(c *compiler.Compiler, Name string, concept Concept) Function {
	if len(concept.Arguments) == 0 {
		//This is easy enough to handle.
		if _, ok := Functions[Name]; !ok {
			f := c.Function(func() {
				c.CompileCache(concept.Cache)
			})
			f.Promote(Name) 
			Functions[Name] = f
		}

		return Functions[Name]
	}
	
	panic("concept.Create(): Cannot capture a concept without function arguments.")
}

func CreateAndCall(c *compiler.Compiler, Name string, concept Concept) compiler.Type {
	if len(concept.Arguments) == 0 {
		c.Expecting("(")
		c.Expecting(")")
		
		f := Create(c, Name, concept)
		return f.Call()

	}
	
	//TODO change name depending on types.
	c.Expecting("(")
	
	var Arguments []Type
	var Reflections []reflect.Type
	for c.Peek() != ")" {
		
		var expression = c.ScanExpression()
		Arguments = append(Arguments, expression)
		Reflections = append(Reflections, reflect.TypeOf(expression))
		
		//TODO do fancy stuff like lists, arrays and function subtypes. 
		//Convert them into Go equivilents.
		
		if c.Peek() != ")" {
			c.Expecting(",")
		}
	}
	
	c.Expecting(")")
	
	var ReflectedFunctionType = reflect.FuncOf(Reflections, nil, false)
	
	if _, ok := Functions[Name]; !ok {
		f := c.Function(reflect.MakeFunc(ReflectedFunctionType, func(args []reflect.Value) (results []reflect.Value) {
			c.GainScope()
			for i := range Arguments {
				c.SetVariable(concept.Arguments[i], args[i].Interface().(Type))
			}
			
			var b bytes.Buffer
			c.StdErr = append(c.StdErr, &b)
			var linenumber = c.Scanners[len(c.Scanners)-1].Line
			
			defer func(c *compiler.Compiler, line int) {
				
				c.StdErr = c.StdErr[:len(c.StdErr)-1] // restoring the real stdout
				
				if r := recover(); r != nil {
					if r == "error" {
						c.Errors = true
						
						if os.Getenv("PANIC") == "1" { 
							fmt.Println(b.String())
							panic("PANIC=1")
						}
						
						if len(Arguments) > 0 {
							c.Scanners[len(c.Scanners)-1].Line = line-c.LineOffset
							c.RaiseError(compiler.Translatable{
								compiler.English: "Cannot pass those arguments to '"+Name+
								"' because\n\n"+b.String(),
							})
						}
						fmt.Println(b.String())
						panic("error")
						
					} else if r != "done" {
						fmt.Println(r, string(debug.Stack()))
					}
				}
			}(c, linenumber)
			
			c.CompileCache(concept.Cache)
			c.LoseScope()
			
			return nil
		}).Interface())
		f.NameArguments(concept.Arguments)
		f.Promote(Name)
		
		Functions[Name] = f
	}
	
	f := Functions[Name]
	return f.Call(Arguments...) //TODO deal with returns
}
