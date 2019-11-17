package compiler

import "github.com/qlova/i/compiler/scanner"

type Token = scanner.Token

func s(s string) []byte {
	return []byte(s)
}
