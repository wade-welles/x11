package main

import (
	"fmt"
	"time"

	"github.com/wade-welles/x11"
)

type Application struct {
	X11   *x11.X11
	Delay time.Duration
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
				fmt.Printf("HasActiveWindowChanged() => true\n")

				activeWindow := x11App.X11.ActiveWindow()

				fmt.Printf("  active_window:%v\n", activeWindow)
				fmt.Printf("  active_window_string:%s\n", activeWindowString())

				fmt.Printf("  active_window_string:%s\n", x11App.X11.ActiveWindowString())

				fmt.Printf(
					"*.X11.CacheActiveWindow() => %v\n",
					x11App.X11.CacheActiveWindow(),
				)

			} else {
				fmt.Printf("tick,...")
			}
		}
	}

}
