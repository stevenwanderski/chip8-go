package cpu

import (
	"os"
)

const (
	RAM_SIZE       uint16 = 4096
	SCREEN_WIDTH   uint16 = 64
	SCREEN_HEIGHT  uint16 = 32
	SCREEN_TOTAL   uint16 = SCREEN_WIDTH * SCREEN_HEIGHT
	SCREEN_SCALE   uint16 = 15
	REGISTER_COUNT uint8  = 16
	STACK_SIZE     uint8  = 16
	START_ADDRESS  uint16 = 512
)

var fontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Emulator struct {
	ProgramCounter uint16
	Ram            [RAM_SIZE]byte
	Screen         [SCREEN_TOTAL]uint8
	VRegisters     [REGISTER_COUNT]uint8
	IRegister      uint16
	Stack          [STACK_SIZE]uint16
	StackPointer   uint16
	DelayTimer     uint16
	SoundTimer     uint16
	Opcode         uint16
	Keys           [16]uint8
}

func (e *Emulator) Tick() {
	opcode := e.Fetch()
	e.Decode(opcode)
}

func (e *Emulator) TickTimers() {
	if e.DelayTimer > 0 {
		e.DelayTimer -= 1
	}

	if e.SoundTimer > 0 {
		// Check if 1, then make a beep
		e.SoundTimer -= 1
	}
}

func (e *Emulator) LoadRom(filepath string) {
	data, err := os.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	for i, v := range data {
		e.Ram[uint16(i)+START_ADDRESS] = v
	}
}

func (e *Emulator) Push(value uint16) {
	e.Stack[e.StackPointer] = value
	e.StackPointer += 1
}

func (e *Emulator) Pop() uint16 {
	e.StackPointer -= 1
	return e.Stack[e.StackPointer]
}

func (e *Emulator) Key(value uint8, pressed uint8) {
	e.Keys[value] = pressed
}

func (e *Emulator) Fetch() uint16 {
	first_code := e.Ram[e.ProgramCounter]
	second_code := e.Ram[e.ProgramCounter+1]
	e.Opcode = (uint16(first_code) << 8) | uint16(second_code)

	if e.ProgramCounter < 4094 {
		e.ProgramCounter += 2
	} else {
		e.ProgramCounter = START_ADDRESS
	}

	return e.Opcode
}

func (e *Emulator) Decode(opcode uint16) {
	decoder := Decoder{emu: e}
	decoder.Run(opcode)
}

func NewEmulator() Emulator {
	emu := Emulator{}
	emu.ProgramCounter = START_ADDRESS
	emu.DelayTimer = 0
	emu.SoundTimer = 0

	for i, v := range fontSet {
		emu.Ram[i] = v
	}

	return emu
}
