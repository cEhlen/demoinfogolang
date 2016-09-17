package demo

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	unsafe "unsafe"

	"github.com/cehlen/demoinfogolang/util"
)

const (
	// Static info
	MAGIC_NUMBER           string = "HL2DEMO"
	VALID_PROTOCOL_VERSION int32  = 4

	// Byte Sizes
	INT32_SIZE = 4
)

const (
	// Commands
	DEM_SIGNON       uint8 = iota + 1
	DEM_PACKET       uint8 = iota + 1
	DEM_SYNCTICK     uint8 = iota + 1
	DEM_CONSOLECMD   uint8 = iota + 1
	DEM_USERCMD      uint8 = iota + 1
	DEM_DATATABLES   uint8 = iota + 1
	DEM_STOP         uint8 = iota + 1
	DEM_CUSTOMDATA   uint8 = iota + 1
	DEM_STRINGTABLES uint8 = iota + 1
	DEM_LASTCMD      uint8 = DEM_STRINGTABLES
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

func (dh *DemoHeader) MagicNumber() string {
	return string(bytes.Trim(dh.fileInfo[:], "\x00"))
}

type DemoFile struct {
	FileBuffer []byte
	bufferPos  int
	FileName   string
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
	length -= headerSize

	// Check 'MagicNumber' and Version
	if df.DemoHeader.MagicNumber() != MAGIC_NUMBER {
		return errors.New("Magic Number does not match")
	}
	if df.DemoHeader.protocolVersion != VALID_PROTOCOL_VERSION {
		return errors.New("Protocol version not valid")
	}

	// Read file into buffer
	tmpBuf := make([]byte, length)
	_, err = reader.Read(tmpBuf)
	if err != nil {
		return err
	}
	df.FileBuffer = tmpBuf
	df.bufferPos = 0
	df.FileName = fileName

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
	df.FileName = ""

	df.bufferPos = 0
	df.FileBuffer = make([]byte, 1)
}

func (df *DemoFile) readRaw(numBytes int32) []byte {
	slice := df.FileBuffer[df.bufferPos:(df.bufferPos + int(numBytes))]
	df.bufferPos += int(numBytes)
	return slice
}

func (df *DemoFile) ReadRawData(length int32) ([]byte, error) {
	if len(df.FileBuffer) == 0 {
		return nil, nil
	}

	blockSize, err := util.ByteSliceToInt32(df.readRaw(INT32_SIZE))
	if err != nil {
		return nil, err
	}

	if length != 0 && length < blockSize {
		return nil, errors.New("Buffer Overflow")
	}

	if length != 0 {
		buffer := df.readRaw(blockSize)
		return buffer, nil
	}
	return nil, nil
}

func (df *DemoFile) ReadCmdHeader() (uint8, int32, uint8) {
	if len(df.FileBuffer) == 0 {
		return DEM_STOP, -1, 0
	}
	cmdBytes := df.readRaw(2)
	cmd, err := util.ByteSliceToUInt8(cmdBytes)
	if cmd <= 0 {
		fmt.Println("Missing end tag in demo file: ", cmd)
		cmd = DEM_STOP
		return cmd, -1, 0
	}
	tickBytes := df.readRaw(4)
	tick, err := util.ByteSliceToInt32(tickBytes)
	if err != nil {
		panic(err)
	}

	playerBytes := df.readRaw(2)
	playerSlot, err := util.ByteSliceToUInt8(playerBytes)
	if err != nil {
		panic(err)
	}

	return cmd, tick, playerSlot
}

func (df *DemoFile) ReadUserCmd(size int32) ([]byte, int32, error) {
	if len(df.FileBuffer) == 0 {
		return nil, 0, nil
	}

	outgoingSeq, err := util.ByteSliceToInt32(df.readRaw(4))
	if err != nil {
		return nil, 0, err
	}
	data, err := df.ReadRawData(size)
	return data, outgoingSeq, err
}

func (df *DemoFile) ReadCmdInfo() *demoinfo {
	if len(df.FileBuffer) == 0 {
		return nil
	}
	var demInfo demoinfo
	size := unsafe.Sizeof(demInfo)
	data := df.readRaw(int32(size))
	dem := ((*demoinfo)(unsafe.Pointer(&data[0])))
	return dem
}

func (df *DemoFile) ReadSequenceInfo() (int32, int32, error) {
	if len(df.FileBuffer) == 0 {
		return 0, 0, nil
	}

	seqNrIn, err := util.ByteSliceToInt32(df.readRaw(4))
	if err != nil {
		return 0, 0, err
	}
	seqNrOut, err := util.ByteSliceToInt32(df.readRaw(4))
	if err != nil {
		return 0, 0, err
	}
	return seqNrIn, seqNrOut, nil
}
