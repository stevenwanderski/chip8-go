package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpcode00E0(t *testing.T) {
	emu := NewEmulator()
	emu.Screen[0] = 1
	emu.Screen[1] = 1

	emu.Decode(0x00E0)

	got := uint8(0)
	want := uint8(0)

	for _, v := range emu.Screen {
		if v == 1 {
			got = 1
		}
	}

	assert.Equal(t, want, got)
}

func TestOpcode00EE(t *testing.T) {
	emu := NewEmulator()
	emu.ProgramCounter = uint16(0x0001)
	emu.Stack[0] = uint16(0x0010)
	emu.Stack[1] = uint16(0x0020)
	emu.StackPointer = uint16(2)
	emu.Decode(0x00EE)

	assert.Equal(t, uint16(0x0020), emu.ProgramCounter)
	assert.Equal(t, uint16(1), emu.StackPointer)
}

func TestOpcode1nnn(t *testing.T) {
	emu := NewEmulator()
	emu.Decode(0x1228)

	got := emu.ProgramCounter
	want := uint16(0x228)
	assert.Equal(t, want, got)
}

func TestOpcode2nnn(t *testing.T) {
	emu := NewEmulator()
	emu.ProgramCounter = uint16(0x900)
	emu.Decode(0x2228)

	assert.Equal(t, uint16(0x228), emu.ProgramCounter)
	assert.Equal(t, uint16(0x900), emu.Stack[0])
	assert.Equal(t, uint16(1), emu.StackPointer)
}

func TestOpcode3xnn(t *testing.T) {
	t.Run("VRegister[x] equals nn", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 0x900
		emu.VRegisters[2] = 0xFF
		emu.Decode(0x32FF)

		assert.Equal(t, uint16(0x902), emu.ProgramCounter)
	})

	t.Run("VRegister[x] does not equal nn", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 0x900
		emu.VRegisters[2] = 0xAA
		emu.Decode(0x32FF)

		assert.Equal(t, uint16(0x900), emu.ProgramCounter)
	})
}

func TestOpcode4xnn(t *testing.T) {
	t.Run("VRegister[x] does not equal nn", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 0x900
		emu.VRegisters[2] = 0xFF
		emu.Decode(0x42AA)

		assert.Equal(t, uint16(0x902), emu.ProgramCounter)
	})

	t.Run("VRegister[x] equals nn", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 0x900
		emu.VRegisters[2] = 0xFF
		emu.Decode(0x42FF)

		assert.Equal(t, uint16(0x900), emu.ProgramCounter)
	})
}

func TestOpcode5xyn(t *testing.T) {
	t.Run("VRegister[x] equals VRegister[y]", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 0x900
		emu.VRegisters[2] = 0xFF
		emu.VRegisters[3] = 0xFF
		emu.Decode(0x5230)

		assert.Equal(t, uint16(0x902), emu.ProgramCounter)
	})

	t.Run("VRegister[x] does not equal VRegister[y]", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 0x900
		emu.VRegisters[2] = 0xFF
		emu.VRegisters[3] = 0xAA
		emu.Decode(0x5230)

		assert.Equal(t, uint16(0x900), emu.ProgramCounter)
	})
}

func TestOpcode6xnn(t *testing.T) {
	emu := NewEmulator()
	emu.Decode(0x6105)

	got := emu.VRegisters[1]
	want := uint8(0x05)

	assert.Equal(t, want, got)
}

func TestOpcode7xnn(t *testing.T) {
	emu := NewEmulator()
	emu.VRegisters[1] = 10
	emu.Decode(0x7105)

	got := emu.VRegisters[1]
	want := uint8(15)

	assert.Equal(t, want, got)
}

func TestOpcode8xy0(t *testing.T) {
	emu := NewEmulator()
	emu.VRegisters[2] = 0x1
	emu.VRegisters[3] = 0x2
	emu.Decode(0x8230)

	assert.Equal(t, uint8(0x2), emu.VRegisters[2])
}

func TestOpcode8xy1(t *testing.T) {
	emu := NewEmulator()
	emu.VRegisters[2] = 0xcc
	emu.VRegisters[3] = 0xaa
	emu.Decode(0x8231)

	assert.Equal(t, uint8(0xee), emu.VRegisters[2])
}

func TestOpcode8xy2(t *testing.T) {
	emu := NewEmulator()
	emu.VRegisters[2] = 0xcc
	emu.VRegisters[3] = 0xaa
	emu.Decode(0x8232)

	assert.Equal(t, uint8(0x88), emu.VRegisters[2])
}

