package demo

type vector struct {
	X float32
	Y float32
	Z float32
}

func NewVector(x, y, z float32) *vector {
	return &vector{x, y, z}
}
