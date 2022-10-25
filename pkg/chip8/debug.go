package chip8

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var instructionRegExp = map[string]*regexp.Regexp{
	"00E0": regexp.MustCompile("00E0"),
	"00EE": regexp.MustCompile("00EE"),
	"1NNN": regexp.MustCompile("1(\\w\\w\\w)"),
	"2NNN": regexp.MustCompile("2(\\w\\w\\w)"),
	"3XNN": regexp.MustCompile("3(\\w)(\\w\\w)"),
	"4XNN": regexp.MustCompile("4(\\w)(\\w\\w)"),
	"5XY0": regexp.MustCompile("5(\\w)(\\w)0"),
	"6XNN": regexp.MustCompile("6(\\w)(\\w\\w)"),
	"7XNN": regexp.MustCompile("7(\\w)(\\w\\w)"),
	"8XY0": regexp.MustCompile("8(\\w)(\\w)0"),
	"8XY1": regexp.MustCompile("8(\\w)(\\w)1"),
	"8XY2": regexp.MustCompile("8(\\w)(\\w)2"),
	"8XY3": regexp.MustCompile("8(\\w)(\\w)3"),
	"8XY4": regexp.MustCompile("8(\\w)(\\w)4"),
	"8XY5": regexp.MustCompile("8(\\w)(\\w)5"),
	"8XY6": regexp.MustCompile("8(\\w)(\\w)6"),
	"8XYE": regexp.MustCompile("8(\\w)(\\w)E"),
	"9XY0": regexp.MustCompile("9(\\w)(\\w)0"),
	"ANNN": regexp.MustCompile("A(\\w\\w\\w)"),
	"BNNN": regexp.MustCompile("B(\\w\\w\\w)"),
	"DXYN": regexp.MustCompile("D(\\w)(\\w)(\\w)"),
	"FX33": regexp.MustCompile("F(\\w)33"),
	"FX55": regexp.MustCompile("F(\\w)55"),
	"FX65": regexp.MustCompile("F(\\w)65"),
}

func DisassembleProgram(romFilepath string, startAddress uint16) {
	bytes := loadByteFile(romFilepath)

	for address := uint16(0); address < uint16(len(bytes)); address++ {

		binaryBitsText := strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%08b", bytes[address]), "0", "░"), "1", "█")
		if (address % 2) == 0 {
			instructionCode := uint16(bytes[address+0])<<8 | uint16(bytes[address+1])

			if matchesAnyInstructionSyntax(instructionCode) {
				fmt.Printf("0x%03X:  0x%02X  %s    %04X    %s\n", startAddress+address, bytes[address], binaryBitsText, instructionCode, explanation(instructionCode))
			} else {
				fmt.Printf("0x%03X:  0x%02X  %s\n", startAddress+address, bytes[address], binaryBitsText)
			}
		} else {
			fmt.Printf("0x%03X:  0x%02X  %s\n", startAddress+address, bytes[address], binaryBitsText)
		}

	}
}

func printInstructionDebugInfo(address uint16, instruction uint16) {
	// "eternal loop" == "jump to the same address"
	eternalLoop := (((instruction & 0xF000) >> 12) == 1) && (address == (instruction & 0x0FFF))

	if !eternalLoop {
		fmt.Printf("0x%03X: %04X   # %s\n", address, instruction, explanation(instruction))
	}
}

func matchesAnyInstructionSyntax(instruction uint16) bool {
	instructionText := fmt.Sprintf("%04X", instruction)

	for _, regexp := range instructionRegExp {
		if regexp.MatchString(instructionText) {
			return true
		}
	}

	return false
}

