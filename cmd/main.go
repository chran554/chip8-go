package main

import (
	"chip8/pkg/chip8"
	"fmt"
)

func main() {
	fmt.Println("Up 'n' running...")

	//screen := chip8.NewScreen()
	//chip8 := chip8.NewChip8(screen)
	//chip8.LoadROM("roms/IBM Logo.ch8")
	//chip8.LoadROM("roms/test_opcode.ch8")
	//chip8.Run()

	romFilepath := "roms/IBM Logo.ch8"
	fmt.Printf("Disassembly of \"%s\":\n", romFilepath)
	chip8.DisassembleProgram(romFilepath, 0x200)

	fmt.Println("Done.")
}
