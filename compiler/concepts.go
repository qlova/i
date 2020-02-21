package compiler

import (
	"errors"

	"github.com/qlova/usm"
)

//Concept is a generic function definition.
type Concept struct {
	Cache
	Args []Argument
}

//CallConcept calls the given concept.
func (c *Compiler) CallConcept(concept Concept) (expression Expression, err error) {
	if !c.ScanIf('(') {
		err = errors.New("expecting '('")
		return
	}

	//Arguments can be named or sequential.
	var args = make(map[string]Expression, len(concept.Args))

	if !c.ScanIf(')') {
		var i = -1
		for {
			i++

			arg, err := c.ScanExpression()
			if err != nil {
				return expression, err
			}

			args[concept.Args[i].Name] = arg

			if c.ScanIf(')') {
				break
			}

			if !c.ScanIf(',') {
				err = c.NewError("expecting ',' or ')'")
				return expression, err
			}
		}
	}

	var f Function

	var given = args

	f.Label = c.Define(len(concept.Args), func() {
		c.Push(NewPackageCtx())
		c.GainScope()

		//Place the arguments in scope.
		for i, arg := range concept.Args {

			var ArgumentType Type = arg.Default.Type

			if exp, ok := given[arg.Name]; ok {
				ArgumentType = exp.Type
			}

			c.DefineVariable(arg.Name, Variable{
				Type:     ArgumentType,
				Register: usm.Arg(uint(i)),
			})
		}

		err = c.CompileCacheWithoutCtx(concept.Cache)

		if c.Throws {
			f.Throws = true
			c.Throws = false
		}

		c.LoseScope()
		c.Pop()
	})

	//We need to prepare the arguments incase some are missing.
	var prepared []usm.Value

	//Place the arguments in scope.
	for _, arg := range concept.Args {
		if ex, ok := given[arg.Name]; ok {
			prepared = append(prepared, ex.Value)

		} else if arg.Default.Value != nil {
			prepared = append(prepared, arg.Default.Value)

		} else {
			err = c.NewError("missing argument '", arg.Name, "' in concept call")
			return
		}
	}

	expression.Value = c.Call(f.Label, prepared...)

	return
}

//ScanConcept scans and returns a concept.
func (c *Compiler) ScanConcept() (concept Concept, err error) {
	if !c.ScanIf('(') {
		return concept, c.NewError("expecting '('")
	}

	//Scan arguments.
	if !c.ScanIf(')') {
		for {
			var arg Argument

			arg.Name = c.Scan().String()
			if c.ScanIf(':') {
				arg.Default, err = c.ScanExpression()
				if err != nil {
					return concept, err
				}
			}

			concept.Args = append(concept.Args, arg)

			if c.ScanIf(')') {
				break
			}

			if !c.ScanIf(',') {
				return concept, c.NewError("expecting ',' or ')'")
			}
		}
	}

	concept.Cache, err = c.CacheBlock()

	return
}
