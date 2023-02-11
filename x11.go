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

// //////////////////////////////////////////////////////////////////////////////
type X11 struct {
	Client  *x11.Conn
	Windows []Window

	// TODO: Maybe just cache the active window name so we do simple name
	// comparison, but this leads to a bug where two windows with the same name
	// are considered the name window
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

// TODO: If we can move some of these to be methods of Window struct, it would
// be better organized but there will be obvious limitations we have to work
// through
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
