package Concept

import "github.com/qlova/script"
import "github.com/qlova/script/compiler"
import "github.com/qlova/i/syntax/errors"

//Format hacks.
import (
	"os"
	"runtime/debug"
	"fmt"
	"bytes"
	//"strings"
)

var Name = compiler.Translatable{
	compiler.English: "concept",
}

//Urggy this is kinda hacky.
type Data struct {
	Children []*compiler.Function
}

func Overload(Function *compiler.Function, Arguments []compiler.Type) *compiler.Function {
	if len(Function.Tokens) == 0 {
		return Function
	}
	
	if Function.Data == nil {
		Function.Data = &Data{}
	}
	var data = Function.Data.(*Data)

	Function.Arguments = Arguments
	
	here:
	for _, f := range data.Children {
		
		if len(f.Arguments) != len(Arguments) {
			continue
		}
		
		for i := range f.Arguments {
			if !f.Arguments[i].Equals(Arguments[i]) {
				continue here
			}
		}
		
		return f
	}
	
	var f compiler.Function
	f = *Function
	f.Arguments = Arguments
	f.Name = Function.Name
	
	
	//Generate a name for the function using underscores to differentiate from user functions.
	for i := range f.Name {
		
		var args_string string
		for j := range Arguments {
			args_string += Arguments[j].Type.Name()
		}
		
		f.Name[i] = f.Name[i]+"_"+args_string
	}
	
	data.Children = append(data.Children, &f)
	
	return &f
}

func ScanCall(c *compiler.Compiler, Function *compiler.Function) *compiler.Type {
	var Arguments []compiler.Type
				
	c.Expecting("(")
	
	for i := 0; i < len(Function.Tokens); i++ {
		Arguments = append(Arguments, c.ScanExpression())
		if i < len(Function.Tokens)-1 {
			c.Expecting(",")
		}
	}
	
	c.Expecting(")")

	var overloaded = Overload(Function, Arguments)
	
	var old = c.CurrentFunction 
	
	//Hacky, way to get overloaded function modified by the return statement...
	c.CurrentFunction = overloaded
	
	//Legacy method to call function.Compile()
	c.Call(overloaded)
	
	var ScriptTypeArguments []script.Type
	if len(Arguments) > 0 {
		for i := range Arguments {
			ScriptTypeArguments = append(ScriptTypeArguments, Arguments[i].Type)
		}
	}
	
	var ReturnType script.Type
	if len(overloaded.Returns) > 0 {
		ReturnType = overloaded.Returns[0].Type
		
		//TODO deal with arguments.
		return compiler.ScriptType(c.Script.Call(overloaded.Name[0], ScriptTypeArguments, ReturnType))
	} else {
		c.Script.Run(overloaded.Name[0], ScriptTypeArguments)
	}

	c.CurrentFunction = old
	
	if len(overloaded.Returns) > 0 {
		return &overloaded.Returns[0]
	} else {
		return nil
	}
	return nil
}

var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		var ConceptName = c.Scan()
		if errors.IsInvalidName(ConceptName) {
			c.RaiseError(errors.InvalidName(ConceptName))
		}
		
		//Create a function
		var Function compiler.Function
		Function.Name = compiler.NoTranslation(ConceptName)
		
		//The arguments of the concept need to be stored.
		c.Expecting("(")
		
		if c.Peek() != ")" {
			Function.Tokens = append(Function.Tokens, c.Scan())
			for {
				if c.Peek() == "," {
					c.Scan()
					Function.Tokens = append(Function.Tokens, c.Scan())
				} else {
					c.Expecting(")")
					break
				}
			}
		} else {
			c.Scan()
		}
		
		//Store the body of the concept inside a compiler Cache.
		var filename = c.Scanners[len(c.Scanners)-1].Filename
		var line = c.Scanners[len(c.Scanners)-1].Line-1
		
		var cache compiler.Cache
		if c.ScanIf(":") {
			cache = c.NewCache("", "\n")
		} else {
			cache = c.NewCache("", "}")
		}
		
		//This actually creates a function.
		Function.Compile = func(c *compiler.Compiler) {
			
			//Pull arguments...
			var ScriptTypeArguments []script.Type
			if len(Function.Arguments) > 0 {
				for i := range Function.Arguments {
					ScriptTypeArguments = append(ScriptTypeArguments, c.Script.Raw(Function.Arguments[i].Type, Function.Tokens[i]) )
				}
			}
			
			var args_string string
			if len(ScriptTypeArguments) > 0 {
				args_string = "_"
			}
			for j := range ScriptTypeArguments {
				args_string += ScriptTypeArguments[j].Name()
			}
			
			CompileCode := func(script *script.Script) {
				
				var old = c.Script
				c.Script = script
				
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
							
							if len(Function.Arguments) > 0 {
								c.Scanners[len(c.Scanners)-1].Line = line-c.LineOffset
								c.RaiseError(compiler.Translatable{
									compiler.English: "Cannot pass those arguments to '"+Function.Name[0]+
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

				c.GainScope()
				
				//Set argument variables here:
				for i := range ScriptTypeArguments {
					c.SetVariable(Function.Tokens[i], *compiler.ScriptType(ScriptTypeArguments[i]))
				}
				
				c.CompileCache(cache, filename, line)
				c.LoseScope()

				c.Script = old
				
			}
			
			//We need to figure out the return types.. So.. lets do some crazy stuff here.
			var FakeScript = c.Script.Fake()
				CompileCode(FakeScript)
			
			var ReturnTypes []script.Type
			if len(c.CurrentFunction.Returns) > 0 {
				ReturnTypes = append(ReturnTypes, c.CurrentFunction.Returns[0].Type)
			}
			
			c.Script.Function(ConceptName+args_string, CompileCode, ScriptTypeArguments, ReturnTypes)
		}
		
		c.RegisterFunction(&Function)
	},
	
	Detect: func(c *compiler.Compiler) bool {		
		for _, f := range c.Functions {
			if f.Name[c.Language] == c.Token() {
				
				ScanCall(c, f)
				
				//Deal with returns...
				
				return true
			}
		}
		
		return false
	},
}

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) *compiler.Type {
		
		for _, Function := range c.Functions {
			if Function.Name[c.Language] == c.Token() {
				
				//More hacky legacy trash.
				var ScriptTypeReturns []script.Type
				if len(Function.Returns) > 0 {
					for i := range Function.Returns {
						ScriptTypeReturns = append(ScriptTypeReturns, c.Script.Raw(Function.Returns[i].Type, "") )
					}
				}
				
				if c.Peek() != "(" {
					return compiler.ScriptType(c.Script.ScopeFunctionType(c.Token(), nil, ScriptTypeReturns))
				} else {
					
					var ret = ScanCall(c, Function)
					
					return compiler.ScriptType(ret.Type)
				}
				
				return nil
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
		
		if c.Peek() == "\n" {
			c.Script.Return(nil)
			return
		}
		var t = c.ScanExpression()
		c.Script.Return(t.Type)
		
		if len(c.CurrentFunction.Returns) == 0 {
			c.CurrentFunction.Returns = append(c.CurrentFunction.Returns, t)
		} else {
			if !c.CurrentFunction.Returns[len(c.CurrentFunction.Returns)-1].Equals(t) {
				c.RaiseError(errors.Inconsistent(t, c.CurrentFunction.Returns[len(c.CurrentFunction.Returns)-1]))
			}
		}
	},
}
