package compiler

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

var Trace = os.Getenv("TRACE") != ""
var Panic = os.Getenv("PANIC") != ""
var Counter = 2

type Error struct {
	Formatted string
	Message   string
}

func (err Error) Error() string {
	return err.Formatted
}

func (c *Compiler) NewError(format ...interface{}) error {
	var msg = fmt.Sprint(format...)

	var wdir string
	if runtime.GOOS != "js" {
		wdir, _ = os.Getwd()
	}

	var rpath, _ = filepath.Rel(wdir, c.Directory)

	if len(rpath) > len(c.Directory) {
		rpath = c.Directory
	}

	var RestOfTheLine string

	if c.Column == 0 {
		RestOfTheLine = string(c.LastLine)
		c.Column = len(c.LastLine)
	} else {
		RestOfTheLine, _ = c.Reader.ReadString('\n')
	}

	var formatted = fmt.Sprint(rpath, c.Filename, ":",
		c.LineNumber, ": ", string(c.Line), RestOfTheLine, "\n",
		strings.Repeat(" ", c.Column+2+len(rpath)+len(c.Filename)), "^\n", msg)

	if Trace {
		var stacktrace = debug.Stack()
		var reader = bufio.NewReader(bytes.NewBuffer(stacktrace))

		const count = 7
		for i := 0; i < count; i++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				formatted += string(stacktrace)
				break
			}
			if i == count-1 {
				formatted += "\n(" + strings.TrimSpace(line) + ")\n"
			}
		}

	}

	if Panic {
		Counter--
		if Counter == 0 {
			panic(formatted)
		}
	}

	return Error{formatted, msg}
}

//Unimplemented is an error describing that the component is unimplemented.
func (c *Compiler) Unimplemented(component []byte) error {
	return c.NewError("unimplemented " + string(component))
}

//Undefined is an error describing that the name is undefined.
func (c *Compiler) Undefined(name []byte) error {
	return c.NewError("undefined " + string(name))
}

//Expecting returns an error in the form "expecting [token]"
func (c *Compiler) Expecting(symbol byte) error {
	return c.NewError("expecting " + string(symbol) + " but found " + c.Scan().String())
}
