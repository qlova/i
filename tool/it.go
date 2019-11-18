package main

import (
	"fmt"
	"os"

	"github.com/qlova/i/compiler"
	"github.com/qlova/usm/target/golang"
	"github.com/qlova/usm/target/runtime"
)

type Runnable interface {
	Run() error
}

func main() {
	var c = compiler.New()

	c.Target = new(runtime.Target)

	for _, arg := range os.Args {
		switch arg {
		case "go":
			c.Target = new(golang.Target)
		}
	}

	if len(os.Args) > 2 {
		c.Directory = os.Args[len(os.Args)-1]
	}

	err := c.Compile()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c.Target.(Runnable).Run()
}
