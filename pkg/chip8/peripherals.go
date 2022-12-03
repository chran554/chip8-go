package chip8

import (
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"net"
	"os"
	"sync"
)

type Peripherals struct {
	state                *PeripheralsState
	screenConnection     net.Conn
	lock                 sync.Mutex
	keyStateListenerPort int
}

type PeripheralsState struct {
	sound  bool
	keys   uint16       // keys are a 16 bit bitmask for all pressed keys, "0" through "F"
	screen ScreenBuffer // ScreenBuffer is the screen memory, the pixel memory representation
}

type peripheralStateMessage struct {
	Sound        bool   `msgpack:"sound"`
	Keys         uint16 `msgpack:"keys"`
	Screen       []byte `msgpack:"screen"`
	ScreenWidth  byte   `msgpack:"screenWidth"`
	ScreenHeight byte   `msgpack:"screenHeight"`
}

func NewPeripherals(screenAddress string, keyStateListenerPort int) Peripherals {
	screenConnection, err := net.Dial("udp", screenAddress)
	if err != nil {
		fmt.Printf("Could not create screenConnection to screen %v", err)
		os.Exit(2)
	}

	state := PeripheralsState{
		sound:  false,              // No sound
		keys:   0b0000000000000000, // No keys pressed
		screen: NewScreenBuffer(),  // Empty (black) screen
	}

	return Peripherals{
		screenConnection:     screenConnection,
		keyStateListenerPort: keyStateListenerPort,
		state:                &state,
	}
}

func (p *Peripherals) StartKeyPadListener() {
	go listenForPeripheralKeyPadInput(p)
}

func listenForPeripheralKeyPadInput(p *Peripherals) {
	keyPadMaxDatagramSize := 256

	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", p.keyStateListenerPort))
	sock, _ := net.ListenUDP("udp", addr)
	sock.SetReadBuffer(keyPadMaxDatagramSize)

	buffer := make([]byte, keyPadMaxDatagramSize)

	// Loop forever reading from the socket
	for {
		numBytes, _, err := sock.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("Read from UDP failed:", err)
		}

		if numBytes != 2 {
			log.Fatalf("chip-8 key state listener: illegal input length: %d bytes (expected 2 bytes)", numBytes)
		}

		keyPadState := (uint16(buffer[0]) << 8) | (uint16(buffer[1]) << 0) // Convert byte input data to key pad state
		// fmt.Printf("Got key state: %016b\n", keyPadState)
		p.state.keys = keyPadState
	}
}

func (p *Peripherals) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	if err := p.screenConnection.Close(); err != nil {
		fmt.Printf("could close screenConnection: %s\n", err.Error())
	}
}

func (p *Peripherals) UpdateSound(newSoundState bool) {
	if p.state.sound != newSoundState {
		p.state.sound = newSoundState
		serializedMessage := getSerializedSoundAndKeysMessage(p.state)

		if _, err := p.screenConnection.Write(serializedMessage); err != nil {
			fmt.Printf("could not update peripherals sound (plus key) state: %s\n", err.Error())
			fmt.Println("(is screen application up and running?)")
		}
	}
}

func (p *Peripherals) UpdateSoundAndKeys() {
	serializedMessage := getSerializedSoundAndKeysMessage(p.state)
	if _, err := p.screenConnection.Write(serializedMessage); err != nil {
		fmt.Printf("could not update peripherals sound and key state: %s\n", err.Error())
		fmt.Println("(is screen application up and running?)")
	}
}

func (p *Peripherals) UpdateScreen() {
	serializedMessage := getSerializedScreenMessage(p.state)

	if _, err := p.screenConnection.Write(serializedMessage); err != nil {
		fmt.Printf("could not update peripherals screen (plus sound and key) state: %s\n", err.Error())
		fmt.Println("(is screen application up and running?)")
	}
}

func getSerializedSoundAndKeysMessage(state *PeripheralsState) []byte {
	message := getSoundAndKeysMessage(state)

	serializedMessage, err := msgpack.Marshal(&message)
	if err != nil {
		fmt.Printf("Could not marshal data: %+v\n", message)
	}

	return serializedMessage
}

func getSoundAndKeysMessage(state *PeripheralsState) *peripheralStateMessage {
	width := state.screen.Width
	height := state.screen.Height

	// Create struct as soon as possible to capture sound state
	message := peripheralStateMessage{
		Sound:        state.sound,
		Keys:         state.keys,
		Screen:       nil,
		ScreenWidth:  width,
		ScreenHeight: height,
	}

	return &message
}

func getSerializedScreenMessage(state *PeripheralsState) []byte {
	message := getSoundAndKeysMessage(state)

	width := state.screen.Width
	height := state.screen.Height
	screenBitBuffer := make([]byte, int(width)*int(height)/8)

	for y := uint8(0); y < height; y++ {
		for x := uint8(0); x < width; x++ {
			pixelIndex := int(y)*int(width) + int(x)
			byteIndex := pixelIndex / 8
			bitIndex := 7 - pixelIndex%8
			screenBitBuffer[byteIndex] |= (state.screen.buffer[x][y] & 0b00000001) << bitIndex
		}
	}

	message.Screen = screenBitBuffer

	serializedMessage, err := msgpack.Marshal(&message)
	if err != nil {
		fmt.Printf("Could not marshal data: %+v\n", message)
	}

	return serializedMessage
}
