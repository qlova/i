//+build wasm

package i

import (
	"bufio"
	"syscall/js"
)

var Stdin = bufio.NewReader(jsStdin{})

type jsStdin struct{}

func (stdin jsStdin) Read(p []byte) (n int, err error) {
	var channel = make(chan string)
	go func() {
		var callback js.Func
		callback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			channel <- args[0].String()
			callback.Release()
			return nil
		})
		js.Global().Get("stdin").Call("Read", len(p), callback)
	}()

	var s = <-channel
	copy(p, s)
	return len(s), nil

}
