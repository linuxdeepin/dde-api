package main

// #cgo CFLAGS: -g -Wall
// #cgo pkg-config: x11 xi
// #include "xinput.h"
import "C"

func init() {
	go C.start_listen()
}

//export go_handle_raw_event
func go_handle_raw_event(evt_type int, event *C.XIRawEvent, x, y, mask int32) {
	switch event.evtype {
	case C.XI_RawKeyPress:
		GetManager().handleKeyboardEvent(int32(event.detail), true, x, y)
	case C.XI_RawKeyRelease:
		GetManager().handleKeyboardEvent(int32(event.detail), false, x, y)
	case C.XI_RawTouchBegin:
		GetManager().handleButtonEvent(1, true, x, y)
	case C.XI_RawButtonPress:
		GetManager().handleButtonEvent(int32(event.detail), true, x, y)
	case C.XI_RawTouchEnd:
		GetManager().handleButtonEvent(1, false, x, y)
	case C.XI_RawButtonRelease:
		GetManager().handleButtonEvent(int32(event.detail), false, x, y)

	case C.XI_RawTouchUpdate:
		GetManager().handleMotionEvent(x, y, false)
	case C.XI_RawMotion:
		if mask != 0 {
			GetManager().handleMotionEvent(x, y, true)
		} else {
			GetManager().handleMotionEvent(x, y, false)
		}
	}
}

func getButtonState(event *C.XIDeviceEvent) []int {
	var buttons []int
	for i := 0; i < int(event.buttons.mask_len)*8; i++ {
		if C.xi_mask_is_set(event.buttons.mask, C.char(i)) != 0 {
			buttons = append(buttons, i)
		}
	}
	return buttons
}
