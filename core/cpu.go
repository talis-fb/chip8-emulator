package main

import (
	"math/rand"
)

const WIDTH uint16 = 64
const HEIGHT uint16 = 32

type Chip8 struct {
	I          uint16
	V          [16]byte
	DelayTimer uint8
	SoundTimer uint8
	Memory     [4096]byte
	PC         uint16
	Stack      [16]uint16
	SP         uint16
	Keyboard   [16]bool
	Display    [64 * 32]bool
}

func (c *Chip8) ExecuteOpcode(opcode uint16) {
	n4 := (opcode & 0xF000) >> 12
	n3 := (opcode & 0x0F00) >> 8
	n2 := (opcode & 0x00F0) >> 4
	n1 := opcode & 0x000F

	switch n4 {
	case 0x0:
		switch opcode {
		case 0x00E0:
			c.execute_00E0()
		case 0x00EE:
			c.execute_00EE()
		}
	case 0x1:
		c.execute_1NNN(n3, n2, n1)
	case 0x2:
		c.execute_2NNN(n3, n2, n1)
	case 0x3:
		c.execute_3XKK(n3, n2, n1)
	case 0x4:
		c.execute_4XKK(n3, n2, n1)
	case 0x5:
		c.execute_5XY0(n3, n2)
	case 0x6:
		c.execute_6XKK(n3, n2, n1)
	case 0x7:
		c.execute_7XKK(n3, n2, n1)
	case 0x8:
		switch n1 {
		case 0x0:
			c.execute_8XY0(n3, n2)
		case 0x1:
			c.execute_8XY1(n3, n2)
		case 0x2:
			c.execute_8XY2(n3, n2)
		case 0x3:
			c.execute_8XY3(n3, n2)
		case 0x4:
			c.execute_8XY4(n3, n2)
		case 0x5:
			c.execute_8XY5(n3, n2)
		case 0x6:
			c.execute_8XY6(n3, n2)
		case 0x7:
			c.execute_8XY7(n3, n2)
		case 0xE:
			c.execute_8XYE(n3, n2)
		}
	case 0x9:
		c.execute_9XY0(n3, n2)
	case 0xA:
		c.execute_ANNN(n3, n2, n1)
	case 0xB:
		c.execute_BNNN(n3, n2, n1)
	case 0xC:
		c.execute_CXKK(n3, n2, n1)
	case 0xD:
		c.execute_DXYN(n3, n2, n1)
	case 0xE:
		n2n1 := (n2 << 4) | n1
		switch n2n1 {
		case 0x9E:
			c.execute_EX9E(n3)
		case 0xA1:
			c.execute_EXA1(n3)
		}
	case 0xF:
		n2n1 := (n2 << 4) | n1
		switch n2n1 {
		case 0x07:
			c.execute_FX07(n3)
		case 0x0A:
			c.execute_FX0A(n3)
		case 0x15:
			c.execute_FX15(n3)
		case 0x18:
			c.execute_FX18(n3)
		case 0x1E:
			c.execute_FX1E(n3)
		case 0x29:
			c.execute_FX29(n3)
		case 0x33:
			c.execute_FX33(n3)
		case 0x55:
			c.execute_FX55(n3)
		case 0x65:
			c.execute_FX65(n3)
		}
	default:
		println("[WASM][ERROR] Unknown opcode: ", opcode)
	}
}

// 00E0 - Clear the display
func (c *Chip8) execute_00E0() {
	c.Display = [64 * 32]bool{}
}

// 00EE - Return from a subroutine
func (c *Chip8) execute_00EE() {
	c.SP -= 1
	c.PC = uint16(c.Stack[c.SP])
}

// 1NNN - Jump to address NNN
func (c *Chip8) execute_1NNN(n3 uint16, n2 uint16, n1 uint16) {
	address := (uint16(n3) << 8) | (uint16(n2) << 4) | uint16(n1)
	c.PC = address
}

// 2NNN - Call subroutine at address NNN
func (c *Chip8) execute_2NNN(n3 uint16, n2 uint16, n1 uint16) {
	address := (uint16(n3) << 8) | (uint16(n2) << 4) | uint16(n1)
	c.Stack[c.SP] = c.PC
	c.SP += 1
	c.PC = address
}

// 3XKK - Skip next instruction if VX == KK
func (c *Chip8) execute_3XKK(X uint16, K2 uint16, K1 uint16) {
	kk := (K2 << 4) | K1
	if c.V[X] == byte(kk) {
		c.PC += 2
	}
}

// 4XKK - Skip next instruction if VX != KK
func (c *Chip8) execute_4XKK(X uint16, K2 uint16, K1 uint16) {
	kk := (K2 << 4) | K1
	if c.V[X] != byte(kk) {
		c.PC += 2
	}
}

// 5XY0 - Skip next instruction if VX == VY
func (c *Chip8) execute_5XY0(X uint16, Y uint16) {
	if c.V[X] == c.V[Y] {
		c.PC += 2
	}
}

// 6XKK - Set VX = KK
func (c *Chip8) execute_6XKK(X uint16, K2 uint16, K1 uint16) {
	kk := (K2 << 4) | K1
	c.V[X] = byte(kk)
}

// 7XKK - Set VX += KK
func (c *Chip8) execute_7XKK(X uint16, K2 uint16, K1 uint16) {
	kk := (K2 << 4) | K1
	c.V[X] += byte(kk)
}

// 8XY0 - Set VX = VY
func (c *Chip8) execute_8XY0(X uint16, Y uint16) {
	c.V[X] = byte(c.V[Y])
}

// 8XY1 - Set VX |= VY
func (c *Chip8) execute_8XY1(X uint16, Y uint16) {
	c.V[X] |= byte(c.V[Y])
}

