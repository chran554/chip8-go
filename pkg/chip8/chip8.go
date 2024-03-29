package chip8

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

const romAddressDefault = 0x200
const romAddressEti660 = 0x600
const fontAddressDefault = 0x050
const flagRegisterIndex = 0xF

type Configuration struct {
	Debug                 bool // Debug mode prints, in more or less natural language, the instructions performed during the program execution
	ModeRomCompatibility  bool // ModeRomCompatibility The preferred mode setting for most ROM compatibility
	ModeStrictCosmac      bool // ModeStrictCosmac infers strict original instruction execution as COSMAC was designed (far from all ROM adhere to this)
	EndOnInfiniteLoop     bool // EndOnInfiniteLoop ends the program if an infinite loop is detected (some program ends with infinite loop and require restart to run again)
	RestartOnInfiniteLoop bool // RestartOnInfiniteLoop restarts the program if an infinite loop is detected (some program ends with infinite loop and require restart to run again)

	Disassemble          bool // Disassemble do execute the ROM program but rather prints it to stdout with, more or less, natural language explanation to each instruction
	DisassembleEveryByte bool // DisassembleEveryByte try to disassemble instructions at all bytes not just at even addresses. Some programs have parts of the code based at uneven addresses.
}

type Chip8 struct {
	Memory           []byte
	PC               uint16
	I                uint16
	Stack            stack
	Timer            uint8
	SoundTimer       uint8
	V                []uint8
	fontStartAddress uint16
	peripherals      *Peripherals
}

func NewChip8(peripherals *Peripherals) *Chip8 {
	chip8 := Chip8{
		Memory:           make([]byte, 0xFFF+1), // 4kB of memory (0x000-0xFFF)
		PC:               romAddressDefault,
		I:                0,
		Stack:            newStack(12), // Original RCA 1802 implementation had 12 levels of nesting
		Timer:            0,
		SoundTimer:       0,
		V:                make([]uint8, 0xF+1), // 16 registers of 8 bit each. Named V0,V1,..,V9,VA,..,VF
		fontStartAddress: fontAddressDefault,
		peripherals:      peripherals,
	}

	addFont(chip8)

	return &chip8
}

