package dxinput

import (
	"pkg.deepin.io/dde/api/dxinput/utils"
)

const (
	// detail see: man libinput
	// 1 boolean value (8 bit, 0 or 1).
	libinputPropTapEnabled = "libinput Tapping Enabled"
	// 1 boolean value (8 bit, 0 or 1).
	libinputPropDragEnabled = "libinput Tapping Drag Enabled"
	// 1 boolean value (8 bit, 0 or 1).
	libinputPropDragLockEnabled = "libinput Tapping Drag Lock Enabled"
	// Either one 8-bit value specifying the meta drag lock button, or a list of button pairs.
	// See section Button Drag Lock for details.
	libinputPropDragLockButtons = "libinput Tapping Drag Lock Buttons"
	// 1 32-bit float value
	// Sets the pointer acceleration speed within the range [-1, 1]
	libinputPropAccelSpeed = "libinput Accel Speed"
	// 1 boolean value (8 bit, 0 or 1).
	libinputPropNaturalEnabled = "libinput Natural Scrolling Enabled"
	// 1 boolean value (8 bit, 0 or 1).
	libinputPropLeftHandedEnabled = "libinput Left Handed Enabled"
	// 3 boolean values (8 bit, 0 or 1), in order "two-finger", "edge", "button".
	// Indicates which scroll method is currently enabled on this device.
	libinputPropScrollMethodsEnabled = "libinput Scroll Method Enabled"
	// 3 boolean values (8 bit, 0 or 1), in order "two-finger", "edge", "button".
	// Indicates which scroll methods are available on this device.
	libinputPropScrollMethodsAvailable = "libinput Scroll Methods Available"
	// 2 boolean values (8 bit, 0 or 1), in order "buttonareas", "clickfinger"
	// Indicates which click methods are available on this device.
	// TODO
	libinputPropClickMethodEnabled = "libinput Click Method Enabled"
	// 2 boolean values (8 bit, 0 or 1), in order "buttonareas", "clickfinger"
	// Indicates which click methods are available on this device.
	// TODO
	libinputPropClickMethodAvailable = "libinput Click Methods Available"
	// 1 32-bit value
	libinputPropButtonScrollingButton = "libinput Button Scrolling Button"
	// 1 boolean value (8 bit, 0 or 1).
	libinputPropMiddleEmulationEnabled = "libinput Middle Emulation Enabled"
	// 1 boolean value (8 bit, 0 or 1).
	libinputPropHorizScrollEnabled = "libinput Horizontal Scroll Enabled"
)

// scroll methods: two-finger, edge, button. button only for trackpoint
func libinputCanScroll(id int32) (bool, bool, bool) {
	available, err := getInt8Prop(id, libinputPropScrollMethodsAvailable, 3)
	if err != nil {
		return false, false, false
	}

	values, err := getInt8Prop(id, libinputPropScrollMethodsEnabled, 3)
	if err != nil {
		return false, false, false
	}
	return (available[0] == 1) && (values[0] == 1),
		(available[1] == 1) && (values[1] == 1),
		(available[2] == 1) && (values[2] == 1)
}

func libinputEnableScrollTwoFinger(id int32, enabled bool) error {
	values, err := getInt8Prop(id, libinputPropScrollMethodsEnabled, 3)
	if err != nil {
		return err
	}

	if enabled {
		if values[0] == 1 {
			return nil
		}
		// These scroll methods are mutually exclusive.
		values[0] = 1
		values[1] = 0
		values[2] = 0
	} else {
		if values[0] == 0 {
			return nil
		}
		values[0] = 0
	}
	return utils.SetInt8Prop(id, libinputPropScrollMethodsEnabled, values)
}

func libinputEnableScrollEdge(id int32, enabled bool) error {
	values, err := getInt8Prop(id, libinputPropScrollMethodsEnabled, 3)
	if err != nil {
		return err
	}

	if enabled {
		if values[1] == 1 {
			return nil
		}
		values[1] = 1
		values[0] = 0
		values[2] = 0
	} else {
		if values[1] == 0 {
			return nil
		}
		values[1] = 0
	}
	return utils.SetInt8Prop(id, libinputPropScrollMethodsEnabled, values)
}

func libinputEnableScrollButton(id int32, enabled bool) error {
	values, err := getInt8Prop(id, libinputPropScrollMethodsEnabled, 3)
	if err != nil {
		return err
	}

	if enabled {
		if values[2] == 1 {
			return nil
		}
		values[2] = 1
		values[0] = 0
		values[1] = 0
	} else {
		if values[2] == 0 {
			return nil
		}
		values[2] = 0
	}
	return utils.SetInt8Prop(id, libinputPropScrollMethodsEnabled, values)
}

func libinputEnableScroll(id int32, twoFinger, edge, button bool) error {
	available, err := getInt8Prop(id, libinputPropScrollMethodsAvailable, 3)
	if err != nil {
		return err
	}

	var values = []int8{1, 1, 1}

	if available[0] == 0 {
		twoFinger = false
	}
	if available[1] == 0 {
		edge = false
	}
	if available[2] == 0 {
		button = false
	}

	if !twoFinger {
		values[0] = 0
	}
	if !edge {
		values[1] = 0
	}
	if !button {
		values[2] = 0
	}

	return utils.SetInt8Prop(id, libinputPropScrollMethodsEnabled, values)
}

func libinputGetAccel(id int32) (float32, error) {
	values, err := getFloat32Prop(id, libinputPropAccelSpeed, 1)
	if err != nil {
		return 0, err
	}
	return values[0], nil
}

func libinputSetAccel(id int32, accel float32) error {
	// range [-1 ~ 1]
	if accel > 1 {
		accel = 1
	} else if accel < -1 {
		accel = -1
	}
	v, _ := libinputGetAccel(id)
	if v == accel {
		return nil
	}
	return utils.SetFloat32Prop(id, libinputPropAccelSpeed, []float32{accel})
}

func libinputGetScrollButton(id int32) (int32, error) {
	values, err := getInt32Prop(id, libinputPropButtonScrollingButton, 1)
	if err != nil {
		return -1, err
	}
	return values[0], nil
}

func libinputSetScrollButton(id, btn int32) error {
	return utils.SetInt32Prop(id, libinputPropButtonScrollingButton, []int32{btn})
}

func libinputInt8PropCan(id int32, prop string) bool {
	values, err := getInt8Prop(id, prop, 1)
	if err != nil {
		return false
	}
	return values[0] == 1
}

func libinputInt8PropSet(id int32, prop string, enabled bool) error {
	var v int8 = 1
	if !enabled {
		v = 0
	}
	return utils.SetInt8Prop(id, prop, []int8{v})
}
