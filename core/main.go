// This setup VsCode LSP to GO + WASM
//go:build js && wasm

package main

import (
	"fmt"
	"sync"
	"syscall/js"
)

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

var chip8 Chip8

var m sync.Mutex

// it assumes the client is running it 60Hz
func cycle(this js.Value, args []js.Value) interface{} {
	println("[WASM] Cycle")

	c := make(chan interface{}, 0)

	// Probaly this is not the best way to do it
	go func() {
		m.Lock()
		defer m.Unlock()

		opcode := uint16(chip8.Memory[chip8.PC])<<8 | uint16(chip8.Memory[chip8.PC+1])

		hexStr := fmt.Sprintf("%04X", opcode)
		println("[WASM] Executing opcode:", hexStr, " | PC:", chip8.PC)

		chip8.ExecuteOpcode(opcode)

		if (opcode & 0xF000) == 0xD000 {
			invokeDrawCallback()
		}

		if chip8.SoundTimer > 0 {
			soundCallback.Invoke()
		}

		chip8.PC += 2

		c <- opcode
	}()

	return <-c
}

func setRom(this js.Value, args []js.Value) interface{} {
	println("[WASM] setRom called")

	jsArray := args[0]
	romData := make([]byte, jsArray.Length())

	js.CopyBytesToGo(romData, jsArray)

	for i := 0; i < len(romData); i++ {
		chip8.Memory[i+0x200] = romData[i]
	}

	chip8.PC = 0x200

	println("ROM loaded:", len(romData), "bytes")

	return nil
}

var drawCallback js.Value

func onDraw(this js.Value, args []js.Value) interface{} {
	println("[WASM] onDraw called")
	if args[0].Type() == js.TypeFunction {
		drawCallback = args[0]
	}
	return nil
}

func invokeDrawCallback() {
	screenBuffer := chip8.Display

	jsArray := js.Global().Get("Uint8Array").New(len(screenBuffer))
	for i, pixel := range screenBuffer {
		jsArray.SetIndex(i, pixel)
	}

	if drawCallback.Type() == js.TypeFunction {
		drawCallback.Invoke(jsArray)
	}
}

var soundCallback js.Value

func onSound(this js.Value, args []js.Value) interface{} {
	println("[WASM] onSound called")
	if args[0].Type() == js.TypeFunction {
		soundCallback = args[0]
	}

	return nil
}

func setKey(this js.Value, args []js.Value) interface{} {
	println("[WASM] setKey called")

	key := args[0].Int()
	isPressed := args[1].Bool()

	println(key, isPressed)

	chip8.Keyboard[key] = isPressed

	return nil
}

func reset(this js.Value, args []js.Value) interface{} {
	println("[WASM] reset called")
	chip8 = Chip8{}
	return nil
}
