package util

import (
	"bytes"
	"encoding/binary"
	"math"
)

func ByteSliceToInt32(data []byte) (int32, error) {
	var result int32
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.LittleEndian, &result)
	return result, err
}

func ByteSliceToFloat32(data []byte) float32 {
	bits := binary.LittleEndian.Uint32(data)
	result := math.Float32frombits(bits)
	return result
}
