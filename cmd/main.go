package main

import (
	"chip8/pkg/chip8"
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("You need to supply at least 1 argument to the program. The file path to a ROM file.")
		os.Exit(1)
	}

	screenAddress := flag.String("screenAddress", "localhost:9999", "The socket address of the screen application. Format: \"127.0.0.1:9999\". Default value: \"127.0.0.1:9999\".")
	listenKeyStatePort := flag.Int("keystatePort", 9998, "The port where to listen for key press state changes. Format: \"9998\". Default value \"9998\".")
	flag.Parse() // add this line
	romFilepath := flag.Arg(0)

	configuration := chip8.Configuration{
		Disassemble:          false,
		Debug:                false,
		EndOnInfiniteLoop:    true,
		ModeRomCompatibility: true,
		ModeStrictCosmac:     false,
	}

	if !configuration.Disassemble {
		fmt.Println()
		fmt.Printf("CHIP-8 execution of ROM file \"%s\"\n", romFilepath)
		fmt.Printf("Using screen address:                    %s\n", *screenAddress)
		fmt.Printf("Listening to key state changes on port:  %d\n", *listenKeyStatePort)
		fmt.Println()
		fmt.Printf("Configuration:\n%+v\n", configuration)
		fmt.Println()

		peripherals := chip8.NewPeripherals(*screenAddress, *listenKeyStatePort)
		peripherals.StartKeyPadListener()

		machine := chip8.NewChip8(&peripherals)
		machine.LoadROM(romFilepath)
		machine.Run(configuration)
	} else {
		fmt.Printf("CHIP-8 disassembly of \"%s\":\n", romFilepath)
		fmt.Printf("%+v\n", configuration)
		chip8.DisassembleProgram(romFilepath, 0x200, configuration)
	}
}
