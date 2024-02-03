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
