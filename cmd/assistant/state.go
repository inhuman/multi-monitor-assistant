package main

import (
	"fmt"
	"github.com/MarinX/keylogger"
	"github.com/go-vgo/robotgo"
	"github.com/sanity-io/litter"
	"log"
	"slices"
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

type KeysState map[uint16]KeyState

func (ks KeysState) slice() []uint16 {
	var sl []uint16
	for code := range ks {
		sl = append(sl, code)
	}

	return sl
}

type CurrentCombinationKeys struct {
	combinationStarted bool
	keys               KeysState
	presets            []CombinationPreset
}

func NewCurrentCombinationKeys() *CurrentCombinationKeys {
	return &CurrentCombinationKeys{
		keys: make(map[uint16]KeyState),
	}
}

func (cck *CurrentCombinationKeys) AddPresets(presets []CombinationPreset) {
	cck.presets = presets
}

func (cck *CurrentCombinationKeys) SetPressed(keyCode uint16) {
	if isControlKey(keyCode) {
		cck.combinationStarted = true
	}

	if cck.combinationStarted {
		cck.keys[keyCode] = KeyPressed
	}
}

func (cck *CurrentCombinationKeys) SetReleased(keyCode uint16) {
	cck.keys[keyCode] = KeyReleased

	preset, isPreset := cck.IsPreset()
	if isPreset {
		log.Println("fired:", preset.Name)
		err := preset.Action()
		if err != nil {
			log.Printf("error doing action '%s': %s", preset.Name, err.Error())
		}

		cck.Flush()
	}

	// TODO: если в комбинации больше двух кнопок - третья сбрасывается

	cck.Flush()
}

func (cck *CurrentCombinationKeys) IsPreset() (CombinationPreset, bool) {
	for _, preset := range cck.presets {
		if preset.IsEqual(cck.keys) {
			return preset, true
		}
	}

	return CombinationPreset{}, false
}

func (cck *CurrentCombinationKeys) Flush() {
	log.Printf("state: %+v\n", *cck)

	for key, status := range cck.keys {
		if status == KeyReleased {
			delete(cck.keys, key)
		}
	}

	cck.combinationStarted = false
}

//func (cck *CurrentCombinationKeys) PrintCurrentCombination() {
//	combStr := strings.Builder{}
//
//	for keyCode := range cck.keys {
//		combStr.WriteString(fmt.Sprintf("%s ", keyToString(keyCode)))
//	}
//
//	fmt.Println("combination", combStr.String())
//}

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

func (cck *CurrentCombinationKeys) isLastControlKeyReleased() bool {
	for key, status := range cck.keys {
		if isControlKey(key) && status == KeyPressed {
			return false
		}
	}

	return true
}

func keyToString(keyCode uint16) string {
	return (&keylogger.InputEvent{Code: keyCode}).KeyString()
}

type CombinationPreset struct {
	Codes  []uint16
	Action ShortcutAction
	Name   string
}

type ShortcutAction func() error

func (cp CombinationPreset) IsEqual(currentCombination KeysState) bool {
	slices.Sort(cp.Codes)

	codes := currentCombination.slice()
	slices.Sort(codes)

	return equal(codes, cp.Codes)
}

func equal(a, b []uint16) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
