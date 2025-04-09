package backend

type KeyboardHalfRole int

const (
	Central KeyboardHalfRole = iota
	Peripheral
)

func (k KeyboardHalfRole) Toggle() KeyboardHalfRole {
	switch k {
	case Central:
		return Peripheral
	case Peripheral:
		return Central
	default:
		return k
	}
}

func (k KeyboardHalfRole) String() string {
	switch k {
	case Central:
		return "Central"
	case Peripheral:
		return "Peripheral"
	default:
		return "Unknown"
	}
}
