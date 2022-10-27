package chip8

import (
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"net"
	"os"
	"sync"
)

type Peripherals struct {
	connection net.Conn
	lock       sync.Mutex
}

type PeripheralsState struct {
	sound        bool
	keys         uint16
	screenBuffer [64][32]uint8
	screenWidth  uint8
	screenHeight uint8
}

func NewPeripherals() Peripherals {
	// In IPv4, any address between 224.0.0.0 to 239.255.255.255 can be used as a multicast address.
	address := "230.0.0.0:9999"

	connection, err := net.Dial("udp", address)
	if err != nil {
		fmt.Printf("Could not create multicast connection to render monitor %v", err)
		os.Exit(2)
	}

	return Peripherals{connection: connection}
}

func (p *Peripherals) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	if err := p.connection.Close(); err != nil {
		fmt.Printf("could close connection: %s\n", err.Error())
	}
}

func (p *Peripherals) Update(ps PeripheralsState) {
	serializedMessage := getSerializedMessage(ps)

	if _, err := p.connection.Write(serializedMessage); err != nil {
		fmt.Printf("could not update peripherals: %s\n%+v\n", err.Error(), ps)
	}
}

func getSerializedMessage(state PeripheralsState) []byte {
	width := int(state.screenWidth)
	height := int(state.screenHeight)
	screenBitBuffer := make([]byte, width*height/8)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixelIndex := y*width + x
			byteIndex := pixelIndex / 8
			bitIndex := 7 - pixelIndex%8
			screenBitBuffer[byteIndex] |= (state.screenBuffer[x][y] & 0b00000001) << bitIndex
		}
	}

	message := struct {
		Sound        bool   `msgpack:"sound"`
		Keys         int    `msgpack:"keys"`
		Screen       []byte `msgpack:"screen"`
		ScreenWidth  int    `msgpack:"screenWidth"`
		ScreenHeight int    `msgpack:"screenHeight"`
	}{
		Sound:        state.sound,
		Keys:         int(state.keys),
		Screen:       screenBitBuffer,
		ScreenWidth:  width,
		ScreenHeight: height,
	}

	serializedMessage, err := msgpack.Marshal(&message)
	if err != nil {
		fmt.Printf("Could not marshal data: %+v\n", message)
	}

	return serializedMessage
}
