package main

import (
	"bufio"
	"errors"
	"os"
	unsafe "unsafe"

	"github.com/cehlen/demoinfogolang/util"
)

const (
	// Static info
	MAGIC_NUMBER string = "HL2DEMO"
	DEMO_PROTO   int32  = 4

	// Messages

)

type DemoHeader struct {
	fileInfo               [8]byte
	protocolVersion        int32
	networkProtocolVersion int32
	server                 [260]byte
	client                 [260]byte
	mapName                [260]byte
	gameDir                [260]byte
	playbackTime           float32
	playbackTicks          int32
	playbackFrames         int32
	signonLength           int32
}

type DemoFile struct {
	FileBuffer []byte
	bufferPos  int
	fileName   string
	DemoHeader DemoHeader
}

// Open starts the process of parsing the demo file.
func (df *DemoFile) Open(fileName string) error {
	// We want to reset the data / close the previous opened file
	df.Reset()

	// Open file
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	// Get the length to do some basic checking
	fileStats, err := f.Stat()
	if err != nil {
		return err
	}
	length := fileStats.Size()

	// Check if it can be a valid demo by compaing the size of the header to
	// the size of the file
	headerSize := (int64)(unsafe.Sizeof(df.DemoHeader))
	if length < headerSize {
		return errors.New("File to small")
	}

	// Go to the beginning again
	f.Seek(0, 0)
	reader := bufio.NewReader(f)

	// Try to read the header
	headerBytes := make([]byte, headerSize)
	_, err = reader.Read(headerBytes)
	if err != nil {
		return err
	}

	df.fillHeader(headerBytes)
	return nil
}

func (df *DemoFile) fillHeader(header []byte) error {
	var hdr []byte = header
	copy(df.DemoHeader.fileInfo[:], hdr[:7])
	var err error
	df.DemoHeader.protocolVersion, err = util.ByteSliceToInt32(hdr[8:12])
	if err != nil {
		return err
	}
	df.DemoHeader.networkProtocolVersion, err = util.ByteSliceToInt32(hdr[12:16])
	if err != nil {
		return err
	}

	copy(df.DemoHeader.server[:], hdr[16:276])
	copy(df.DemoHeader.client[:], hdr[276:536])
	copy(df.DemoHeader.mapName[:], hdr[536:796])
	copy(df.DemoHeader.gameDir[:], hdr[796:1056])

	// playback
	df.DemoHeader.playbackTime = util.ByteSliceToFloat32(hdr[1056:1060])
	df.DemoHeader.playbackTicks, err = util.ByteSliceToInt32(hdr[1060:1064])
	if err != nil {
		return err
	}
	df.DemoHeader.playbackTicks, err = util.ByteSliceToInt32(hdr[1064:1068])
	if err != nil {
		return err
	}
	df.DemoHeader.signonLength, err = util.ByteSliceToInt32(hdr[1068:])
	if err != nil {
		return err
	}
	return nil
}

// Reset resets the structure so we can use it multiple times
func (df *DemoFile) Reset() {
	df.fileName = ""

	df.bufferPos = 0
	df.FileBuffer = make([]byte, 1)
}