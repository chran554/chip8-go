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
	state      *PeripheralsState
	connection net.Conn
	lock       sync.Mutex
}

type PeripheralsState struct {
	sound  bool
	keys   uint16       // keys are a 16 bit bitmask for all pressed keys, "0" through "F"
	screen ScreenBuffer // screen is the screen memory, the pixel memory representation
}

type peripheralStateMessage struct {
	Sound        bool   `msgpack:"sound"`
	Keys         uint16 `msgpack:"keys"`
	Screen       []byte `msgpack:"screen"`
	ScreenWidth  byte   `msgpack:"screenWidth"`
	ScreenHeight byte   `msgpack:"screenHeight"`
}

func NewPeripherals() Peripherals {
	// In IPv4, any address between 224.0.0.0 to 239.255.255.255 can be used as a multicast address.
	address := "230.0.0.0:9999"

	connection, err := net.Dial("udp", address)
	if err != nil {
		fmt.Printf("Could not create multicast connection to render monitor %v", err)
		os.Exit(2)
	}

	state := PeripheralsState{
		sound:  false,              // No sound
		keys:   0b0000000000000000, // No keys pressed
		screen: NewScreenBuffer(),  // Empty (black) screen
	}

	return Peripherals{connection: connection, state: &state}
}

func (p *Peripherals) StartKeyPadListener() {
	go listenForPeripheralKeyPadInput("230.0.0.0:9998", p)
}

func listenForPeripheralKeyPadInput(address string, p *Peripherals) {
	keyPadMaxDatagramSize := 256
	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(err)
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err = conn.SetReadBuffer(keyPadMaxDatagramSize); err != nil {
		log.Fatalf("could not set receive buffer size: %s", err.Error())
	}

	buffer := make([]byte, keyPadMaxDatagramSize)

	// Loop forever reading from the socket
	for {
		numBytes, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}

		if numBytes != 2 {
			log.Fatalf("Illegal input length: %d (expected 2)", numBytes)
		}

		keyPadState := (uint16(buffer[0]) << 8) | (uint16(buffer[1]) << 0) // Convert byte input data to key pad state
		//fmt.Printf("Got key state: %016b\n", keyPadState)
		p.state.keys = keyPadState
	}
}

func (p *Peripherals) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	if err := p.connection.Close(); err != nil {
		fmt.Printf("could close connection: %s\n", err.Error())
	}
}

func (p *Peripherals) UpdateSound(newSoundState bool) {
	if p.state.sound != newSoundState {
		p.state.sound = newSoundState
		serializedMessage := getSerializedSoundAndKeysMessage(p.state)

		if _, err := p.connection.Write(serializedMessage); err != nil {
			fmt.Printf("could not update peripherals: %s\n%+v\n", err.Error(), p.state)
		}
	}
}

func (p *Peripherals) UpdateSoundAndKeys() {
	serializedMessage := getSerializedSoundAndKeysMessage(p.state)
	if _, err := p.connection.Write(serializedMessage); err != nil {
		fmt.Printf("could not update peripherals: %s\n%+v\n", err.Error(), p.state)
	}
}

func (p *Peripherals) UpdateScreen() {
	serializedMessage := getSerializedScreenMessage(p.state)

	if _, err := p.connection.Write(serializedMessage); err != nil {
		fmt.Printf("could not update peripherals: %s\n%+v\n", err.Error(), p.state)
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
