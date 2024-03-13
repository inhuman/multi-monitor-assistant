package main

import (
	"fmt"
	"github.com/MarinX/keylogger"
	"github.com/go-vgo/robotgo"
	_ "github.com/go-vgo/robotgo/base"
	_ "github.com/go-vgo/robotgo/key"
	_ "github.com/go-vgo/robotgo/mouse"
	_ "github.com/go-vgo/robotgo/screen"
	_ "github.com/go-vgo/robotgo/window"
	"github.com/robotn/xgbutil"
	"github.com/robotn/xgbutil/ewmh"
	"github.com/robotn/xgbutil/icccm"
	"github.com/robotn/xgbutil/xwindow"
	"log"
)

const keyboard = "/dev/input/by-path/pci-0000:00:14.0-usb-0:3:1.0-event-kbd"

var xu *xgbutil.XUtil

func main() {
	var err error
	xu, err = xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting multi-monitor assistant")
	displaysCount := robotgo.DisplaysNum()

	fmt.Println("Displays count:", displaysCount)
	fmt.Println("Main display id:", robotgo.GetMainId())

	for i := 0; i < displaysCount; i++ {
		x, y, w, h := robotgo.GetDisplayBounds(i)

		fmt.Printf("\nDisplay %d bounds:\n", i)
		fmt.Printf("x: %d\n", x)
		fmt.Printf("y: %d\n", y)
		fmt.Printf("w: %d\n", w)
		fmt.Printf("h: %d\n", h)
	}

	state := NewState()

	InitState(state)

	log.Println("Found a keyboard at", keyboard)
	// init keylogger with keyboard
	k, err := keylogger.New(keyboard)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer k.Close()

	events := k.Read()

	// range of events
	for e := range events {
		switch e.Type {
		// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
		// check the input_event.go for more events
		case keylogger.EvKey:

			// if the state of key is pressed
			if e.KeyPress() {
				log.Println("[event] press key:", e.KeyString(), "key code:", e.Code)
			}

			// if the state of key is released
			if e.KeyRelease() {
				log.Println("[event] release key:", e.KeyString())

				title := robotgo.GetTitle()
				fmt.Println("title:", title)

				pid := robotgo.GetPid()
				fmt.Println("pid:", pid)

				x, y, w, h := robotgo.GetClient(pid)
				fmt.Printf("x: %d\n", x)
				fmt.Printf("y: %d\n", y)
				fmt.Printf("w: %d\n", w)
				fmt.Printf("h: %d\n", h)

				// C key
				if e.Code == 46 {
					fmt.Println("move window:", title)
					err := MoveWindow(pid, x+10, y+10)
					if err != nil {
						fmt.Println("move window err:", err)
					}
				}

				// V key
				if e.Code == 47 {
					fmt.Println("resize window:", title)
					err := ResizeWindow(pid, w+10, h+10)
					if err != nil {
						fmt.Println("resize window err:", err)
					}
				}

				// B key
				if e.Code == 48 {

				}
			}

			break
		}
	}
}

//func MaximizeWindow(pid int) error {
//	xid, err := robotgo.GetXid(xu, pid)
//	if err != nil {
//		return err
//	}
//
//}

func ResizeWindow(pid int, width, height int) error {
	xid, err := robotgo.GetXid(xu, pid)
	if err != nil {
		return err
	}

	err = ewmh.ResizeWindow(xu, xid, width, height)
	if err != nil {
		return err
	}

	return nil
}

func MoveWindow(pid int, x, y int) error {
	xid, err := robotgo.GetXid(xu, pid)
	if err != nil {
		return err
	}

	err = ewmh.MoveWindow(xu, xid, x, y)
	if err != nil {
		return err
	}

	return nil
}

func test1() {
	// Get a list of all client ids.
	clientids, err := ewmh.ClientListGet(xu)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through each client, find its name and find its size.
	for _, clientid := range clientids {
		fmt.Printf("window id: %d\n", clientid)

		name, err := ewmh.WmNameGet(xu, clientid)

		// If there was a problem getting _NET_WM_NAME or if its empty,
		// try the old-school version.
		if err != nil || len(name) == 0 {
			name, err = icccm.WmNameGet(xu, clientid)

			// If we still can't find anything, give up.
			if err != nil || len(name) == 0 {
				name = "N/A"
			}
		}

		// Now find the geometry, including decorations, of the client window.
		// Note that DecorGeometry actually traverses the window tree by
		// issuing QueryTree requests until a top-level window (i.e., its
		// parent is the root window) is found. The geometry of *that* window
		// is then returned.
		dgeom, err := xwindow.New(xu, clientid).DecorGeometry()
		if err != nil {
			log.Printf("Could not get geometry for %s (0x%X) because: %s",
				name, clientid, err)
			continue
		}

		fmt.Printf("%s (0x%x)\n", name, clientid)
		fmt.Printf("\tGeometry: %s\n", dgeom)
	}
}
