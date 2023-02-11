package x11

import (
	"fmt"
	"strings"
	"time"

	x11 "github.com/linuxdeepin/go-x11-client"
	ewmh "github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
)

// TODO: Our system doesn't work when the windows have the same name, which is a
// clear issue. We need a way to distinguish windows better than the name.

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

// //////////////////////////////////////////////////////////////////////////////
type X11 struct {
	Client  *x11.Conn
	Windows []Window

	// TODO: Maybe just cache the active window name so we do simple name
	// comparison
	ActiveWindowName      string
	ActiveWindowChangedAt time.Time
}

func ConnectToX11() *x11.Conn {
	client, err := x11.NewConn()
	if err != nil {
		panic(err)
	}
	return client
}

func (x *X11) HasActiveWindowChanged() bool {
	return !(x.ActiveWindowName == x.ActiveWindow().Name)
}

func (x *X11) ActiveWindow() Window {
	windowName, err := ewmh.GetActiveWindow(x.Client).Reply(x.Client)
	if err != nil {
		fmt.Println("error(ewmh.GetActiveWindow(x.Client)...):", err)
		return UndefinedWindow
	}

	activeWindowName, err := ewmh.GetWMName(x.Client, windowName).Reply(x.Client)
	if err != nil {
		fmt.Println("error(ewmh.GetWMName(x.Client, windowName)...):", err)
		return UndefinedWindow
	}

	activeWindow := Window{
		Name: activeWindowName,
	}

	// TODO: Switch case to determine the window type, this will be useful for
	// simplifying automation. Needs to also detect Tor/Firefox/etc
	//   * Switch case so we can load Browser if chromium/firefox/tor etc
	switch {
	case strings.HasSuffix(strings.ToLower(activeWindowName), "chromium"):
		activeWindow.Type = Browser
	case strings.HasSuffix(strings.ToLower(activeWindowName), "user@host"):
		activeWindow.Type = Terminal
	default:
		activeWindow.Type = UndefinedType
	}

	return activeWindow
}

func (x *X11) InitActiveWindow() Window {
	activeWindow := x.ActiveWindow()
	x.ActiveWindowName = activeWindow.Name
	x.ActiveWindowChangedAt = time.Now()
	return activeWindow
}

func (x *X11) CacheActiveWindow() Window {
	activeWindow := x.ActiveWindow()
	x.ActiveWindowName = x.ActiveWindow().Name
	x.ActiveWindowChangedAt = time.Now()
	return activeWindow
}

func (x *X11) IsActiveWindowType(windowType WindowType) bool {
	return x.ActiveWindow().Type == windowType
}

func (x *X11) IsActiveWindow(searchText string) bool {
	return strings.HasSuffix(strings.ToLower(x.ActiveWindowName), searchText)
}