func TestOpcode8xy3(t *testing.T) {
	emu := NewEmulator()
	emu.VRegisters[2] = 12
	emu.VRegisters[3] = 6
	emu.Decode(0x8233)

	assert.Equal(t, uint8(10), emu.VRegisters[2])
}

func TestOpcode8xy4(t *testing.T) {
	t.Run("Total is less than 8-bits (255)", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = 12
		emu.VRegisters[3] = 6
		emu.Decode(0x8234)

		assert.Equal(t, uint8(18), emu.VRegisters[2])
		assert.Equal(t, uint8(0), emu.VRegisters[0xF])
	})

	t.Run("Total is greater than 8-bits (255)", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = 255
		emu.VRegisters[3] = 3
		emu.Decode(0x8234)

		assert.Equal(t, uint8(2), emu.VRegisters[2])
		assert.Equal(t, uint8(1), emu.VRegisters[0xF])
	})
}

func TestOpcode8xy5(t *testing.T) {
	t.Run("Vx is greater than Vy", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = 12
		emu.VRegisters[3] = 8
		emu.Decode(0x8235)

		assert.Equal(t, uint8(4), emu.VRegisters[2])
		assert.Equal(t, uint8(1), emu.VRegisters[0xF])
	})

	t.Run("Vx is less than Vy", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = 12
		emu.VRegisters[3] = 14
		emu.Decode(0x8235)

		// When subtraction causes a negative overflow
		// 255 comes after 0
		assert.Equal(t, uint8(254), emu.VRegisters[2])
		assert.Equal(t, uint8(0), emu.VRegisters[0xF])
	})
}

func TestOpcode8xy6(t *testing.T) {
	t.Run("The shifted bit is 1", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = 0x00C1
		emu.Decode(0x8206)

		assert.Equal(t, uint8(0x60), emu.VRegisters[2])
		assert.Equal(t, uint8(1), emu.VRegisters[0xF])
	})

	t.Run("The shifted bit is 0", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = 0x00CC
		emu.Decode(0x8206)

		assert.Equal(t, uint8(0x66), emu.VRegisters[2])
		assert.Equal(t, uint8(0), emu.VRegisters[0xF])
	})
}

func TestOpcode8xy7(t *testing.T) {
	t.Run("Vy is greater than Vx", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = 10
		emu.VRegisters[3] = 12
		emu.Decode(0x8237)

		assert.Equal(t, uint8(2), emu.VRegisters[2])
		assert.Equal(t, uint8(1), emu.VRegisters[0xF])
	})

	t.Run("Vy is less than Vx", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = 12
		emu.VRegisters[3] = 10
		emu.Decode(0x8237)

		assert.Equal(t, uint8(0), emu.VRegisters[2])
		assert.Equal(t, uint8(0), emu.VRegisters[0xF])
	})
}

func TestOpcode8xyE(t *testing.T) {
	t.Run("The shifted bit is 1", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = 0x80
		emu.Decode(0x820E)

		assert.Equal(t, uint8(0x0), emu.VRegisters[2])
		assert.Equal(t, uint8(1), emu.VRegisters[0xF])
	})

	t.Run("The shifted bit is 0", func(t *testing.T) {
		emu := NewEmulator()
		emu.VRegisters[2] = 0x3
		emu.Decode(0x820E)

		assert.Equal(t, uint8(0x6), emu.VRegisters[2])
		assert.Equal(t, uint8(0), emu.VRegisters[0xF])
	})
}

func TestOpcode9xy0(t *testing.T) {
	t.Run("Vx != Vy", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 2
		emu.VRegisters[2] = 0xcc
		emu.VRegisters[3] = 0xaa
		emu.Decode(0x9230)

		assert.Equal(t, uint16(4), emu.ProgramCounter)
	})

	t.Run("Vx == Vy", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 2
		emu.VRegisters[2] = 0xcc
		emu.VRegisters[3] = 0xcc
		emu.Decode(0x9230)

		assert.Equal(t, uint16(2), emu.ProgramCounter)
	})
}

func TestOpcodeAnnn(t *testing.T) {
	emu := NewEmulator()
	emu.Decode(0xA105)

	got := emu.IRegister
	want := uint16(0x105)

	assert.Equal(t, want, got)
}

func TestOpcodeBnnn(t *testing.T) {
	emu := NewEmulator()
	emu.ProgramCounter = 2
	emu.VRegisters[0] = 4
	emu.Decode(0xB003)

	assert.Equal(t, uint16(9), emu.ProgramCounter)
}

func TestOpcodeCxnn(t *testing.T) {
	// TODO: Find a way to mock rand.Intn
}

