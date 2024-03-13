package main

import (
	"github.com/go-vgo/robotgo"
)

type State struct {
	Displays []DisplayInfo
	Keyboard string
}

func NewState() *State {
	return &State{}
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
