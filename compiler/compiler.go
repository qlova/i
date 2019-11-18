package compiler

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/qlova/usm"
)

type Compiler struct {
	usm.Target

	//Packages
	PackageCtx
	Stack    []PackageCtx
	Packages map[string]PackageCtx

	//Testing
	ExpectedOutput []byte
	ProvidedInput  []byte
}

//New returns a new initialized compiler.
func New() Compiler {
	var c Compiler
	c.Packages = make(map[string]PackageCtx)
	return c
}

//Compile package located at Compiler.Dir or current working directory if empty.
func (c *Compiler) Compile() error {
	if c.Directory == "" {
		c.Directory = "."
	}

	PackageDirectory, err := os.Open(c.Directory)
	if err != nil {
		return fmt.Errorf("could not open %v: %w", c.Directory, err)
	}

	if info, err := PackageDirectory.Stat(); err != nil {
		return fmt.Errorf("could not stat %v: %w", c.Directory, err)
	} else if !info.IsDir() {
		PackageDirectory.Close()
		return c.CompileFile(c.Directory)
	}

	files, err := ioutil.ReadDir(c.Directory)
	if err != nil {
		return fmt.Errorf("could not read directory: %v (%w)", c.Directory, err)
	}

	for _, file := range files {
		if path.Ext(file.Name()) == ".i" {
			err := c.CompileFile(c.Directory + "/" + file.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//CompileFile compiles a file.
func (c *Compiler) CompileFile(location string) error {
	c.Filename = filepath.Base(location)

	file, err := os.Open(location)
	if err != nil {
		//Return to the last frame.
		if len(c.Stack) > 0 {
			c.Pop()
		}
		return fmt.Errorf("could not open file %v: %w", location, err)
	}
	defer file.Close()

	return c.CompileReader(file)
}

//CompileReader compiles a reader.
func (c *Compiler) CompileReader(reader io.Reader) error {
	if reader == nil {
		return c.NewError("null reader")
	}

	c.SetReader(reader)

	for {
		err := c.CompileStatement()
		if err != nil {

			//Return to the last frame.
			if len(c.Stack) > 0 {
				c.Pop()

				if err != io.EOF {
					return err
				}

				return nil
			} else if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
	}
}

//CompileBlock compiles an 'i' code block.
func (c *Compiler) CompileBlock() error {
	c.GainScope()
	defer c.LoseScope()

	if c.ScanIf(':') {
		if err := c.CompileStatement(); err != nil {
			return err
		}
		return nil
	}

	if !c.ScanIf('\n') {
		return c.NewError("block must start with a newline")
	}

	for {
		if c.Peek().Is("}") {
			c.Scan()
			break
		}

		err := c.CompileStatement()
		if err != nil {
			return err
		}
	}

	return nil
}
