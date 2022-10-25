package main

import (
	"chip8/pkg/chip8"
	"fmt"
)

func main() {
	fmt.Println("Up 'n' running...")

	//romFilepath := "roms/test_opcode.ch8"
	//romFilepath := "roms/IBM Logo.ch8"
	//romFilepath := "roms/PONG.ch8"
	romFilepath := "roms/BRIX.ch8"

	screen := chip8.NewScreen()
	chip8 := chip8.NewChip8(screen)
	chip8.LoadROM(romFilepath)
	chip8.Run(false)

	//fmt.Printf("Disassembly of \"%s\":\n", romFilepath)
	//chip8.DisassembleProgram(romFilepath, 0x200)

	fmt.Println("Done.")
}
