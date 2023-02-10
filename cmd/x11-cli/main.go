package main

import (
	"fmt"
	"time"

	"github.com/wade-welles/streamkit/x11"
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
		Delay: 8 * time.Second,
	}

	x11App.X11.InitActiveWindow()

	// TODO: Probably want to load some settings from a YAML config to make things
	// easier

	fmt.Printf("x11App:%v \n", x11App)

	tick := time.Tick(x11App.Delay)
	for {
		select {
		case <-tick:
			if x11App.X11.HasActiveWindowChanged() {
				fmt.Printf("HasActiveWindowChanged() => true\n")
			}
		default:
			fmt.Printf("tick..\n")
		}
	}

}
