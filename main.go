package main

import (
	"chip-8/cpu"
	"os"
)

func main() {
	rom_path := os.Args[1]
	emu := cpu.NewEmulator()
	emu.LoadRom(rom_path)

	display := cpu.Display{}
	display.Run(emu)
}
