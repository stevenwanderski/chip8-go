package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	RAM_SIZE       uint16 = 4096
	SCREEN_WIDTH   uint16 = 64
	SCREEN_HEIGHT  uint16 = 32
	SCREEN_SCALE   uint16 = 10
	REGISTER_COUNT uint8  = 16
	STACK_SIZE     uint8  = 16
	KEY_COUNT      uint8  = 16
	START_ADDRESS  uint16 = 512
)

type Emulator struct {
	ProgramCounter uint16
	Ram            [RAM_SIZE]byte
	Screen         [SCREEN_WIDTH * SCREEN_HEIGHT]bool
	VRegisters     [REGISTER_COUNT]uint16
	IRegister      uint16
	Stack          [STACK_SIZE]uint16
	StackPointer   uint16
	DelayTimer     uint16
	SoundTimer     uint16
	Opcode         uint16
}

type Display struct {
}

func (d *Display) Run(emulator Emulator) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"CHIP-8",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(SCREEN_WIDTH*SCREEN_SCALE),
		int32(SCREEN_HEIGHT*SCREEN_SCALE),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	renderer.SetDrawColor(186, 177, 144, 255)
	renderer.Clear()

	running := true
	for running {
		emulator.Cycle()

		renderer.SetDrawColor(186, 177, 144, 255)
		renderer.Clear()

		for i, v := range emulator.Screen {
			rect := sdl.Rect{
				X: (int32(i) % int32(SCREEN_WIDTH)) * int32(SCREEN_SCALE),
				Y: (int32(i) / int32(SCREEN_WIDTH)) * int32(SCREEN_SCALE),
				W: int32(SCREEN_SCALE),
				H: int32(SCREEN_SCALE),
			}

			if v == true {
				renderer.SetDrawColor(108, 149, 117, 255)
			} else {
				renderer.SetDrawColor(186, 177, 144, 255)
			}

			renderer.FillRect(&rect)
		}

		renderer.Present()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			}
		}

		sdl.Delay(1000 / 60)
	}
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
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			for i := range e.Screen {
				e.Screen[i] = false
			}

			fmt.Printf("Opcode %x: Clear the screen\n", opcode)
			break

		case 0x00EE:
			e.ProgramCounter = e.Pop()

			fmt.Printf("Opcode %x: Set the ProgramCounter to stack address\n", opcode)
			break
		}

	case 0x1000:
		nnn := opcode & 0x0FFF
		e.ProgramCounter = nnn

		fmt.Printf("Opcode %x: Set ProgramCounter to %d\n", opcode, nnn)
		break

	case 0x2000:
		oldValue := e.ProgramCounter
		e.Push(oldValue)

		nnn := opcode & 0x0FFF
		e.ProgramCounter = nnn

		fmt.Printf("Opcode %x: Set ProgramCounter to %d and add %d to the stack \n", opcode, nnn, oldValue)
		break

	case 0x6000:
		x := (opcode & 0x0F00) >> 8
		nn := opcode & 0x00FF
		e.VRegisters[x] = nn

		fmt.Printf("Opcode %x: Set VRegister %d to %d\n", opcode, x, nn)
		break

	case 0x7000:
		x := (opcode & 0x0F00) >> 8
		nn := opcode & 0x00FF
		e.VRegisters[x] += nn

		fmt.Printf("Opcode %x: Add %d to VRegister %d\n", opcode, nn, x)
		break

	case 0xA000:
		nnn := opcode & 0x0FFF
		e.IRegister = nnn

		fmt.Printf("Opcode %x: Set IRegister to %d\n", opcode, nnn)
		break

	case 0xD000:
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		n := opcode & 0x000F

		start_addr := e.IRegister

		var i uint16 = 0
		var j uint16 = 0
		for i = 0; i < n; i++ {
			pixels := e.Ram[start_addr+i]

			for j = 0; j < 8; j++ {
				if pixels&(0b10000000>>j) != 0 {
					x_position := (e.VRegisters[x] + j) % SCREEN_WIDTH
					y_position := (e.VRegisters[y] + i) % SCREEN_HEIGHT

					screen_index := (y_position * SCREEN_WIDTH) + x_position
					e.Screen[screen_index] = true
				}
			}
		}

		fmt.Printf("Opcode %x: Draw %d rows high at X: %d, Y: %d\n", opcode, n, e.VRegisters[x], e.VRegisters[y])
		break
	}
}

func NewEmulator() Emulator {
	emu := Emulator{}
	emu.ProgramCounter = START_ADDRESS
	return emu
}

func main() {
	emu := NewEmulator()
	emu.LoadRom("./roms/ibm-logo.ch8")

	display := Display{}
	display.Run(emu)
}
