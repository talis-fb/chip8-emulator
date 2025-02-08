// This setup VsCode LSP to GO + WASM
//go:build js && wasm

package main

import (
	"syscall/js"
)

var drawCallback js.Value

func onDraw(this js.Value, args []js.Value) interface{} {
	println("onDraw called")
	if args[0].Type() == js.TypeFunction {
		drawCallback = args[0]
	}

	screenBuffer := []byte{0x00, 0x01, 0xFF, 0x02}

	jsArray := js.Global().Get("Uint8Array").New(len(screenBuffer))
	for i, pixel := range screenBuffer {
		jsArray.SetIndex(i, pixel)
	}

	if drawCallback.Type() == js.TypeFunction {
		drawCallback.Invoke(jsArray)
	}

	return nil
}

var keyboard [16]bool

func setKey(this js.Value, args []js.Value) interface{} {
	println("[WASM] setKey called")

	key := args[0].Int()
	isPressed := args[1].Bool()

	println(key, isPressed)

	keyboard[key] = isPressed

	return nil
}

func reset(this js.Value, args []js.Value) interface{} {
	println("[WASM] reset called")
	return nil
}

func cycle(this js.Value, args []js.Value) interface{} {
	println("[WASM] cycle called")
	return nil
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("helloWasm", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		println("Hello, World from Go WebAssembly!")
		return nil
	}))

	js.Global().Set("onDraw", js.FuncOf(onDraw))

	js.Global().Set("setKey", js.FuncOf(setKey))

	js.Global().Set("reset", js.FuncOf(reset))

	js.Global().Set("cycle", js.FuncOf(cycle))

	<-c // Block the Go runtime from exiting
}
