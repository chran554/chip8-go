package chip8

import "fmt"

type Screen struct {
	Width  uint8
	Height uint8
	gfxMem [64][32]uint8
}

func NewScreen() Screen {
	return Screen{
		Width:  64,
		Height: 32,
		gfxMem: [64][32]uint8{},
	}
}

func (s *Screen) XorPixel(x, y uint8, v uint8) byte {
	s.gfxMem[x][y] ^= v
	return s.gfxMem[x][y]
}

func (s *Screen) Value(x, y uint8) uint8 {
	return s.gfxMem[x][y]
}

func (s *Screen) Clear() {
	for _, uint8s := range s.gfxMem {
		for i := range uint8s {
			uint8s[i] = 0
		}
	}
}

func (s *Screen) Print() {
	fmt.Println("================================")
	for y := uint8(0); y < s.Height; y++ {
		for x := uint8(0); x < s.Width; x++ {
			pixelValue := s.Value(x, y)
			if pixelValue == 0 {
				fmt.Print("░")
			} else {
				fmt.Print("█")
			}
		}
		fmt.Println()
	}
	fmt.Println("================================")
}
