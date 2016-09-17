package demo

type qAngle struct {
	X float32
	Y float32
	Z float32
}

func NewQAngle(x, y, z float32) *qAngle {
	return &qAngle{x, y, z}
}
