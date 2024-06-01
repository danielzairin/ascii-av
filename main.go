package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

var HEADER_SIZE = 10

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// ID3 2.4.0 specs: https://mutagen-specs.readthedocs.io/en/latest/id3/id3v2.4.0-structure.html
func main() {
	file, err := os.Open("song.mp3")
	defer file.Close()
	check(err)

	reader := bufio.NewReader(file)
	header := make([]byte, HEADER_SIZE)
	for i := 0; i < HEADER_SIZE; i++ {
		header[i], err = reader.ReadByte()
		check(err)
	}

	fmt.Printf("ID3 tag header: %v\n", header)

	last4Bytes := [4]byte(header[HEADER_SIZE-4:])
	tagSize := HEADER_SIZE + int(parseSynchsafeInt(last4Bytes))

	fmt.Printf("ID3 tag size: %d\n", tagSize)
}

// TODO: Test this
func parseSynchsafeInt(bytes [4]byte) uint32 {
	for i := len(bytes) - 1; i >= 0; i-- {
		if i > 0 && bytes[i-1] > 0 {
			lsbOfNextByte := bytes[i-1] & 0b00000001
			if lsbOfNextByte == 1 {
				bytes[i] = bytes[i] | 0b10000000 // set the most significant bit to 1
			}
			bytes[i-1] = bytes[i-1] >> 1 // shift the next byte by 1 bit to the right
		}
	}

	return binary.BigEndian.Uint32(bytes[:])
}
