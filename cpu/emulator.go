package cpu

import (
	"os"
)

const (
	RAM_SIZE       uint16 = 4096
	SCREEN_WIDTH   uint16 = 64
	SCREEN_HEIGHT  uint16 = 32
	SCREEN_SCALE   uint16 = 10
	REGISTER_COUNT uint8  = 16
	STACK_SIZE     uint8  = 16
	START_ADDRESS  uint16 = 512
)

type Emulator struct {
	ProgramCounter uint16
	Ram            [RAM_SIZE]byte
	Screen         [SCREEN_WIDTH * SCREEN_HEIGHT]bool
	VRegisters     [REGISTER_COUNT]uint8
	IRegister      uint16
	Stack          [STACK_SIZE]uint16
	StackPointer   uint16
	DelayTimer     uint16
	SoundTimer     uint16
	Opcode         uint16
}

func (e *Emulator) Cycle() {
	opcode := e.Fetch()
	e.Decode(opcode)
}

func (e *Emulator) LoadRom(filepath string) {
	data, _ := os.ReadFile(filepath)

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
	return emu
}
