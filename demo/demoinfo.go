package demo

const (
	MAX_SPLITSCREEN_CLIENTS = 2
)

type split struct {
	Flags int32

	// original origin/viewangles
	viewOrigin      vector
	viewAngles      qAngle
	localViewAngles qAngle

	// Resampled origin/viewangles
	viewOrigin2      vector
	viewAngles2      qAngle
	localViewAngles2 qAngle
}

type demoinfo struct {
	Splits [MAX_SPLITSCREEN_CLIENTS]split
}
