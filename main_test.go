package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {
	t.Run("00E0: Clears the screen", func(t *testing.T) {
		emu := NewEmulator()
		emu.Screen[0] = true
		emu.Screen[1] = true

		emu.Decode(0x00E0)

		got := false
		want := false

		for _, v := range emu.Screen {
			got = got && v
		}

		assert.Equal(t, want, got)
	})

	t.Run("00EE: Sets program counter to the stack address", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = uint16(0x0001)
		emu.Stack[0] = uint16(0x0010)
		emu.Stack[1] = uint16(0x0020)
		emu.StackPointer = uint16(2)
		emu.Decode(0x00EE)

		assert.Equal(t, uint16(0x0020), emu.ProgramCounter)
		assert.Equal(t, uint16(1), emu.StackPointer)
	})

	t.Run("1nnn: Sets the program counter to nnn", func(t *testing.T) {
		emu := NewEmulator()
		emu.Decode(0x1228)

		got := emu.ProgramCounter
		want := uint16(0x228)
		assert.Equal(t, want, got)
	})

	t.Run("2nnn: Sets the program counter to nnn and adds the previous value to the stack", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = uint16(0x900)
		emu.Decode(0x2228)

		assert.Equal(t, uint16(0x228), emu.ProgramCounter)
		assert.Equal(t, uint16(0x900), emu.Stack[0])
		assert.Equal(t, uint16(1), emu.StackPointer)
	})

	t.Run("3xnn: Skip next instruction if VRegister[x] equals nn", func(t *testing.T) {
		t.Run("VRegister[x] equals nn", func(t *testing.T) {
			emu := NewEmulator()
			emu.ProgramCounter = uint16(0x900)
			emu.VRegisters[2] = uint16(0xFF)
			emu.Decode(0x32FF)

			assert.Equal(t, uint16(0x902), emu.ProgramCounter)
		})

		t.Run("VRegister[x] does not equal nn", func(t *testing.T) {
			emu := NewEmulator()
			emu.ProgramCounter = uint16(0x900)
			emu.VRegisters[2] = uint16(0xAA)
			emu.Decode(0x32FF)

			assert.Equal(t, uint16(0x900), emu.ProgramCounter)
		})
	})

	t.Run("4xnn: Skip next instruction if VRegister[x] does not equal nn", func(t *testing.T) {
		t.Run("VRegister[x] does not equal nn", func(t *testing.T) {
			emu := NewEmulator()
			emu.ProgramCounter = uint16(0x900)
			emu.VRegisters[2] = uint16(0xFF)
			emu.Decode(0x42AA)

			assert.Equal(t, uint16(0x902), emu.ProgramCounter)
		})

		t.Run("VRegister[x] equals nn", func(t *testing.T) {
			emu := NewEmulator()
			emu.ProgramCounter = uint16(0x900)
			emu.VRegisters[2] = uint16(0xFF)
			emu.Decode(0x42FF)

			assert.Equal(t, uint16(0x900), emu.ProgramCounter)
		})
	})

	t.Run("5xy0: Skip next instruction if VRegister[x] equals VRegister[y]", func(t *testing.T) {
		t.Run("VRegister[x] equals VRegister[y]", func(t *testing.T) {
			emu := NewEmulator()
			emu.ProgramCounter = uint16(0x900)
			emu.VRegisters[2] = uint16(0xFF)
			emu.VRegisters[3] = uint16(0xFF)
			emu.Decode(0x5230)

			assert.Equal(t, uint16(0x902), emu.ProgramCounter)
		})

		t.Run("VRegister[x] does not equal VRegister[y]", func(t *testing.T) {
			emu := NewEmulator()
			emu.ProgramCounter = uint16(0x900)
			emu.VRegisters[2] = uint16(0xFF)
			emu.VRegisters[3] = uint16(0xAA)
			emu.Decode(0x5230)

			assert.Equal(t, uint16(0x900), emu.ProgramCounter)
		})
	})

	t.Run("6xnn: Assigns nn to v register x", func(t *testing.T) {
		emu := NewEmulator()
		emu.Decode(0x6105)

		got := emu.VRegisters[1]
		want := uint16(0x05)

		assert.Equal(t, want, got)
	})

	t.Run("7xnn: Adds nn to v register x", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[1] = uint16(10)
		emu.Decode(0x7105)

		got := emu.VRegisters[1]
		want := uint16(15)

		assert.Equal(t, want, got)
	})

	t.Run("8xy0: Set v register x to the value of v register y", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = uint16(0x1)
		emu.VRegisters[3] = uint16(0x2)
		emu.Decode(0x8230)

		assert.Equal(t, uint16(0x2), emu.VRegisters[2])
	})

	t.Run("8xy1: Set VRegister[x] to the value of VRegister[x] OR VRegister[y]", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = uint16(0xcc)
		emu.VRegisters[3] = uint16(0xaa)
		emu.Decode(0x8231)

		assert.Equal(t, uint16(0xee), emu.VRegisters[2])
	})

	t.Run("8xy2: Set VRegister[x] to the value of VRegister[x] AND VRegister[y]", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = uint16(0xcc)
		emu.VRegisters[3] = uint16(0xaa)
		emu.Decode(0x8232)

		assert.Equal(t, uint16(0x88), emu.VRegisters[2])
	})

	t.Run("8xy3: Set Vx to Vx XOR Vy", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = uint16(12)
		emu.VRegisters[3] = uint16(6)
		emu.Decode(0x8233)

		assert.Equal(t, uint16(10), emu.VRegisters[2])
	})

	t.Run("8xy4: Set Vx to Vx + Vy", func(t *testing.T) {
		t.Run("Total is less than 8-bits (255)", func(t *testing.T) {
			emu := NewEmulator()
			emu.VRegisters[2] = uint16(12)
			emu.VRegisters[3] = uint16(6)
			emu.Decode(0x8234)

			assert.Equal(t, uint16(18), emu.VRegisters[2])
			assert.Equal(t, uint16(0), emu.VRegisters[0xF])
		})

		t.Run("Total is greater than 8-bits (255)", func(t *testing.T) {
			emu := NewEmulator()
			emu.VRegisters[2] = uint16(255)
			emu.VRegisters[3] = uint16(1)
			emu.Decode(0x8234)

			assert.Equal(t, uint16(255), emu.VRegisters[2])
			assert.Equal(t, uint16(1), emu.VRegisters[0xF])
		})
	})

	t.Run("8xy5: Set Vx to Vx - Vy", func(t *testing.T) {
		t.Run("Vx is greater than Vy", func(t *testing.T) {
			emu := NewEmulator()
			emu.VRegisters[2] = uint16(12)
			emu.VRegisters[3] = uint16(8)
			emu.Decode(0x8235)

			assert.Equal(t, uint16(4), emu.VRegisters[2])
			assert.Equal(t, uint16(0), emu.VRegisters[0xF])
		})

		t.Run("Vx is less than Vy", func(t *testing.T) {
			emu := NewEmulator()
			emu.VRegisters[2] = uint16(12)
			emu.VRegisters[3] = uint16(14)
			emu.Decode(0x8235)

			assert.Equal(t, uint16(0), emu.VRegisters[2])
			assert.Equal(t, uint16(1), emu.VRegisters[0xF])
		})
	})

	t.Run("8xy6: Right shift Vx by 1 and store that bit in VF", func(t *testing.T) {
		t.Run("The shifted bit is 1", func(t *testing.T) {
			emu := NewEmulator()
			emu.VRegisters[2] = uint16(0x00C1)
			emu.Decode(0x8206)

			assert.Equal(t, uint16(0x60), emu.VRegisters[2])
			assert.Equal(t, uint16(1), emu.VRegisters[0xF])
		})

		t.Run("The shifted bit is 0", func(t *testing.T) {
			emu := NewEmulator()
			emu.VRegisters[2] = uint16(0x00CC)
			emu.Decode(0x8206)

			assert.Equal(t, uint16(0x66), emu.VRegisters[2])
			assert.Equal(t, uint16(0), emu.VRegisters[0xF])
		})
	})

	t.Run("8xy7: Set Vx to Vy - Vx", func(t *testing.T) {
		t.Run("Vy is greater than Vx", func(t *testing.T) {
			emu := NewEmulator()
			emu.VRegisters[2] = uint16(10)
			emu.VRegisters[3] = uint16(12)
			emu.Decode(0x8237)

			assert.Equal(t, uint16(2), emu.VRegisters[2])
			assert.Equal(t, uint16(1), emu.VRegisters[0xF])
		})

		t.Run("Vy is less than Vx", func(t *testing.T) {
			emu := NewEmulator()
			emu.VRegisters[2] = uint16(12)
			emu.VRegisters[3] = uint16(10)
			emu.Decode(0x8237)

			assert.Equal(t, uint16(0), emu.VRegisters[2])
			assert.Equal(t, uint16(0), emu.VRegisters[0xF])
		})
	})

	t.Run("8xyE: Left shift Vx by 1 and store that bit in VF", func(t *testing.T) {
		t.Run("The shifted bit is 1", func(t *testing.T) {
			emu := NewEmulator()
			emu.VRegisters[2] = uint16(0b10000000)
			emu.Decode(0x820E)

			assert.Equal(t, uint16(0b100000000), emu.VRegisters[2])
			assert.Equal(t, uint16(1), emu.VRegisters[0xF])
		})

		t.Run("The shifted bit is 0", func(t *testing.T) {
			emu := NewEmulator()
			emu.VRegisters[2] = uint16(0b01000000)
			emu.Decode(0x820E)

			assert.Equal(t, uint16(0b10000000), emu.VRegisters[2])
			assert.Equal(t, uint16(0), emu.VRegisters[0xF])
		})
	})

	t.Run("Annn: Sets IRegister to nnn", func(t *testing.T) {
		emu := NewEmulator()
		emu.Decode(0xA105)

		got := emu.IRegister
		want := uint16(0x105)

		assert.Equal(t, want, got)
	})

	t.Run("Dxyn: Adds a sprite to the Screen array", func(t *testing.T) {
		emu := NewEmulator()
		emu.Decode(0xD123)
		// TODO: hmmm
	})
}
