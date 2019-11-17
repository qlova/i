package compiler

import (
	"bytes"
	"os"
)

func (c *Compiler) compileComments(comment Token) error {
	//Special output comment for tests.
	if len(comment) > len("//output: ") && bytes.Equal(comment[:len("//output: ")], []byte("//output: ")) {
		c.ExpectedOutput = comment[len("//output: "):]
		c.ExpectedOutput = bytes.Replace(c.ExpectedOutput, []byte(`\n`), []byte("\n"), -1)
		c.ExpectedOutput = bytes.Replace(c.ExpectedOutput, []byte(`$HOME`), []byte(os.Getenv("HOME")), -1)
		c.ExpectedOutput = bytes.Replace(c.ExpectedOutput, []byte(`$USER`), []byte(os.Getenv("USER")), -1)
		c.ExpectedOutput = bytes.Replace(c.ExpectedOutput, []byte(`$PATH`), []byte(os.Getenv("PATH")), -1)
	}

	//Special output comment for tests.
	if len(comment) > len("//input: ") && bytes.Equal(comment[:len("//input: ")], []byte("//input: ")) {
		c.ProvidedInput = comment[len("//input: "):]
		c.ProvidedInput = bytes.Replace(c.ProvidedInput, []byte(`\n`), []byte("\n"), -1)
	}

	return nil
}
