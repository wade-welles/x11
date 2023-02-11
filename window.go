package x11

import (
	"strings"
)

type Windows []*Window

// TODO: Maybe desktop, always ontop, always on desktop, etc, definitely should
// include order, and desktop, and other facts.
// We will absolutely need a way to run javascript on the browser window page so
// we can set it to 160p preferably or eventually randomize to 160p to 320
type Window struct {
	Name    string
	Type    WindowType
	Focused bool // aka Active
}

var UndefinedWindow = Window{Name: "undefined", Type: UndefinedType}

type WindowType uint8 // 0..255

const (
	UndefinedType WindowType = iota
	Terminal
	Browser
	Other
)

func (wt WindowType) String() string {
	switch wt {
	case Terminal:
		return "terminal"
	case Browser:
		return "browser"
	case Other:
		return "other"
	default: // UndefinedType
		return "undefined"
	}
}

func MarshalWindowType(wt string) WindowType {
	switch strings.ToLower(wt) {
	case Terminal.String():
		return Terminal
	case Browser.String():
		return Browser
	case Other.String():
		return Other
	default:
		return UndefinedType
	}
}