func TestOpcodeDxyn(t *testing.T) {
	emu := NewEmulator()
	emu.Decode(0xD123)
	// TODO: hmmm
}

func TestOpcodeEx9E(t *testing.T) {
	t.Run("Key is pressed", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 900
		emu.VRegisters[1] = 0x3
		emu.Keys[0x3] = 1
		emu.Decode(0xE19E)

		assert.Equal(t, uint16(902), emu.ProgramCounter)
	})

	t.Run("Key is not pressed", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 900
		emu.VRegisters[1] = 0x3
		emu.Keys[0x3] = 0
		emu.Decode(0xE19E)

		assert.Equal(t, uint16(900), emu.ProgramCounter)
	})
}

func TestOpcodeEx1A(t *testing.T) {
	t.Run("Key is not pressed", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 900
		emu.VRegisters[1] = 0x3
		emu.Keys[0x3] = 0
		emu.Decode(0xE1A1)

		assert.Equal(t, uint16(902), emu.ProgramCounter)
	})

	t.Run("Key is pressed", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 900
		emu.VRegisters[1] = 0x3
		emu.Keys[0x3] = 1
		emu.Decode(0xE1A1)

		assert.Equal(t, uint16(900), emu.ProgramCounter)
	})
}

func TestOpcodeFx07(t *testing.T) {
	emu := NewEmulator()
	emu.DelayTimer = 10
	emu.VRegisters[2] = 0
	emu.Decode(0xF207)

	assert.Equal(t, uint8(10), emu.VRegisters[2])
}

func TestOpcodeFx0A(t *testing.T) {
	t.Run("no key is pressed", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 10
		emu.Decode(0xF20A)

		assert.Equal(t, uint16(8), emu.ProgramCounter)
	})

	t.Run("a key is pressed", func(t *testing.T) {
		emu := NewEmulator()
		emu.ProgramCounter = 10
		emu.Keys[0xC] = 1
		emu.Decode(0xF20A)

		assert.Equal(t, uint16(10), emu.ProgramCounter)
		assert.Equal(t, uint8(0xC), emu.VRegisters[2])
	})
}

func TestOpcodeFx15(t *testing.T) {
	emu := NewEmulator()
	emu.DelayTimer = 0
	emu.VRegisters[2] = 10
	emu.Decode(0xF215)

	assert.Equal(t, uint16(10), emu.DelayTimer)
}

func TestOpcodeFx18(t *testing.T) {
	emu := NewEmulator()
	emu.SoundTimer = 0
	emu.VRegisters[2] = 10
	emu.Decode(0xF218)

	assert.Equal(t, uint16(10), emu.SoundTimer)
}

func TestOpcodeFx1E(t *testing.T) {
	emu := NewEmulator()
	emu.IRegister = 5
	emu.VRegisters[2] = 10
	emu.Decode(0xF21E)

	assert.Equal(t, uint16(15), emu.IRegister)
}

func TestOpcodeFx29(t *testing.T) {
	emu := NewEmulator()
	emu.IRegister = 0
	emu.VRegisters[2] = 0x4
	emu.Decode(0xF229)

	assert.Equal(t, uint16(20), emu.IRegister)
}

func TestOpcodeFx33(t *testing.T) {
	emu := NewEmulator()
	emu.IRegister = 999
	emu.VRegisters[2] = 216
	emu.Decode(0xF233)

	assert.Equal(t, uint8(2), emu.Ram[999])
	assert.Equal(t, uint8(1), emu.Ram[1000])
	assert.Equal(t, uint8(6), emu.Ram[1001])
}

func TestOpcodeFx55(t *testing.T) {
	emu := NewEmulator()
	emu.IRegister = 999
	emu.VRegisters[0] = 100
	emu.VRegisters[1] = 200
	emu.VRegisters[2] = 255
	emu.Decode(0xF255)

	assert.Equal(t, uint8(100), emu.Ram[999])
	assert.Equal(t, uint8(200), emu.Ram[1000])
	assert.Equal(t, uint8(255), emu.Ram[1001])
}

func TestOpcodeFx65(t *testing.T) {
	emu := NewEmulator()
	emu.IRegister = 999
	emu.Ram[999] = 100
	emu.Ram[1000] = 200
	emu.Ram[1001] = 255
	emu.Decode(0xF265)

	assert.Equal(t, uint8(100), emu.VRegisters[0])
	assert.Equal(t, uint8(200), emu.VRegisters[1])
	assert.Equal(t, uint8(255), emu.VRegisters[2])
}
