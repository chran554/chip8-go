package chip8

import (
	"bytes"
	"fmt"
)

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
	if (x >= s.Width) || (y >= s.Height) {
		return 0x00
	} else {
		s.gfxMem[x][y] ^= v
		return s.gfxMem[x][y]
	}
}

func (s *Screen) Value(x, y uint8) uint8 {
	if (x >= s.Width) || (y >= s.Height) {
		return 0x00
	} else {
		return s.gfxMem[x][y]
	}
}

func (s *Screen) Clear() {
	for _, uint8s := range s.gfxMem {
		for i := range uint8s {
			uint8s[i] = 0x00
		}
	}
}

func (s *Screen) Print() {
	var buffer bytes.Buffer
	buffer.WriteString("\n")
	for y := uint8(0); y < s.Height; y++ {
		for x := uint8(0); x < s.Width; x++ {
			pixelValue := s.Value(x, y)
			if pixelValue == 0 {
				buffer.WriteString("░░")
			} else {
				buffer.WriteString("██")
			}
		}
		buffer.WriteString("\n")
	}

	fmt.Println(buffer.String())
}
