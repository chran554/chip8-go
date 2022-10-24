package chip8

import (
	"fmt"
	"regexp"
)

func printInstructionDebugInfo(address uint16, instruction uint16) {
	// "eternal loop" == "jump to the same address"
	eternalLoop := (((instruction & 0xF000) >> 12) == 1) && (address == (instruction & 0x0FFF))

	if !eternalLoop {
		fmt.Printf("0x%03X: %04X   # %s\n", address, instruction, explanation(instruction))
	}
}

func explanation(instruction uint16) string {
	instructionText := fmt.Sprintf("%04X", instruction)
	regExp_00E0 := regexp.MustCompile("00E0")
	regExp_00EE := regexp.MustCompile("00EE")
	regExp_6XNN := regexp.MustCompile("6(\\w)(\\w\\w)")
	regExp_7XNN := regexp.MustCompile("7(\\w)(\\w\\w)")
	regExp_ANNN := regexp.MustCompile("A(\\w\\w\\w)")
	regExp_DXYN := regexp.MustCompile("D(\\w)(\\w)(\\w)")
	regExp_1NNN := regexp.MustCompile("1(\\w\\w\\w)")

	if regExp_00E0.MatchString(instructionText) {
		return "Clear screen"
	}

	if regExp_00EE.MatchString(instructionText) {
		return "Return from subroutine"
	}

	if regExp_1NNN.MatchString(instructionText) {
		matches := regExp_1NNN.FindStringSubmatch(instructionText)
		return fmt.Sprintf("Jump to address 0x%s", matches[1])
	}

	if regExp_6XNN.MatchString(instructionText) {
		matches := regExp_6XNN.FindStringSubmatch(instructionText)
		return fmt.Sprintf("Set register V%s to value 0x%s", matches[1], matches[2])
	}

	if regExp_7XNN.MatchString(instructionText) {
		matches := regExp_7XNN.FindStringSubmatch(instructionText)
		return fmt.Sprintf("Add value 0x%s to register V%s", matches[2], matches[1])
	}

	if regExp_ANNN.MatchString(instructionText) {
		matches := regExp_ANNN.FindStringSubmatch(instructionText)
		return fmt.Sprintf("Set register I to point at address 0x%s", matches[1])
	}

	if regExp_DXYN.MatchString(instructionText) {
		matches := regExp_DXYN.FindStringSubmatch(instructionText)
		return fmt.Sprintf("Xor draw sprite of pixel size 8x%s, from address pointed to by register I, at screen position (V%s, V%s)", matches[3], matches[1], matches[2])
	}

	return ""
}
