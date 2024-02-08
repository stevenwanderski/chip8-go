package cpu

import (
	"github.com/veandco/go-sdl2/sdl"
)

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

	setBackgroundColor(renderer)
	renderer.Clear()

	running := true
	for running {
		for i := 0; i < 10; i++ {
			emulator.Tick()
		}

		emulator.TickTimers()
		d.DrawScreen(renderer, emulator.Screen)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
				break

			case *sdl.KeyboardEvent:
				switch t.State {
				case sdl.RELEASED:
					switch t.Keysym.Sym {
					case sdl.K_1:
						emulator.Key(0x1, 0)
					case sdl.K_2:
						emulator.Key(0x2, 0)
					case sdl.K_3:
						emulator.Key(0x3, 0)
					case sdl.K_4:
						emulator.Key(0xC, 0)
					case sdl.K_q:
						emulator.Key(0x4, 0)
					case sdl.K_w:
						emulator.Key(0x5, 0)
					case sdl.K_e:
						emulator.Key(0x6, 0)
					case sdl.K_r:
						emulator.Key(0xD, 0)
					case sdl.K_a:
						emulator.Key(0x7, 0)
					case sdl.K_s:
						emulator.Key(0x8, 0)
					case sdl.K_d:
						emulator.Key(0x9, 0)
					case sdl.K_f:
						emulator.Key(0xE, 0)
					case sdl.K_z:
						emulator.Key(0xA, 0)
					case sdl.K_x:
						emulator.Key(0xB, 0)
					case sdl.K_c:
						emulator.Key(0x0, 0)
					case sdl.K_v:
						emulator.Key(0xf, 0)
					}

				case sdl.PRESSED:
					switch t.Keysym.Sym {
					case sdl.K_1:
						emulator.Key(0x1, 1)
					case sdl.K_2:
						emulator.Key(0x2, 1)
					case sdl.K_3:
						emulator.Key(0x3, 1)
					case sdl.K_4:
						emulator.Key(0xC, 1)
					case sdl.K_q:
						emulator.Key(0x4, 1)
					case sdl.K_w:
						emulator.Key(0x5, 1)
					case sdl.K_e:
						emulator.Key(0x6, 1)
					case sdl.K_r:
						emulator.Key(0xD, 1)
					case sdl.K_a:
						emulator.Key(0x7, 1)
					case sdl.K_s:
						emulator.Key(0x8, 1)
					case sdl.K_d:
						emulator.Key(0x9, 1)
					case sdl.K_f:
						emulator.Key(0xE, 1)
					case sdl.K_z:
						emulator.Key(0xA, 1)
					case sdl.K_x:
						emulator.Key(0xB, 1)
					case sdl.K_c:
						emulator.Key(0x0, 1)
					case sdl.K_v:
						emulator.Key(0xf, 1)
					}
				}
			}
		}

		sdl.Delay(1000 / 60)
	}
}

func (d *Display) DrawScreen(renderer *sdl.Renderer, screen [SCREEN_TOTAL]uint8) {
	setBackgroundColor(renderer)
	renderer.Clear()

	for i, v := range screen {
		rect := sdl.Rect{
			X: (int32(i) % int32(SCREEN_WIDTH)) * int32(SCREEN_SCALE),
			Y: (int32(i) / int32(SCREEN_WIDTH)) * int32(SCREEN_SCALE),
			W: int32(SCREEN_SCALE),
			H: int32(SCREEN_SCALE),
		}

		if v == 1 {
			setForegroundColor(renderer)
		} else {
			setBackgroundColor(renderer)
		}

		renderer.FillRect(&rect)
	}

	renderer.Present()
}

func setBackgroundColor(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 0)
}

func setForegroundColor(renderer *sdl.Renderer) {
	renderer.SetDrawColor(15, 255, 80, 255)
}
