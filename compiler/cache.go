package compiler

import (
	"bytes"
	"io"
)

//Cache is a storage container that contains code. It can be compiled at a later point in time.
type Cache struct {
	Data []byte

	Filename   string
	LineNumber int

	compiler *Compiler
}

//CompileCacheWithoutCtx compiles a cache.
func (c *Compiler) CompileCacheWithoutCtx(cache Cache) error {
	c.SetReader(bytes.NewReader(cache.Data))
	c.Filename = cache.Filename
	c.LineNumber = cache.LineNumber

	for {
		err := c.CompileStatement()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

//CompileCache compiles a cache with the new context.
func (c *Compiler) CompileCache(cache Cache) error {
	c.Push(c.PackageCtx)
	defer c.Pop()

	return c.CompileCacheWithoutCtx(cache)
}

//CacheBlock create a cache out of the next 'i' block.
func (c *Compiler) CacheBlock() (cache Cache, err error) {
	var buffer bytes.Buffer

	cache.Filename = c.Filename
	cache.LineNumber = c.LineNumber

	var sort = c.Scan()

	var depth = 1

	if sort.Is(":") {

		var column = c.Column
		for {
			var token = c.Scan()

			if token.Is("\n") {
				buffer.Write(c.LastLine[column:])
				break
			}

			if token == nil {
				err = c.NewError("unexpected eof")
				break
			}
		}
	} else {
		if sort.Is("{") {
			c.Scan()
		}

		cache.LineNumber++
		for {
			if c.Peek().Is("||") || c.Peek().Is("|") {
				if depth == 1 {
					if len(c.Line) > 0 {
						buffer.Write(c.Line[:len(c.Line)-1])
					}
					break
				}
			}

			var token = c.Scan()

			if token == nil {
				err = c.NewError("unexpected eof")
				break
			}

			if token.Is("\n") {
				buffer.Write(c.LastLine)
			}

			if token.Is(":") || token.Is("}") || token == nil {
				depth--
				if depth == 0 {
					if len(c.Line) > 0 {
						buffer.Write(c.Line[:len(c.Line)-1])
					}
					break
				}
			} else {
				switch token.String() {
				case "for", "if", "catch", "try", "{", "main":
					depth++
				}

			}
		}
	}

	cache.Data = buffer.Bytes()

	return
}
