package main

import (
	"fmt"
	"time"

	"github.com/wade-welles/x11"
)

type Application struct {
	X11   *x11.X11
	Delay time.Duration
	Paths map[PathType]Path
}

// TODO: Initialize some Paths for our basic application, which would include,
// config (~/.config/$APP_NAME/*) => so config file could be the existance of a
// .conf file in the folder.
// Local is the local cache stored in (~/.local/share/$APP_NAME/*)
// System is only necessary if the application is capable or even desirable to
// ever be ran as root.
// Log is the path to the logs, with the ability to output to a variety of log
// files (see how CLI framework works)
type PathType int

const (
	Config PathType = iota
	Local
	System
	Log
)

type Path struct {
	Type PathType
}

func main() {
	fmt.Printf("x11-cli\n")
	fmt.Printf("=======\n")
	fmt.Printf("Currently not yet implemented.\n")

	x11App := Application{
		X11: &x11.X11{
			Client: x11.ConnectToX11(),
		},
		Delay: 2 * time.Second,
	}

	x11App.X11.InitActiveWindow()

	// TODO: Probably want to load some settings from a YAML config to make things
	// easier

	fmt.Printf("x11App:\n")

	tick := time.Tick(x11App.Delay)
	for {
		select {
		case <-tick:
			if x11App.X11.HasActiveWindowChanged() {
				fmt.Printf("HasActiveWindowChanged(): true\n")

				activeWindow := x11App.X11.ActiveWindow()
				fmt.Printf("  active_window_name: %s\n", activeWindow.Name)

				fmt.Printf("  x11.ActiveWindowName: %v\n", x11App.X11.ActiveWindowName)

				// NOTE: This worked to prevent it from repeating
				// HasActiveWindowChanged() over and over
				x11App.X11.CacheActiveWindow()

			} else {
				fmt.Printf("tick,...\n")
				fmt.Printf("  x11.ActiveWindowName: %v\n", x11App.X11.ActiveWindowName)
				fmt.Printf(
					"  x11.ActiveWindow().Type.String(): %v\n",
					x11App.X11.ActiveWindow().Type.String(),
				)
			}
		}
	}

}