func explanation(instruction uint16) string {
	instructionText := fmt.Sprintf("%04X", instruction)

	if instructionRegExp["00E0"].MatchString(instructionText) {
		return "00E0: Clear screen"
	}

	if instructionRegExp["00EE"].MatchString(instructionText) {
		return "00EE: Return from subroutine"
	}

	if instructionRegExp["1NNN"].MatchString(instructionText) {
		matches := instructionRegExp["1NNN"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("1NNN: Jump to address 0x%s", matches[1])
	}

	if instructionRegExp["2NNN"].MatchString(instructionText) {
		matches := instructionRegExp["2NNN"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("2NNN: Jump to subroutine at address 0x%s", matches[1])
	}

	if instructionRegExp["3XNN"].MatchString(instructionText) {
		matches := instructionRegExp["3XNN"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("3XNN: Skip next instruction if register V%s equals 0x%s", matches[1], matches[2])
	}

	if instructionRegExp["4XNN"].MatchString(instructionText) {
		matches := instructionRegExp["4XNN"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("4XNN: Skip next instruction if register V%s NOT equals 0x%s", matches[1], matches[2])
	}

	if instructionRegExp["5XY0"].MatchString(instructionText) {
		matches := instructionRegExp["5XY0"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("5XY0: Skip next instruction if register V%s equals register V%s", matches[1], matches[2])
	}

	if instructionRegExp["6XNN"].MatchString(instructionText) {
		matches := instructionRegExp["6XNN"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("6XNN: Set register V%s to value 0x%s", matches[1], matches[2])
	}

	if instructionRegExp["7XNN"].MatchString(instructionText) {
		matches := instructionRegExp["7XNN"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("7XNN: Add value 0x%s to register V%s", matches[2], matches[1])
	}

	if instructionRegExp["8XY0"].MatchString(instructionText) {
		matches := instructionRegExp["8XY0"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("8XY0: V%s is set to value of V%s. V%s is not affected.", matches[1], matches[2], matches[2])
	}

	if instructionRegExp["8XY1"].MatchString(instructionText) {
		matches := instructionRegExp["8XY1"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("8XY1: V%s is set to the bitwise/binary logical disjunction (OR) of V%s and V%s. V%s is not affected.", matches[1], matches[1], matches[2], matches[2])
	}

	if instructionRegExp["8XY2"].MatchString(instructionText) {
		matches := instructionRegExp["8XY2"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("8XY2: V%s is set to the bitwise/binary logical conjunction (AND) of V%s and V%s. V%s is not affected.", matches[1], matches[1], matches[2], matches[2])
	}

	if instructionRegExp["8XY3"].MatchString(instructionText) {
		matches := instructionRegExp["8XY3"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("8XY3: V%s is set to the bitwise/binary exclusive OR (XOR) of V%s and V%s. V%s is not affected.", matches[1], matches[1], matches[2], matches[2])
	}

	if instructionRegExp["8XY4"].MatchString(instructionText) {
		matches := instructionRegExp["8XY4"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("8XY4: V%s is set to the value of V%s + V%s. V%s is not affected. Carry flag in register VF is set if overflow", matches[1], matches[1], matches[2], matches[2])
	}

	if instructionRegExp["8XY5"].MatchString(instructionText) {
		matches := instructionRegExp["8XY5"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("8XY5: subtract V%s from V%s and put the result in V%s. V%s is not affected.", matches[2], matches[1], matches[1], matches[2])
	}

	if instructionRegExp["8XY6"].MatchString(instructionText) {
		matches := instructionRegExp["8XY6"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("8XY6: Copy V%s to V%s and shift V%s 1 bit to the RIGHT. VF is set to the bit that was shifted out.", matches[2], matches[1], matches[1])
	}

	if instructionRegExp["8XYE"].MatchString(instructionText) {
		matches := instructionRegExp["8XYE"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("8XYE: Copy V%s to V%s and shift V%s 1 bit to the LEFT. VF is set to the bit that was shifted out.", matches[2], matches[1], matches[1])
	}

	if instructionRegExp["9XY0"].MatchString(instructionText) {
		matches := instructionRegExp["9XY0"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("9XY0: Skip next instruction if register V%s NOT equals register V%s", matches[1], matches[2])
	}

	if instructionRegExp["ANNN"].MatchString(instructionText) {
		matches := instructionRegExp["ANNN"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("ANNN: Set register I to point at address 0x%s", matches[1])
	}

	if instructionRegExp["BNNN"].MatchString(instructionText) {
		matches := instructionRegExp["BNNN"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("BNNN: Jump to address 0x%s plus offset found in register V0", matches[1])
	}

	if instructionRegExp["DXYN"].MatchString(instructionText) {
		matches := instructionRegExp["DXYN"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("DXYN: Xor draw sprite of pixel size 8x%s, from address pointed to by register I, at screen position (V%s, V%s)", matches[3], matches[1], matches[2])
	}

	if instructionRegExp["FX33"].MatchString(instructionText) {
		matches := instructionRegExp["FX33"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("FX33: Binary-coded decimal conversion, store decimal digits of value found in register V%s in addresses pointed to by register I, I+1, and I+2", matches[1])
	}

	if instructionRegExp["FX55"].MatchString(instructionText) {
		matches := instructionRegExp["FX55"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("FX65: Store registers V0 through V%s to memory locations pointed to by register I through I+V%s", matches[1], matches[1])
	}

	if instructionRegExp["FX65"].MatchString(instructionText) {
		matches := instructionRegExp["FX65"].FindStringSubmatch(instructionText)
		return fmt.Sprintf("FX65: Load registers V0 through V%s from memory locations pointed to by register I through I+V%s", matches[1], matches[1])
	}

	return ""
}

func loadByteFile(filepath string) []byte {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("could not load byte file \"%s\": %s\n", filepath, err.Error())
	}

	return bytes
}
