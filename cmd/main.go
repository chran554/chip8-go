package main

import (
	"chip8/pkg/chip8"
	"fmt"
)

func main() {
	fmt.Println("Up 'n' running...")

	screen := chip8.NewScreen()
	chip8 := chip8.NewChip8(screen)
	chip8.LoadROM("roms/IBM Logo.ch8")
	chip8.Run()

	fmt.Println("Done.")
}
