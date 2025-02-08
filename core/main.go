// This setup VsCode LSP to GO + WASM
//go:build js && wasm

package main

import (
	"syscall/js"
)

var romData [4096]byte

func setRom(this js.Value, args []js.Value) interface{} {
	println("[WASM] setRom called")

	jsArray := args[0]
	tempRomData := make([]byte, jsArray.Length())

	js.CopyBytesToGo(tempRomData, jsArray)

	for i := 0; i < len(tempRomData); i++ {
		romData[i+0x200] = tempRomData[i]
	}

	println("ROM loaded:", len(romData), "bytes")

	return nil
}

var drawCallback js.Value

func onDraw(this js.Value, args []js.Value) interface{} {
	println("[WASM] onDraw called")
	if args[0].Type() == js.TypeFunction {
		drawCallback = args[0]
	}

	// screenBuffer := []byte{0x00, 0x01, 0xFF, 0x02}

	// jsArray := js.Global().Get("Uint8Array").New(len(screenBuffer))
	// for i, pixel := range screenBuffer {
	// 	jsArray.SetIndex(i, pixel)
	// }

	// if drawCallback.Type() == js.TypeFunction {
	// 	drawCallback.Invoke(jsArray)
	// }

	return nil
}

var soundCallback js.Value

func onSound(this js.Value, args []js.Value) interface{} {
	println("[WASM] onSound called")
	if args[0].Type() == js.TypeFunction {
		soundCallback = args[0]
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

	js.Global().Set("setRom", js.FuncOf(setRom))

	js.Global().Set("onDraw", js.FuncOf(onDraw))

	js.Global().Set("onSound", js.FuncOf(onSound))

	js.Global().Set("setKey", js.FuncOf(setKey))

	js.Global().Set("reset", js.FuncOf(reset))

	js.Global().Set("cycle", js.FuncOf(cycle))

	<-c // Block the Go runtime from exiting
}
