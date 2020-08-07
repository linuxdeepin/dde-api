/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package dxinput

import (
	"pkg.deepin.io/dde/api/dxinput/utils"
	"errors"
)

const (
	// detail see: man libinput
	// 1 boolean value (8 bit, 0 or 1).
	libinputPropTapEnabled = "libinput Tapping Enabled"
	// 2 boolean values (8 bit, 0 or 1), in order "adaptive", "flat".
	// Indicates which acceleration profiles are available on this device.
	libinputPropAccelProfileAvaliable = "libinput Accel Profiles Available"
	// 2 boolean values (8 bit, 0 or 1), in order "adaptive", "flat".
	// Indicates which acceleration profile is currently enabled on this device.
	libinputPropAccelProfileEnabled = "libinput Accel Profile Enabled"
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
	// 1 32-bit value
	libinputPropButtonScrollingButton = "libinput Button Scrolling Button"
	// 1 boolean value (8 bit, 0 or 1).
	libinputPropMiddleEmulationEnabled = "libinput Middle Emulation Enabled"
	// 1 boolean value (8 bit, 0 or 1).
	libinputPropHorizScrollEnabled = "libinput Horizontal Scroll Enabled"
	// 1 boolean value (8 bit, 0 or 1).
	libinputPropDisableWhileTyping = "libinput Disable While Typing Enabled"
)

// for mouse: check if both "adaptive" and "flat" profile are avaliable
func libinputIsBothAccelProfileAvaliable(id int32) bool {
	available, err := getInt8Prop(id, libinputPropAccelProfileAvaliable, 2)
	if err != nil {
		return false
	}

	return (available[0] == 1) && (available[1] == 1)
}

// for mouse: get enabled accel profile, in order "adaptive", "flat".
func libinputGetAccelProfile(id int32) (bool, bool) {
	enabled, err := getInt8Prop(id, libinputPropAccelProfileEnabled, 2)
	if err != nil {
		return false, false
	}

	return enabled[0] == 1, enabled[1] == 1
}

// for mouse: set enabled accel profile, in order "adaptive", "flat".
func libinputSetAccelProfile(id int32, useAdaptiveProfile bool) error {
	if !libinputIsBothAccelProfileAvaliable(id) {
		return errors.New("dde-api: device doesn't support both accel profile")
	}

	prop, err := getInt8Prop(id, libinputPropAccelProfileEnabled, 2)
	if err != nil {
		return err
	}

	if useAdaptiveProfile {
		prop[0] = 1
		prop[1] = 0
	} else {
		prop[0] = 0
		prop[1] = 1
	}

	return utils.SetInt8Prop(id, libinputPropAccelProfileEnabled, prop)
}

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