func (chip8 *Chip8) Run(configuration Configuration) {
	var err error

	go timerCounter(chip8)

	for true {
		// Processor stage: Fetch

		time.Sleep(time.Duration(3) * time.Millisecond)

		// Chip8 is big endian
		instructionCode := uint16(chip8.Memory[chip8.PC])<<8 | uint16(chip8.Memory[chip8.PC+1])
		if configuration.Debug {
			printInstructionDebugInfo(chip8.PC, instructionCode, configuration)
		}

		// Processor stage: Decode(-ish)

		instructionType := uint8((instructionCode & 0xF000) >> 12)
		x := uint8((instructionCode & 0x0F00) >> 8)
		y := uint8((instructionCode & 0x00F0) >> 4)
		z := uint8((instructionCode & 0x000F) >> 0)
		n := uint8(instructionCode & 0x000F)
		nn := uint8(instructionCode & 0x00FF)
		nnn := instructionCode & 0x0FFF

		// Processor stage: Execute

		chip8.PC += 2

		switch instructionType {
		case 0x0:
			if nnn == 0x0EE {
				// 00EE: Return from a subroutine
				chip8.PC, err = chip8.Stack.Pop()
				if err != nil {
					fmt.Printf("Error returning from subroutine (popping return address): %s\n", err.Error())
				}
			} else if nnn == 0x0E0 {
				// 00E0: Clear screen
				chip8.peripherals.state.screen.Clear()
				go chip8.UpdateScreen()
			} else {
				fmt.Println("Machine code execution not available/not implemented")
				os.Exit(1)
			}

		case 0x1:
			// 1NNN: Jump to address NNN
			if configuration.EndOnInfiniteLoop && ((chip8.PC - 2) == nnn) {
				fmt.Println("Terminated emulator and program on detected infinite loop")

				chip8.peripherals.state.sound = false
				chip8.peripherals.state.keys = 0b0000000000000000
				chip8.UpdateSoundAndKeys()

				os.Exit(0)
			}

			chip8.PC = nnn

		case 0x2:
			// 2NNN: Jump to subroutine (see also 00EE)
			err := chip8.Stack.Push(chip8.PC)
			if err != nil {
				fmt.Printf("Error jumping to subroutine (pushing return address): %s\n", err.Error())
			}
			chip8.PC = nnn

		case 0x3:
			// 3XNN: Skip next instruction if register X equals NN (see also 4XNN)
			if chip8.V[x] == nn {
				chip8.PC += 2
			}

		case 0x4:
			// 4XNN: Skip next instruction if register X NOT equals NN (see also 3XNN)
			if chip8.V[x] != nn {
				chip8.PC += 2
			}

		case 0x5:
			if z == 0 {
				// 5XY0: Skip next instruction if register X equals register Y (see also 9XY0)
				if chip8.V[x] == chip8.V[y] {
					chip8.PC += 2
				}
			}

		case 0x6:
			// 6XNN: Set register X to value NN
			chip8.V[x] = nn

		case 0x7:
			// 7XNN: Add the value NN to VX.
			// NOTE: overflow flag is not affected by this instruction if result > 0xFF. If result wraps to zero when overflow i.e. VX = (VX + NN) % 0xFF.
			chip8.V[x] += nn

		case 0x8:
			if z == 0x0 {
				// 8XY0: Set VX to the value of VY
				chip8.V[x] = chip8.V[y]
			} else if z == 0x1 {
				// 8XY1: VX is set to the bitwise/binary logical disjunction (OR) of VX and VY. VY is not affected.
				chip8.V[x] |= chip8.V[y]
			} else if z == 0x2 {
				// 8XY2: VX is set to the bitwise/binary logical conjunction (AND) of VX and VY. VY is not affected.
				chip8.V[x] &= chip8.V[y]
			} else if z == 0x3 {
				// 8XY3: VX is set to the bitwise/binary exclusive OR (XOR) of VX and VY. VY is not affected.
				chip8.V[x] ^= chip8.V[y]
			} else if z == 0x4 {
				// 8XY4: VX is set to the value of VX plus the value of VY. VY is not affected. Carry flag in register VF is set if overflow
				result := uint16(chip8.V[x]) + uint16(chip8.V[y])
				if result > 0xFF {
					chip8.V[flagRegisterIndex] = 1
				} else {
					chip8.V[flagRegisterIndex] = 0
				}
				chip8.V[x] = uint8(result % 0x100)
			} else if z == 0x5 {
				// 8XY5: subtract VY from VX and put the result in VX. VY is not affected.
				if chip8.V[x] > chip8.V[y] {
					chip8.V[flagRegisterIndex] = 1
				} else {
					chip8.V[flagRegisterIndex] = 0
				}
				chip8.V[x] = chip8.V[x] - chip8.V[y]
			} else if z == 0x6 {
				// 8XY6: (Strict COSMAC: Copy VY to VX and) shift VX 1 bit to the right. VF is set to the bit that was shifted out.
				if configuration.ModeStrictCosmac {
					chip8.V[x] = chip8.V[y]
				}
				chip8.V[flagRegisterIndex] = (chip8.V[x] & 0b00000001) >> 0
				chip8.V[x] = chip8.V[x] >> 1
			} else if z == 0x7 {
				// 8XY7: subtract VX from VY and put the result in VX. VY is not affected.
				if chip8.V[y] > chip8.V[x] {
					chip8.V[flagRegisterIndex] = 1
				} else {
					chip8.V[flagRegisterIndex] = 0
				}
				chip8.V[x] = chip8.V[y] - chip8.V[x]
			} else if z == 0xE {
				// 8XYE: (Strict COSMAC: Copy VY to VX and) shift VX 1 bit to the left. VF is set to the bit that was shifted out.
				if configuration.ModeStrictCosmac {
					chip8.V[x] = chip8.V[y]
				}
				chip8.V[flagRegisterIndex] = (chip8.V[x] & 0b10000000) >> 7
				chip8.V[x] = chip8.V[x] << 1
			}

		case 0x9:
			if z == 0x0 {
				// 9XY0: Skip next instruction if register X NOT equals register Y (see also 5XY0)
				if chip8.V[x] != chip8.V[y] {
					chip8.PC += 2
				}
			}

		case 0xA:
			// ANNN: Sets the index register I to the value NNN.
			chip8.I = nnn

		case 0xB:
			if configuration.ModeStrictCosmac || configuration.ModeRomCompatibility {
				// BNNN: Jump to the address NNN plus the value in the register V0.
				chip8.PC = nnn + uint16(chip8.V[0x0])
			} else {
				// B(X)NNN: Jump to the address NNN plus the value in the register VX.
				chip8.PC = nnn + uint16(chip8.V[x])
			}

		case 0xC:
			// CXNN: Generates a random number, binary ANDs it with the value NN, and puts the result in VX.
			chip8.V[x] = uint8(rand.Uint32()&0x000000FF) & nn

		case 0xD:
			// DXYN: Draw an N pixels tall sprite from the memory location that the I-index register is holding to the screen,
			// at the horizontal X coordinate in VX and the Y coordinate in VY.
			pixelX := chip8.V[x] % chip8.peripherals.state.screen.Width
			pixelY := chip8.V[y] % chip8.peripherals.state.screen.Height
			chip8.V[flagRegisterIndex] = 0

			for spriteY := uint8(0); spriteY < n; spriteY++ {
				pixelBitValues := chip8.Memory[chip8.I+uint16(spriteY)]
				for spriteX := uint8(0); spriteX < 8; spriteX++ {
					if (spriteX < chip8.peripherals.state.screen.Width) && (spriteY < chip8.peripherals.state.screen.Height) {
						pixelValue := (pixelBitValues >> spriteX) & 0b00000001

						resultPixelValue := chip8.peripherals.state.screen.XorPixel(pixelX+(7-spriteX), pixelY+spriteY, pixelValue)
						if (pixelValue == 1) && (resultPixelValue == 0) {
							chip8.V[flagRegisterIndex] = 1
						}
					}
				}
			}

			go chip8.UpdateScreen()

		case 0xE:
			if nn == 0x9E {
				// EX9E: Skip next instruction if key denoted by VX is pressed at the moment
				if chip8.isKeyPressed(chip8.V[x]) {
					chip8.PC += 2
				}
			} else if nn == 0xA1 {
				// EXA1: Skip next instruction if key denoted by VX is NOT pressed at the moment
				if !chip8.isKeyPressed(chip8.V[x]) {
					chip8.PC += 2
				}
			}

		case 0xF:
			if nn == 0x07 {
				// FX07: Sets VX to the current value of the delay timer
				chip8.V[x] = chip8.Timer
			} else if nn == 0x15 {
				// FX15: Sets the delay timer to the value in VX
				chip8.Timer = chip8.V[x]
			} else if nn == 0x18 {
				// FX18: Sets the sound timer to the value in VX
				chip8.UpdateSound(chip8.V[x] > 0)
				chip8.SoundTimer = remappedSoundValue(chip8.V[x])
				//fmt.Printf("Sound on (value %d, sound timer set to %d, %d msec)\n", chip8.V[x], chip8.SoundTimer, int(math.Round(float64(chip8.SoundTimer)*1000.0/60.0)))
			} else if nn == 0x1E {
				// FX1E: Add to index. The index register I will get the value in VX added to it.
				result := chip8.I + uint16(chip8.V[x])

				if !configuration.ModeStrictCosmac {
					if result > 0xFFF {
						// Register I would point outside memory range
						chip8.V[flagRegisterIndex] = 1
					} else {
						chip8.V[flagRegisterIndex] = 0
					}
				}
				chip8.I = result & 0x0FFF
			} else if nn == 0x0A {
				// FX0A: This instruction "blocks", it stops executing instructions and wait for key input. Value of key is stored in VX.
				pressedKeyCode := chip8.getPressedKey()
				if pressedKeyCode != 0xFF {
					chip8.V[x] = pressedKeyCode
				} else {
					chip8.PC -= 2 // Do not advance in program, do this instruction over again (loop)
				}
			} else if nn == 0x29 {
				// FX29: Set index register to point at font character address. The character code is stored in VX
				// Each character is 5 bytes in height
				chip8.I = chip8.fontStartAddress + (uint16(chip8.V[x]) * 5)
			} else if nn == 0x33 {
				// FX33: Binary-coded decimal conversion
				// It takes the number in VX and converts it to three decimal digits,
				// storing these digits in memory at the start address in the index register I.
				chip8.Memory[chip8.I+0] = (chip8.V[x] / 100) % 10
				chip8.Memory[chip8.I+1] = (chip8.V[x] / 10) % 10
				chip8.Memory[chip8.I+2] = (chip8.V[x] / 1) % 10
			} else if nn == 0x55 {
				// FX55: Store V registers in memory
				// The value of each variable register from V0 to VX inclusive
				// (if X is 0, then only V0) will be stored in successive memory addresses,
				// starting with the one that’s pointed to by register I.
				if configuration.ModeStrictCosmac && !configuration.ModeRomCompatibility {
					for i := uint8(0); (i <= x) && (i <= 0xF); i++ {
						chip8.Memory[chip8.I] = chip8.V[i]
						chip8.I++
					}
				} else {
					for i := uint8(0); (i <= x) && (i <= 0xF); i++ {
						chip8.Memory[chip8.I+uint16(i)] = chip8.V[i]
					}
				}
			} else if nn == 0x65 {
				// FX65: Load registers from memory
				// Takes the value stored at the memory addresses and loads them into the variable registers.
				if configuration.ModeStrictCosmac && !configuration.ModeRomCompatibility {
					for i := uint8(0); (i <= x) && (i <= 0xF); i++ {
						chip8.V[i] = chip8.Memory[chip8.I]
						chip8.I++
					}
				} else {
					for i := uint8(0); (i <= x) && (i <= 0xF); i++ {
						chip8.V[i] = chip8.Memory[chip8.I+uint16(i)]
					}
				}
			}

		default:
			panic(fmt.Sprintf("unknown instruction \"0x%X\" at address 0x%X", instructionCode, chip8.PC-2))
		}

	}
}

