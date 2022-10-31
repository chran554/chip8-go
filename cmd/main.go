package main

import (
	"chip8/pkg/chip8"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("You need to supply at least 1 argument to the program. The file path to a ROM file.")
		os.Exit(1)
	}
	romFilepath := os.Args[1]

	configuration := chip8.Configuration{
		Disassemble:          false,
		Debug:                false,
		EndOnInfiniteLoop:    true,
		ModeRomCompatibility: true,
		ModeStrictCosmac:     false,
	}

	if !configuration.Disassemble {
		fmt.Printf("CHIP-8 execution of \"%s\":\n", romFilepath)
		fmt.Printf("%+v\n", configuration)

		peripherals := chip8.NewPeripherals()
		peripherals.StartKeyPadListener()

		machine := chip8.NewChip8(&peripherals)
		machine.LoadROM(romFilepath)
		machine.Run(configuration)
	} else {
		fmt.Printf("CHIP-8 disassembly of \"%s\":\n", romFilepath)
		chip8.DisassembleProgram(romFilepath, 0x200, configuration)
	}
}
