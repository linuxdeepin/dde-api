package main

// #cgo CFLAGS: -g -Wall
// #cgo pkg-config: x11 xi
// #include "record.h"
import "C"

func init() {
	go C.start_listen()
}

//export go_handle_device_event
func go_handle_device_event(evt_type int, event *C.XIDeviceEvent) {
	switch event.evtype {
	case C.XI_KeyPress:
		GetManager().handleKeyboardEvent(int32(event.detail), true, int32(event.root_x), int32(event.root_y))
	case C.XI_KeyRelease:
		GetManager().handleKeyboardEvent(int32(event.detail), false, int32(event.root_x), int32(event.root_y))
	case C.XI_TouchBegin:
		GetManager().handleButtonEvent(1, true, int32(event.root_x), int32(event.root_y))
	case C.XI_ButtonPress:
		GetManager().handleButtonEvent(int32(event.detail), true, int32(event.root_x), int32(event.root_y))
	case C.XI_TouchEnd:
		GetManager().handleButtonEvent(1, false, int32(event.root_x), int32(event.root_y))
	case C.XI_ButtonRelease:
		GetManager().handleButtonEvent(int32(event.detail), false, int32(event.root_x), int32(event.root_y))
	case C.XI_Motion:
		if len(getButtonState(event)) > 0 {
			GetManager().handleMotionEvent(int32(event.root_x), int32(event.root_y), true)
		} else {
			GetManager().handleMotionEvent(int32(event.root_x), int32(event.root_y), false)
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