func remappedSoundValue(soundDelay uint8) uint8 {
	if soundDelay == 0 {
		return 0
	}

	minimalDelay := 5.0
	v := math.Pow(float64(soundDelay-1), 0.85) / math.Pow(254, 0.85)
	return uint8(math.Round(v*(255.0-minimalDelay) + minimalDelay))
}

func (chip8 *Chip8) _loadROM(filepath string, startAddress int) {
	romBytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("could not load ROM file \"%s\": %s\n", filepath, err.Error())
		os.Exit(1)
	}

	for i, romByte := range romBytes {
		chip8.Memory[startAddress+i] = romByte
	}
}

func (chip8 *Chip8) LoadROM(filepath string) {
	chip8._loadROM(filepath, romAddressDefault)
}

func (chip8 *Chip8) LoadETI660ROM(filepath string) {
	chip8._loadROM(filepath, romAddressEti660)
}

func (chip8 *Chip8) UpdateScreen() {
	// go chip8.peripherals.state.screen.Print()
	chip8.peripherals.UpdateScreen()
}

func (chip8 *Chip8) UpdateSound(soundState bool) {
	chip8.peripherals.UpdateSound(soundState)
}

func (chip8 *Chip8) UpdateSoundAndKeys() {
	chip8.peripherals.UpdateSoundAndKeys()
}

