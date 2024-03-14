package main

import (
	"github.com/MarinX/keylogger"
	"github.com/go-vgo/robotgo"
	_ "github.com/go-vgo/robotgo/base"
	_ "github.com/go-vgo/robotgo/key"
	_ "github.com/go-vgo/robotgo/screen"
	_ "github.com/go-vgo/robotgo/window"
	"github.com/robotn/xgb/xproto"
	"github.com/robotn/xgbutil"
	"github.com/robotn/xgbutil/ewmh"
	"log"
)

var xu *xgbutil.XUtil

func main() {
	var err error
	xu, err = xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting multi-monitor assistant")
	state := NewState()

	presets := []CombinationPreset{
		{
			Codes: []uint16{125, 103},
			Name:  "Win + UP",
			Action: func() error {
				return MoveResize(-2160, -32, 2160, 1920)
			},
		},
		{
			Codes: []uint16{125, 108},
			Name:  "Win + DOWN",
			Action: func() error {
				return MoveResize(-2160, 1920, 2160, 1920)
			},
		},
		{
			Codes: []uint16{125, 29},
			Name:  "Win + L_CTRL",
			Action: func() error {
				pid := robotgo.GetPid()

				printCurrentGeometry(pid)

				return nil
			},
		},
	}

	InitState(state)
	state.CurrentPressedKeys.AddPresets(presets)

	PrintState(state)

	k, err := keylogger.New(keyboard)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(k *keylogger.KeyLogger) {
		err := k.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(k)

	events := k.Read()

	// range of events
	for e := range events {
		switch e.Type {
		// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
		// check the input_event.go for more events

		case keylogger.EvKey:
			// if the state of key is pressed
			if e.KeyPress() {
				state.CurrentPressedKeys.SetPressed(e.Code)
				//log.Println("[event] press key:", e.KeyString(), "key code:", e.Code)
				//log.Printf("state: %+v\n", *state.CurrentPressedKeys)
			}

			// if the state of key is released
			if e.KeyRelease() {
				state.CurrentPressedKeys.SetReleased(e.Code)

				//log.Println("[event] release key:", e.KeyString())
				//log.Printf("state: %+v\n", *state.CurrentPressedKeys)
				//
				//title := robotgo.GetTitle()
				//fmt.Println("title:", title)
				//
				//pid := robotgo.GetPid()
				//fmt.Println("pid:", pid)
				//
				//x, y, w, h := robotgo.GetClient(pid)
				//fmt.Printf("x: %d\n", x)
				//fmt.Printf("y: %d\n", y)
				//fmt.Printf("w: %d\n", w)
				//fmt.Printf("h: %d\n", h)
				//
				//// C key
				//if e.Code == 46 {
				//	fmt.Println("move window:", title)
				//	err := MoveWindow(pid, x+10, y+10)
				//	if err != nil {
				//		fmt.Println("move window err:", err)
				//	}
				//}
				//
				//// V key
				//if e.Code == 47 {
				//	fmt.Println("resize window:", title)
				//	err := ResizeWindow(pid, w+10, h+10)
				//	if err != nil {
				//		fmt.Println("resize window err:", err)
				//	}
				//}
				//
				//// B key
				//if e.Code == 48 {
				//
				//}
			}

			break
		}
	}
}

func MoveResize(x, y, w, h int) error {
	pid := robotgo.GetPid()

	log.Println("pid", pid)
	printCurrentGeometry(pid)

	xid := GetActiveWindowXid()

	log.Println("resize and move")
	err := ewmh.MoveresizeWindow(xu, xid, x, y, w, h)
	if err != nil {
		return err
	}

	printCurrentGeometry(pid)

	return nil
}

func printCurrentGeometry(pid int) {
	cx, cy, cw, ch := robotgo.GetClient(pid)

	log.Printf("x: %d\ty:%d\tw:%d\th:%d\n", cx, cy, cw, ch)
}

func GetActiveWindowXid() xproto.Window {
	return xproto.Window(robotgo.GetActive().XWin)
}

//
//func ResizeWindow(xid xproto.Window, width, height int) error {
//
//	//wmPid, err := ewmh.WmPidGet(xu, xid)
//	//log.Println("XWin:", xid, "wmPid:", wmPid, "err:", err)
//
//	err := ewmh.ResizeWindow(xu, xid, width, height)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

//func MoveWindow(xid xproto.Window, x, y int) error {
//
//	err = ewmh.MoveWindow(xu, xid, x, y)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

//
//func test1() {
//	// Get a list of all client ids.
//	clientids, err := ewmh.ClientListGet(xu)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Iterate through each client, find its name and find its size.
//	for _, clientid := range clientids {
//		fmt.Printf("window id: %d\n", clientid)
//
//		name, err := ewmh.WmNameGet(xu, clientid)
//
//		// If there was a problem getting _NET_WM_NAME or if its empty,
//		// try the old-school version.
//		if err != nil || len(name) == 0 {
//			name, err = icccm.WmNameGet(xu, clientid)
//
//			// If we still can't find anything, give up.
//			if err != nil || len(name) == 0 {
//				name = "N/A"
//			}
//		}
//
//		// Now find the geometry, including decorations, of the client window.
//		// Note that DecorGeometry actually traverses the window tree by
//		// issuing QueryTree requests until a top-level window (i.e., its
//		// parent is the root window) is found. The geometry of *that* window
//		// is then returned.
//		dgeom, err := xwindow.New(xu, clientid).DecorGeometry()
//		if err != nil {
//			log.Printf("Could not get geometry for %s (0x%X) because: %s",
//				name, clientid, err)
//			continue
//		}
//
//		fmt.Printf("%s (0x%x)\n", name, clientid)
//		fmt.Printf("\tGeometry: %s\n", dgeom)
//	}
//}
