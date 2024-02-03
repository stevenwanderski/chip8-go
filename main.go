package main

import (
	"chip-8/cpu"
)

func main() {
	emu := cpu.NewEmulator()
	emu.LoadRom("./roms/test-opcode.ch8")
	// emu.LoadRom("./roms/ibm-logo.ch8")

	display := cpu.Display{}
	display.Run(emu)
}
