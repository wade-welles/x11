package x11

import (
	"fmt"
	"strings"
	"time"

	x11 "github.com/linuxdeepin/go-x11-client"
	ewmh "github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
)

type Window struct {
	Name    WindowName
	Focused bool // aka Active
}

type Windows []*Window

type WindowName uint8 // 0..255

const (
	UndefinedName WindowName = iota
	Terminal
	Browser
	Other
)

// TODO: For now we can extend it by tracking for specific types of windows say:
// FileManager, Browser, Terminal, Other
//
//	Then we store the actual title from the window on it and use it to decipher
//	it into a selction of rebuilt tools
func (wt WindowName) String() string {
	switch wt {
	case Terminal:
		return "terminal"
	case Browser:
		return "browser"
	case Other:
		return "other"
	default: // UndefinedName
		return "undefined"
	}
}

func MarshalWindowName(wt string) WindowName {
	switch strings.ToLower(wt) {
	case Terminal.String():
		return Terminal
	case Browser.String():
		return Browser
	case Other.String():
		return Other
	default: // UndefinedName
		return UndefinedName
	}
}

// //////////////////////////////////////////////////////////////////////////////
type X11 struct {
	Client *x11.Conn

	CurrentWindow WindowName

	ActiveWindowChangedAt time.Time
}

func ConnectToX11() *x11.Conn {
	client, err := x11.NewConn()
	if err != nil {
		panic(err)
	}
	return client
}

// TODO: X11 =(should return a Window() function=> WindowName type
func (x *X11) Window() WindowName {
	x.CacheActiveWindow()
	return x.CurrentWindow
}

func (x *X11) HasActiveWindowChanged() bool {
	// NOTE
	// If we record time last active window change happened, then we can limit
	// the number of times it can change within x amount of time for better
	// reliability under pressure.
	return !x.IsCurrentWindow(x.ActiveWindow())
}

func (x *X11) ActiveWindow() WindowName {
	windowName, err := ewmh.GetActiveWindow(x.Client).Reply(x.Client)
	if err != nil {
		fmt.Println("error(ewmh.GetActiveWindow(x.Client)...):", err)
		return UndefinedName
	}

	activeWindowName, err := ewmh.GetWMName(x.Client, windowName).Reply(x.Client)
	if err != nil {
		fmt.Println("error(ewmh.GetWMName(x.Client, windowName)...):", err)
		return UndefinedName
	}

	if strings.HasSuffix(strings.ToLower(activeWindowName), "chromium") {
		return Browser
	} else {
		return MarshalWindowName(activeWindowName)
	}
}

func (x *X11) ActiveWindowString() string {
	windowName, err := ewmh.GetActiveWindow(x.Client).Reply(x.Client)
	if err != nil {
		fmt.Println("error(ewmh.GetActiveWindow(x.Client)...):", err)
		return UndefinedName
	}

	activeWindowName, err := ewmh.GetWMName(x.Client, windowName).Reply(x.Client)
	if err != nil {
		fmt.Println("error(ewmh.GetWMName(x.Client, windowName)...):", err)
		return UndefinedName
	}
	return activeWindowName
}

func (x *X11) InitActiveWindow() WindowName {
	x.CurrentWindow = Primary
	x.ActiveWindowChangedAt = time.Now()
	return x.CurrentWindow
}

func (x *X11) CacheActiveWindow() WindowName {
	x.CurrentWindow = x.ActiveWindow()
	x.ActiveWindowChangedAt = time.Now()
	return x.CurrentWindow
}

// TODO: Naming issue here bigger than may originally appear
// x.IsActiveWindow(Primary)
func (x *X11) IsActiveWindow(windowName WindowName) bool {
	return x.ActiveWindow() == windowName
}

func (x *X11) IsCurrentWindow(windowName WindowName) bool {
	return x.CurrentWindow == windowName
}