func (chip8 *Chip8) getPressedKey() uint8 {
	keySate := chip8.peripherals.state.keys

	if keySate > 0 {
		for keyIndex := uint8(0); keyIndex <= 0xF; keyIndex++ {
			if (keySate>>keyIndex)&0x1 == 1 {
				return keyIndex
			}
		}
	}

	return 0xFF // return key code of pressed key or 0xFF for "no key pressed"
}

func (chip8 *Chip8) isKeyPressed(keyCode uint8) bool {
	// fmt.Printf("Checking for key: %1X    %016b\n", keyCode, chip8.peripherals.state.keys)
	return (chip8.peripherals.state.keys>>keyCode)&0x1 == 1
}

func addFont(chip Chip8) {
	font := []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	for i, b := range font {
		chip.Memory[chip.fontStartAddress+uint16(i)] = b
	}
}

func timerCounter(chip8 *Chip8) {
	var countDownFrequency = 1000 / 60 // 60 Hz
	for {
		time.Sleep(time.Duration(countDownFrequency) * time.Millisecond)
		if chip8.Timer > 0 {
			chip8.Timer--
		}
		if chip8.SoundTimer > 0 {
			chip8.SoundTimer--
		}

		chip8.UpdateSound(chip8.SoundTimer > 0)
	}
}
