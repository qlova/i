package main

import (
	"fmt"
	"os"

	"github.com/qlova/i/compiler"
	"github.com/qlova/script"
	"github.com/qlova/script/as"
)

func main() {
	var c = compiler.New()

	var golang bool
	for _, arg := range os.Args {
		switch arg {
		case "go":
			golang = true
		}
	}

	if len(os.Args) > 2 {
		c.Directory = os.Args[len(os.Args)-1]
	}

	var program = func(q script.Ctx) {
		c.Ctx = q

		err := c.Compile()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if golang {
		os.Stdout.Write(as.Go(program))
	} else {
		script.Execute(program)
	}
}
