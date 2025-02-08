package main

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
	Display    [64 * 32]byte
}

func (c *Chip8) ExecuteOpcode(opcode uint16) {
	println("Executing opcode:", opcode)

	n4 := (opcode & 0xF000) >> 12
	n3 := (opcode & 0x0F00) >> 8
	n2 := (opcode & 0x00F0) >> 4
	n1 := opcode & 0x000F

	switch n4 {
	case 0x0:
		switch opcode {
		case 0x00E0:
			c.Execute_00E0()
		case 0x00EE:
			c.Execute_00EE()
		}
	case 0x1:
		c.Execute_1NNN(n3, n2, n1)
	case 0x2:
		c.Execute_2NNN(n3, n2, n2)
	case 0x3:
		c.Execute_3XKK(n3, n2, n1)
	case 0x4:
		c.Execute_4XKK(n3, n2, n1)
	case 0x5:
		c.Execute_5XY0(n3, n2, n1)
	case 0x6:
		c.Execute_6XKK(n3, n2, n1)
	case 0x7:
		c.Execute_7XKK(n3, n2, n1)
	case 0x8:
		switch n1 {
		case 0x0:
			c.Execute_8XY0(n3, n2, n1)
		case 0x1:
			c.Execute_8XY1(n3, n2, n1)
		case 0x2:
			c.Execute_8XY2(n3, n2, n1)
		case 0x3:
			c.Execute_8XY3(n3, n2, n1)
		case 0x4:
			c.Execute_8XY4(n3, n2, n1)
		case 0x5:
			c.Execute_8XY5(n3, n2, n1)
		case 0x6:
			c.Execute_8XY6(n3, n2, n1)
		case 0x7:
			c.Execute_8XY7(n3, n2, n1)
		case 0xE:
			c.Execute_8XYE(n3, n2, n1)
		}
	case 0x9:
		c.Execute_9XY0(n3, n2, n1)
	case 0xA:
		c.Execute_ANNN(n3, n2, n1)
	case 0xB:
		c.Execute_BNNN(n3, n2, n1)
	case 0xC:
		c.Execute_CXKK(n3, n2, n1)
	case 0xD:
		c.Execute_DXYN(n3, n2, n1)
	case 0xE:
		switch n1 {
		case 0x9:
			c.Execute_EX9E(n3, n2, n1)
		case 0xA:
			c.Execute_EXA1(n3, n2, n1)
		}
	case 0xF:
		switch n1 {
		case 0x1:
			c.Execute_FX07(n3, n2, n1)
		case 0x2:
			c.Execute_FX0A(n3, n2, n1)
		case 0x3:
			c.Execute_FX15(n3, n2, n1)
		case 0x4:
			c.Execute_FX18(n3, n2, n1)
		case 0x5:
			c.Execute_FX1E(n3, n2, n1)
		case 0x6:
			c.Execute_FX29(n3, n2, n1)
		case 0x7:
			c.Execute_FX33(n3, n2, n1)
		case 0x8:
			c.Execute_FX55(n3, n2, n1)
		case 0x9:
			c.Execute_FX65(n3, n2, n1)
		}
	}
}

// 00E0 - Clear the display
func (c *Chip8) Execute_00E0() {
	c.Display = [64 * 32]byte{} // Reset the display buffer
	c.PC += 2
}

// 00EE - Return from a subroutine
func (c *Chip8) Execute_00EE() {
	c.SP--
	c.PC = c.Stack[c.SP]
}

// 1NNN - Jump to address NNN
func (c *Chip8) Execute_1NNN(n3 uint16, n2 uint16, n1 uint16) {
	address := (uint16(n3) << 8) | (uint16(n2) << 4) | uint16(n1)
	c.PC = address
}

// 2NNN - Call subroutine at address NNN
func (c *Chip8) Execute_2NNN(n3 uint16, n2 uint16, n1 uint16) {
	address := (uint16(n3) << 8) | (uint16(n2) << 4) | uint16(n1)
	c.Stack[c.SP] = c.PC
	c.SP++
	c.PC = address
}

