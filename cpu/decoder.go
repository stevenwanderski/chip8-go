package cpu

import (
	"fmt"
	"math/rand"
)

type Decoder struct {
	emu *Emulator
}

func (d *Decoder) Run(opcode uint16) {
	e := d.emu

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

		// fmt.Printf("Opcode %x: Set ProgramCounter to %d\n", opcode, nnn)
		break

	case 0x2000:
		oldValue := e.ProgramCounter
		e.Push(oldValue)

		nnn := opcode & 0x0FFF
		e.ProgramCounter = nnn

		fmt.Printf("Opcode %x: Set ProgramCounter to %d and add %d to the stack \n", opcode, nnn, oldValue)
		break

	case 0x3000:
		x := (opcode & 0x0F00) >> 8
		nn := uint8(opcode & 0x00FF)

		if e.VRegisters[x] == nn {
			e.ProgramCounter += 2
			// fmt.Printf("Opcode %x: Set ProgramCounter to %d\n", opcode, e.ProgramCounter)
		} else {
			// fmt.Printf("Opcode %x: Skip because VRegister[%d] (%d) does not equal nn (%d)\n", opcode, x, e.VRegisters[x], nn)
		}
		break

	case 0x4000:
		x := (opcode & 0x0F00) >> 8
		nn := uint8(opcode & 0x00FF)

		if e.VRegisters[x] != nn {
			e.ProgramCounter += 2
			fmt.Printf("Opcode %x: Set ProgramCounter to %d\n", opcode, e.ProgramCounter)
		} else {
			fmt.Printf("Opcode %x: Skip because VRegister[%d] (%d) equals nn (%d)\n", opcode, x, e.VRegisters[x], nn)
		}
		break

	case 0x5000:
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4

		if e.VRegisters[x] == e.VRegisters[y] {
			e.ProgramCounter += 2
			fmt.Printf("Opcode %x: Set ProgramCounter to %d\n", opcode, e.ProgramCounter)
		} else {
			fmt.Printf("Opcode %x: Skip because VRegister[%d] (%d) does not equal VRegister[%d] (%d)\n", opcode, x, e.VRegisters[x], y, e.VRegisters[y])
		}
		break

	case 0x6000:
		x := (opcode & 0x0F00) >> 8
		nn := opcode & 0x00FF
		e.VRegisters[x] = uint8(nn)

		fmt.Printf("Opcode %x: Set VRegister %d to %d\n", opcode, x, nn)
		break

	case 0x7000:
		x := (opcode & 0x0F00) >> 8
		nn := opcode & 0x00FF
		e.VRegisters[x] += uint8(nn)

		fmt.Printf("Opcode %x: Add %d to VRegister[%d] (%d)\n", opcode, nn, x, e.VRegisters[x])
		break

	case 0x8000:
		switch opcode & 0x000F {
		case 0:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			e.VRegisters[x] = e.VRegisters[y]

			fmt.Printf("Opcode %x: Set VRegister[%d] to VRegister[%d] (%d)\n", opcode, x, y, e.VRegisters[y])
			break

		case 1:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			e.VRegisters[x] = e.VRegisters[x] | e.VRegisters[y]

			fmt.Printf("Opcode %x: Set VRegister[%d] to VRegister[%d] (%d) OR VRegister[%d] (%d)\n", opcode, x, x, e.VRegisters[x], y, e.VRegisters[y])
			break

		case 2:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			e.VRegisters[x] = e.VRegisters[x] & e.VRegisters[y]

			fmt.Printf("Opcode %x: Set VRegister[%d] to VRegister[%d] (%d) AND VRegister[%d] (%d)\n", opcode, x, x, e.VRegisters[x], y, e.VRegisters[y])
			break

		case 3:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			e.VRegisters[x] = e.VRegisters[x] ^ e.VRegisters[y]

			fmt.Printf("Opcode %x: Set VRegister[%d] to VRegister[%d] (%d) XOR VRegister[%d] (%d)\n", opcode, x, x, e.VRegisters[x], y, e.VRegisters[y])
			break

		case 4:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4

			if e.VRegisters[y] > (255 - e.VRegisters[x]) {
				e.VRegisters[0xF] = uint8(1)
			} else {
				e.VRegisters[0xF] = uint8(0)
			}

			e.VRegisters[x] = e.VRegisters[x] + e.VRegisters[y]

			fmt.Printf("Opcode %x: Set VRegister[%d] to VRegister[%d] (%d) XOR VRegister[%d] (%d)\n", opcode, x, x, e.VRegisters[x], y, e.VRegisters[y])
			break

		case 5:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4

			if e.VRegisters[x] > e.VRegisters[y] {
				e.VRegisters[0xF] = uint8(1)
			} else {
				e.VRegisters[0xF] = uint8(0)
			}

			e.VRegisters[x] = e.VRegisters[x] - e.VRegisters[y]

			fmt.Printf("Opcode %x: Set VRegister[%d] to VRegister[%d] (%d) XOR VRegister[%d] (%d)\n", opcode, x, x, e.VRegisters[x], y, e.VRegisters[y])
			break

		case 6:
			x := (opcode & 0x0F00) >> 8
			v := e.VRegisters[x]

			e.VRegisters[0xF] = v & 1
			e.VRegisters[x] = v >> 1

			// fmt.Printf("Opcode %x: Set VRegister[%d] to VRegister[%d] (%d) XOR VRegister[%d] (%d)\n", opcode, x, x, e.VRegisters[x], y, e.VRegisters[y])
			break

		case 7:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4

			if e.VRegisters[y] > e.VRegisters[x] {
				e.VRegisters[x] = e.VRegisters[y] - e.VRegisters[x]
				e.VRegisters[0xF] = 1
			} else {
				e.VRegisters[x] = 0
				e.VRegisters[0xF] = 0
			}

			// fmt.Printf("Opcode %x: Set VRegister[%d] to VRegister[%d] (%d) XOR VRegister[%d] (%d)\n", opcode, x, x, e.VRegisters[x], y, e.VRegisters[y])
			break

		case 0xE:
			x := (opcode & 0x0F00) >> 8
			v := e.VRegisters[x]

			e.VRegisters[0xF] = v >> 7 & 1
			e.VRegisters[x] = v << 1

			// fmt.Printf("Opcode %x: Set VRegister[%d] to VRegister[%d] (%d) XOR VRegister[%d] (%d)\n", opcode, x, x, e.VRegisters[x], y, e.VRegisters[y])
			break

		}

		break

	case 0x9000:
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4

		if e.VRegisters[x] != e.VRegisters[y] {
			e.ProgramCounter += 2
		}

		// fmt.Printf("Opcode %x: Set IRegister to %d\n", opcode, nnn)
		break

	case 0xA000:
		nnn := opcode & 0x0FFF
		e.IRegister = nnn

		fmt.Printf("Opcode %x: Set IRegister to %d\n", opcode, nnn)
		break

	case 0xB000:
		nnn := uint8(opcode & 0x0FFF)
		v := e.VRegisters[0]
		e.ProgramCounter += uint16(v + nnn)

		// fmt.Printf("Opcode %x: Set IRegister to %d\n", opcode, nnn)
		break

	case 0xC000:
		x := (opcode & 0x0F00) >> 8
		nn := opcode & 0x00FF
		random_number := rand.Intn(256)

		e.VRegisters[x] = uint8(random_number) & uint8(nn)
		break

	case 0xD000:
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		n := uint8(opcode & 0x000F)

		start_addr := e.IRegister

		var i uint8 = 0
		var j uint8 = 0

		// For each row (n)
		for i = 0; i < n; i++ {
			// Get the value from RAM
			pixels := e.Ram[start_addr+uint16(i)]

			// For each bit (0 or 1) in the RAM value
			for j = 0; j < 8; j++ {
				// If the bit equals 1
				if pixels&(0b10000000>>j) != 0 {
					// Get value from VRegister[x] and add the bit position (j).
					// Then modulo it by the width of the screen to simulate wrapping.
					// This way x_position will never be greater than the screen width.
					x_position := uint16(e.VRegisters[x] + j)

					// Get the "Y" screen position from the y VRegister
					y_position := uint16(e.VRegisters[y] + i)

					screen_index := (y_position * SCREEN_WIDTH) + x_position
					e.Screen[screen_index] = true
				}
			}
		}

		fmt.Printf("Opcode %x: Draw %d rows high at X: %d, Y: %d\n", opcode, n, e.VRegisters[x], e.VRegisters[y])
		break

	case 0xE000:
		switch opcode & 0x00FF {
		case 0x9E:
			x := (opcode & 0x0F00) >> 8

			if e.Keys[e.VRegisters[x]] == 1 {
				e.ProgramCounter += 2
			}
			break

		case 0xA1:
			x := (opcode & 0x0F00) >> 8

			if e.Keys[e.VRegisters[x]] == 0 {
				e.ProgramCounter += 2
			}

			break
		}

	case 0xF000:
		x := (opcode & 0x0F00) >> 8

		switch opcode & 0x00FF {
		case 0x07:
			e.VRegisters[x] = uint8(e.DelayTimer)

			break

		case 0x0A:
			pressed := false

			for i, v := range e.Keys {
				if v == 1 {
					pressed = true
					e.VRegisters[x] = uint8(i)
				}
			}

			if !pressed {
				e.ProgramCounter -= 2
			}

			break

		case 0x15:
			e.DelayTimer = uint16(e.VRegisters[x])

			break

		case 0x18:
			e.SoundTimer = uint16(e.VRegisters[x])

			break

		case 0x1E:
			e.IRegister += uint16(e.VRegisters[x])

			break

		case 0x29:
			e.IRegister = uint16(e.VRegisters[x] * 5)

			break

		case 0x33:

			break

		case 0x55:

			break

		case 0x65:

			break

		}
		break
	}
}
