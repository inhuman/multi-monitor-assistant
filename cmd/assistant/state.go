package main

import (
	"fmt"
	"github.com/MarinX/keylogger"
	"github.com/go-vgo/robotgo"
	"github.com/sanity-io/litter"
	"strings"
)

const keyboard = "/dev/input/by-path/pci-0000:00:14.0-usb-0:3:1.0-event-kbd"

type State struct {
	Displays           []DisplayInfo
	Keyboard           string
	CurrentPressedKeys *CurrentCombinationKeys
}

func NewState() *State {
	return &State{
		CurrentPressedKeys: NewCurrentCombinationKeys(),
	}
}

type DisplayInfo struct {
	DisplayNum int
	X          int
	Y          int
	Width      int
	Height     int
	IsPrimary  bool
}

func InitState(state *State) {
	state.Keyboard = keyboard

	displaysCount := robotgo.DisplaysNum()

	for i := 0; i < displaysCount; i++ {
		x, y, w, h := robotgo.GetDisplayBounds(i)
		state.Displays = append(state.Displays, DisplayInfo{
			DisplayNum: i,
			X:          x,
			Y:          y,
			Width:      w,
			Height:     h,
			IsPrimary:  i == robotgo.GetMainId(),
		})
	}
}

func PrintState(state *State) {
	fmt.Println(litter.Sdump(state))
}

type KeyState uint8

const (
	KeyPressed  KeyState = 1
	KeyReleased KeyState = 2
)

type CurrentCombinationKeys struct {
	keys map[uint16]KeyState
}

func NewCurrentCombinationKeys() *CurrentCombinationKeys {
	return &CurrentCombinationKeys{
		keys: make(map[uint16]KeyState),
	}
}

func (cck *CurrentCombinationKeys) SetPressed(keyCode uint16) {
	//fmt.Println("pressed", keyToString(keyCode))
	cck.keys[keyCode] = KeyPressed
}

func (cck *CurrentCombinationKeys) SetReleased(keyCode uint16) {
	//fmt.Println("released", keyToString(keyCode))
	cck.keys[keyCode] = KeyReleased

	if isControlKey(keyCode) {
		//fmt.Println("detected control key released")
		cck.PrintCurrentCombination()
		cck.Flush()
	}
}

func (cck *CurrentCombinationKeys) Flush() {
	cck.keys = make(map[uint16]KeyState)
}

func (cck *CurrentCombinationKeys) PrintCurrentCombination() {
	combStr := strings.Builder{}

	for keyCode := range cck.keys {
		combStr.WriteString(fmt.Sprintf("%s ", keyToString(keyCode)))
	}

	fmt.Println("combination", combStr.String())
}

var controlKeys = [...]uint16{
	29,  //  "L_CTRL"
	42,  //  "L_SHIFT"
	54,  //  "R_SHIFT"
	56,  //  "L_ALT"
	97,  //  "R_CTRL"
	100, //  "R_ALT",
	125, //  "Super"
}

func isControlKey(keyCode uint16) bool {
	for _, code := range controlKeys {
		if keyCode == code {
			return true
		}
	}

	return false
}

func keyToString(keyCode uint16) string {
	return (&keylogger.InputEvent{Code: keyCode}).KeyString()
}