// 3XKK - Skip next instruction if VX == KK
func (c *Chip8) Execute_3XKK(X uint16, K2 uint16, K1 uint16) {
	kk := (K2 << 4) | K1
	if c.V[X] == kk {
		c.PC += 2
	}
	c.PC += 2
}

// 4XKK - Skip next instruction if VX != KK
func (c *Chip8) Execute_4XKK(X uint16, K2 uint16, K1 uint16) {
	kk := (K2 << 4) | K1
	if c.V[X] != kk {
		c.PC += 2
	}
	c.PC += 2
}

// 5XY0 - Skip next instruction if VX == VY
func (c *Chip8) Execute_5XY0(X uint16, Y uint16, n1 uint16) {
	if c.V[X] == c.V[Y] {
		c.PC += 2
	}
	c.PC += 2
}

// 6XKK - Set VX = KK
func (c *Chip8) Execute_6XKK(X uint16, K2 uint16, K1 uint16) {
	kk := (K2 << 4) | K1
	c.V[X] = byte(kk)
	c.PC += 2
}

// 7XKK - Set VX += KK
func (c *Chip8) Execute_7XKK(X uint16, K2 uint16, K1 uint16) {
	kk := (K2 << 4) | K1
	c.V[X] += byte(kk)
	c.PC += 2
}

// 8XY0 - Set VX = VY
func (c *Chip8) Execute_8XY0(X uint16, Y uint16, n1 uint16) {
	c.V[X] = c.V[Y]
	c.PC += 2
}

// 8XY1 - Set VX |= VY
func (c *Chip8) Execute_8XY1(X uint16, Y uint16, n1 uint16) {
	c.V[X] |= c.V[Y]
	c.PC += 2
}

// 8XY2 - Set VX &= VY
func (c *Chip8) Execute_8XY2(X uint16, Y uint16, n1 uint16) {
	c.V[X] &= c.V[Y]
	c.PC += 2
}

// 8XY3 - Set VX ^= VY
func (c *Chip8) Execute_8XY3(X uint16, Y uint16, n1 uint16) {
	c.V[X] ^= c.V[Y]
	c.PC += 2
}

