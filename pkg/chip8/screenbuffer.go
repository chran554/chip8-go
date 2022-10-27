package chip8

import (
	"bytes"
	"fmt"
)

type ScreenBuffer struct {
	Width  uint8
	Height uint8
	buffer [64][32]uint8
}

func NewScreenBuffer() ScreenBuffer {
	return ScreenBuffer{
		Width:  64,
		Height: 32,
		buffer: [64][32]uint8{},
	}
}

func (s *ScreenBuffer) XorPixel(x, y uint8, v uint8) byte {
	if (x >= s.Width) || (y >= s.Height) {
		return 0x00
	} else {
		s.buffer[x][y] ^= v
		return s.buffer[x][y]
	}
}

func (s *ScreenBuffer) Value(x, y uint8) uint8 {
	if (x >= s.Width) || (y >= s.Height) {
		return 0x00
	} else {
		return s.buffer[x][y]
	}
}

func (s *ScreenBuffer) Clear() {
	for _, uint8s := range s.buffer {
		for i := range uint8s {
			uint8s[i] = 0x00
		}
	}
}

func (s *ScreenBuffer) Print() {
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
