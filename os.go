//+build !wasm

package i

import (
	"bufio"
	"os"
)

//Stdin is the stdin input that 'i' will read from.
var Stdin = bufio.NewReader(os.Stdin)