// 8XY2 - Set VX &= VY
func (c *Chip8) execute_8XY2(X uint16, Y uint16) {
	c.V[X] &= byte(c.V[Y])
}

// 8XY3 - Set VX ^= VY
func (c *Chip8) execute_8XY3(X uint16, Y uint16) {
	c.V[X] ^= byte(c.V[Y])
}

// 8XY4 - Set VX += VY, set VF = carry
func (c *Chip8) execute_8XY4(X uint16, Y uint16) {
	result := uint16(c.V[X]) + uint16(c.V[Y])
	if result > 255 {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[X] = byte(result)
}

// 8XY5 - Set VX -= VY, set VF = NOT borrow
func (c *Chip8) execute_8XY5(X uint16, Y uint16) {
	if c.V[X] >= c.V[Y] {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[X] -= c.V[Y]
}

// 8XY6 - Set VX >>= 1
func (c *Chip8) execute_8XY6(X uint16, _ uint16) {
	c.V[0xF] = c.V[X] & 1
	c.V[X] >>= 1
}

// 8XY7 - Set VX = VY - VX, set VF = NOT borrow
func (c *Chip8) execute_8XY7(X uint16, Y uint16) {
	if c.V[Y] >= c.V[X] {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[X] = c.V[Y] - c.V[X]
}

// 8XYE - Set VX <<= 1
func (c *Chip8) execute_8XYE(X uint16, _ uint16) {
	c.V[0xF] = c.V[X] & 0x80
	c.V[X] <<= 1
}

// 9XY0 - Skip next instruction if VX != VY
func (c *Chip8) execute_9XY0(X uint16, Y uint16) {
	if c.V[X] != c.V[Y] {
		c.PC += 2
	}
}

// ANNN - Set I = NNN
func (c *Chip8) execute_ANNN(n3 uint16, n2 uint16, n1 uint16) {
	c.I = (uint16(n3) << 8) | (uint16(n2) << 4) | uint16(n1)
}

// BNNN - Jump to address V0 + NNN
func (c *Chip8) execute_BNNN(n3 uint16, n2 uint16, n1 uint16) {
	address := (uint16(n3) << 8) | (uint16(n2) << 4) | uint16(n1)
	c.PC = uint16(c.V[0]) + address
}

// CXKK - Set VX = random() & KK
func (c *Chip8) execute_CXKK(X uint16, K2 uint16, K1 uint16) {
	kk := (K2 << 4) | K1
	c.V[X] = byte(rand.Intn(256)) & byte(kk)
}

// DXYN - Draw sprite at coordinates (VX, VY) with N bytes of sprite data
func (c *Chip8) execute_DXYN(X uint16, Y uint16, N uint16) {
	collision := false
	sprite := c.Memory[c.I : c.I+N]

	vx := uint16(c.V[X])
	vy := uint16(c.V[Y])

	for i := 0; i < len(sprite); i++ {
		row := sprite[i]
		for j := 0; j < 8; j++ {
			newPixel := row >> (7 - j) & 0x01
			if newPixel == 1 {
				xi := (vx + uint16(j)) % WIDTH
				yi := (vy + uint16(i)) % HEIGHT
				lastPixel := c.Display[xi+yi*WIDTH]
				if lastPixel {
					collision = true
				}
				c.Display[xi+yi*WIDTH] = lastPixel != (newPixel == 1)
			}
		}
	}

	if collision {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}

}

// EX9E - Skip next instruction if key with the value of VX is pressed
func (c *Chip8) execute_EX9E(X uint16) {
	if c.Keyboard[c.V[X]] {
		c.PC += 2
	}
}

// EXA1 - Skip next instruction if key with the value of VX is not pressed
func (c *Chip8) execute_EXA1(X uint16) {
	if !c.Keyboard[c.V[X]] {
		c.PC += 2
	}
}

// FX07 - Set VX = delay timer
func (c *Chip8) execute_FX07(X uint16) {
	c.V[X] = byte(c.DelayTimer)
}

// FX0A - Wait for a key press and store it in VX
func (c *Chip8) execute_FX0A(X uint16) {
	c.PC -= 2
	for i := 0; i < len(c.Keyboard); i++ {
		if c.Keyboard[i] {
			c.V[X] = byte(i)
			c.PC += 2
			break
		}
	}
}

// FX15 - Set delay timer = VX
func (c *Chip8) execute_FX15(X uint16) {
	c.DelayTimer = byte(c.V[X])
}

// FX18 - Set sound timer = VX
func (c *Chip8) execute_FX18(X uint16) {
	c.SoundTimer = byte(c.V[X])
}

// FX1E - Add VX to I
func (c *Chip8) execute_FX1E(X uint16) {
	c.I += uint16(c.V[X])
}

// FX29 - Set I = location of sprite for character in VX
func (c *Chip8) execute_FX29(X uint16) {
	c.I = uint16(c.V[X]) * 5 // Each sprite is 5 bytes wide
}

// FX33 - Store BCD representation of VX in memory at addresses I, I+1, I+2
func (c *Chip8) execute_FX33(X uint16) {
	vx := c.V[X]
	c.Memory[c.I] = vx / 100
	c.Memory[c.I+1] = (vx / 10) % 10
	c.Memory[c.I+2] = (vx % 100) / 10 //????
}

// FX55 - Store V0 to VX in memory starting at address I
func (c *Chip8) execute_FX55(X uint16) {
	for i := uint16(0); i <= X; i++ {
		c.Memory[c.I+i] = byte(c.V[i])
	}
}

// FX65 - Fill V0 to VX from memory starting at address I
func (c *Chip8) execute_FX65(X uint16) {
	for i := uint16(0); i <= X; i++ {
		c.V[i] = byte(c.Memory[c.I+i])
	}
}
