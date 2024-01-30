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
