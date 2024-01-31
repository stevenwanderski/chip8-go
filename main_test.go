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
