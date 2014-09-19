package main

// #cgo CFLAGS: -g -Wall
// #cgo pkg-config: x11 xi
// #include "xinput.h"
import "C"

func init() {
	go C.start_listen()
}

//export go_handle_raw_event
func go_handle_raw_event(evt_type int, detail int32, x, y, mask int32) {
	switch evt_type {
	case C.XI_RawKeyPress:
		GetManager().handleKeyboardEvent(detail, true, x, y)
	case C.XI_RawKeyRelease:
		GetManager().handleKeyboardEvent(detail, false, x, y)
	case C.XI_RawTouchBegin:
		GetManager().handleButtonEvent(1, true, x, y)
	case C.XI_RawButtonPress:
		GetManager().handleButtonEvent(detail, true, x, y)
	case C.XI_RawTouchEnd:
		GetManager().handleButtonEvent(1, false, x, y)
	case C.XI_RawButtonRelease:
		GetManager().handleButtonEvent(detail, false, x, y)

	case C.XI_RawTouchUpdate:
		GetManager().handleCursorEvent(x, y, false)
	case C.XI_RawMotion:
		/**
		* mouse left press: mask = 256
		* mouse right press: mask = 512
		* mouse middle press: mask = 1024
		**/
		if mask >= 256 {
			GetManager().handleCursorEvent(x, y, true)
		} else {
			GetManager().handleCursorEvent(x, y, false)
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
