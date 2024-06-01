package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

// Thank you internet:
// ID3 2.4.0 specs - https://mutagen-specs.readthedocs.io/en/latest/id3/id3v2.4.0-structure.html
// MP3 header specs - http://mpgedit.org/mpgedit/mpeg_format/mpeghdr.htm

const ID3_HEADER_SIZE = 10
const MP3_HEADER_SIZE = 10

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	file, err := os.Open("song.mp3")
	defer file.Close()
	check(err)

	reader := bufio.NewReader(file)
	var id3Header [ID3_HEADER_SIZE]byte
	for i := 0; i < ID3_HEADER_SIZE; i++ {
		id3Header[i], err = reader.ReadByte()
		check(err)
	}

	fmt.Printf("ID3 tag header: %b\n", id3Header)

	last4Bytes := [4]byte(id3Header[ID3_HEADER_SIZE-4:])
	bodySize := int(parseSynchsafeInt(last4Bytes))

	skipped, err := reader.Discard(bodySize)

	fmt.Printf("Skipped: %d bytes\n", skipped)

	var mp3Header [MP3_HEADER_SIZE]byte
	for i := 0; i < len(mp3Header); i++ {
		mp3Header[i], err = reader.ReadByte()
		check(err)
	}

	fmt.Printf("MP3 header: %X\n", mp3Header)
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