// 8XY4 - Set VX += VY, set VF = carry
func (c *Chip8) Execute_8XY4(X uint16, Y uint16, n1 uint16) {
	result := uint16(c.V[X]) + uint16(c.V[Y])
	if result > 255 {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[X] = byte(result)
	c.PC += 2
}

// 8XY5 - Set VX -= VY, set VF = NOT borrow
func (c *Chip8) Execute_8XY5(X uint16, Y uint16, n1 uint16) {
	if c.V[X] > c.V[Y] {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[X] -= c.V[Y]
	c.PC += 2
}

// 8XY6 - Set VX >>= 1
func (c *Chip8) Execute_8XY6(X uint16, Y uint16, n1 uint16) {
	c.V[0xF] = c.V[X] & 1
	c.V[X] >>= 1
	c.PC += 2
}

// 8XY7 - Set VX = VY - VX, set VF = NOT borrow
func (c *Chip8) Execute_8XY7(X uint16, Y uint16, n1 uint16) {
	if c.V[Y] > c.V[X] {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[X] = c.V[Y] - c.V[X]
	c.PC += 2
}

// 8XYE - Set VX <<= 1
func (c *Chip8) Execute_8XYE(X uint16, Y uint16, n1 uint16) {
	c.V[0xF] = c.V[X] >> 7
	c.V[X] <<= 1
	c.PC += 2
}

// 9XY0 - Skip next instruction if VX != VY
func (c *Chip8) Execute_9XY0(X uint16, Y uint16, n1 uint16) {
	if c.V[X] != c.V[Y] {
		c.PC += 2
	}
	c.PC += 2
}

// ANNN - Set I = NNN
func (c *Chip8) Execute_ANNN(n3 uint16, n2 uint16, n1 uint16) {
	c.I = (uint16(n3) << 8) | (uint16(n2) << 4) | uint16(n1)
	c.PC += 2
}

// BNNN - Jump to address V0 + NNN
func (c *Chip8) Execute_BNNN(n3 uint16, n2 uint16, n1 uint16) {
	address := (uint16(n3) << 8) | (uint16(n2) << 4) | uint16(n1)
	c.PC = uint16(c.V[0]) + address
}

// CXKK - Set VX = random() & KK
func (c *Chip8) Execute_CXKK(X uint16, K2 uint16, K1 uint16) {
	kk := (K2 << 4) | K1
	c.V[X] = byte(rand.Intn(256)) & byte(kk)
	c.PC += 2
}

// DXYN - Draw sprite at coordinates (VX, VY) with N bytes of sprite data
func (c *Chip8) Execute_DXYN(X uint16, Y uint16, N uint16) {
	// Assuming the sprite drawing function here
	// We'll need to implement the logic to draw the sprite on the display and set VF for collision detection.
	c.PC += 2
}

// EX9E - Skip next instruction if key with the value of VX is pressed
func (c *Chip8) Execute_EX9E(X uint16, Y uint16, n1 uint16) {
	if c.Keyboard[c.V[X]] {
		c.PC += 2
	}
	c.PC += 2
}

// EXA1 - Skip next instruction if key with the value of VX is not pressed
func (c *Chip8) Execute_EXA1(X uint16, Y uint16, n1 uint16) {
	if !c.Keyboard[c.V[X]] {
		c.PC += 2
	}
	c.PC += 2
}

// FX07 - Set VX = delay timer
func (c *Chip8) Execute_FX07(X uint16, Y uint16, n1 uint16) {
	c.V[X] = c.DelayTimer
	c.PC += 2
}

// FX0A - Wait for a key press and store it in VX
func (c *Chip8) Execute_FX0A(X uint16, Y uint16, n1 uint16) {
	// Wait for a key press, and store the key in VX
	// We'll need to implement waiting for keypress and storing the result
	c.PC += 2
}

// FX15 - Set delay timer = VX
func (c *Chip8) Execute_FX15(X uint16, Y uint16, n1 uint16) {
	c.DelayTimer = c.V[X]
	c.PC += 2
}

// FX18 - Set sound timer = VX
func (c *Chip8) Execute_FX18(X uint16, Y uint16, n1 uint16) {
	c.SoundTimer = c.V[X]
	c.PC += 2
}

// FX1E - Add VX to I
func (c *Chip8) Execute_FX1E(X uint16, Y uint16, n1 uint16) {
	c.I += uint16(c.V[X])
	c.PC += 2
}

// FX29 - Set I = location of sprite for character in VX
func (c *Chip8) Execute_FX29(X uint16, Y uint16, n1 uint16) {
	c.I = uint16(c.V[X]) * 5 // Each sprite is 5 bytes wide
	c.PC += 2
}

// FX33 - Store BCD representation of VX in memory at addresses I, I+1, I+2
func (c *Chip8) Execute_FX33(X uint16, Y uint16, n1 uint16) {
	value := c.V[X]
	c.Memory[c.I] = value / 100
	c.Memory[c.I+1] = (value / 10) % 10
	c.Memory[c.I+2] = value % 10
	c.PC += 2
}

// FX55 - Store V0 to VX in memory starting at address I
func (c *Chip8) Execute_FX55(X uint16, Y uint16, n1 uint16) {
	for i := uint16(0); i <= X; i++ {
		c.Memory[c.I+i] = c.V[i]
	}
	c.PC += 2
}

// FX65 - Fill V0 to VX from memory starting at address I
func (c *Chip8) Execute_FX65(X uint16, Y uint16, n1 uint16) {
	for i := uint16(0); i <= X; i++ {
		c.V[i] = c.Memory[c.I+i]
	}
	c.PC += 2
}
