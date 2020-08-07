/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	"fmt"

	. "pkg.deepin.io/dde/api/dxinput/common"
	"pkg.deepin.io/dde/api/dxinput/kwayland"
	"pkg.deepin.io/dde/api/dxinput/utils"
)

const (
	// only for xf86-input-synaptics
	propOff            string = "Synaptics Off"
	propScrollDistance string = "Synaptics Scrolling Distance"
	propEdgeScroll     string = "Synaptics Edge Scrolling"
	propTwoFingerScrol string = "Synaptics Two-Finger Scrolling"
	propTapAction      string = "Synaptics Tap Action"
	propPalmDetect     string = "Synaptics Palm Detection"
	propPalmDimensions string = "Synaptics Palm Dimensions"
)

type Touchpad struct {
	Id   int32
	Name string

	// -1: unknown, 0: not used, 1: used
	isLibinputUsed bool
}

/**
 * touchpad properties see:
 * http://www.x.org/archive/X11R7.5/doc/man/man4/synaptics.4.html#sect4
 *
 * Also use 'xinput list-props <id>' to list these props.
 **/
func NewTouchpad(id int32) (*Touchpad, error) {
	info := utils.ListDevice().Get(id)
	if info == nil {
		return nil, fmt.Errorf("Invalid device id: %v", id)
	}
	return NewTouchpadFromDevInfo(info)
}

func NewTouchpadFromDevInfo(dev *DeviceInfo) (*Touchpad, error) {
	if dev == nil || dev.Type != DevTypeTouchpad {
		return nil, fmt.Errorf("Not a touchpad device(%d - %s)", dev.Id, dev.Name)
	}

	return &Touchpad{
		Id:             dev.Id,
		Name:           dev.Name,
		isLibinputUsed: utils.IsPropertyExist(dev.Id, libinputPropTapEnabled),
	}, nil
}

/**
 * Property 'Synaptics Off' 8 bit, valid values (0, 1, 2):
 *	Value 0: Touchpad is enabled
 *	Value 1: Touchpad is switched off
 *	Value 2: Only tapping and scrolling is switched off
 **/
func (tpad *Touchpad) Enable(enabled bool) error {
	if globalWayland {
		return kwayland.Enable(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id), enabled)
	}

	err := enableDevice(tpad.Id, enabled)
	if err != nil {
		return err
	}

	if enabled == tpad.IsEnabled() || tpad.isLibinputUsed {
		return nil
	}

	var values []int8
	if enabled {
		values = []int8{0}
	} else {
		values = []int8{1}
	}

	return utils.SetInt8Prop(tpad.Id, propOff, values)
}

func (tpad *Touchpad) IsEnabled() bool {
	if globalWayland {
		return kwayland.CanEnabled(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id))
	}

	if !isDeviceEnabled(tpad.Id) {
		return false
	}

	if tpad.isLibinputUsed {
		return true
	}

	values, err := getInt8Prop(tpad.Id, propOff, 1)
	if err != nil {
		return false
	}

	return (values[0] == 0)
}

func (tpad *Touchpad) EnableLeftHanded(enabled bool) error {
	if enabled == tpad.CanLeftHanded() {
		return nil
	}

	if globalWayland {
		return kwayland.EnableLeftHanded(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id), enabled)
	}
	if tpad.isLibinputUsed {
		return libinputInt8PropSet(tpad.Id, libinputPropLeftHandedEnabled, enabled)
	}
	return utils.SetLeftHanded(uint32(tpad.Id), tpad.Name, enabled)
}

func (tpad *Touchpad) CanLeftHanded() bool {
	if globalWayland {
		return kwayland.CanLeftHanded(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id))
	}
	if tpad.isLibinputUsed {
		return libinputInt8PropCan(tpad.Id, libinputPropLeftHandedEnabled)
	}
	return utils.CanLeftHanded(uint32(tpad.Id), tpad.Name)
}

/**
 * Property 'Synaptics Tap Action' 8 bit,
 * up to MAX_TAP values (see synaptics.h), 0 disables an element.
 * order: RT, RB, LT, LB, F1, F2, F3.
 *	Option "RTCornerButton" "integer":
 *		Which mouse button is reported on a right top corner tap.
 *	Option "RBCornerButton" "integer":
 *		Which mouse button is reported on a right bottom corner tap.
 *	Option "LTCornerButton" "integer":
 *		Which mouse button is reported on a left top corner tap.
 *	Option "LBCornerButton" "integer":
 *		Which mouse button is reported on a left bottom corner tap.
 *	Option "TapButton1" "integer":
 *		Which mouse button is reported on a non-corner one-finger tap.
 *	Option "TapButton2" "integer":
 *		Which mouse button is reported on a non-corner two-finger tap.
 *	Option "TapButton3" "integer":
 *		Which mouse button is reported on a non-corner
 *		three-finger tap.
 **/
func (tpad *Touchpad) EnableTapToClick(enabled bool) error {
	if enabled == tpad.CanTapToClick() {
		return nil
	}

	if globalWayland {
		return kwayland.EnableTapToClick(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id), enabled)
	}
	if tpad.isLibinputUsed {
		// TODO: libinput unsupported tap mapping settings.
		return libinputInt8PropSet(tpad.Id, libinputPropTapEnabled, enabled)
	}

	values, err := getInt8Prop(tpad.Id, propTapAction, 7)
	if err != nil {
		return err
	}

	if !enabled {
		values[4], values[5], values[6] = 0, 0, 0
	} else {
		// disable tap paste, because of conflicts with tap gesture
		if tpad.CanLeftHanded() {
			values[4], values[5], values[6] = 3, 1, 0
		} else {
			values[4], values[5], values[6] = 1, 3, 0
		}
	}

	return utils.SetInt8Prop(tpad.Id, propTapAction, values)
}

func (tpad *Touchpad) CanTapToClick() bool {
	if globalWayland {
		return kwayland.CanTapToClick(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id))
	}
	if tpad.isLibinputUsed {
		return libinputInt8PropCan(tpad.Id, libinputPropTapEnabled)
	}

	values, err := getInt8Prop(tpad.Id, propTapAction, 7)
	if err != nil {
		return false
	}

	if tpad.CanLeftHanded() {
		if values[4] == 3 && values[5] == 1 {
			return true
		}
	} else {
		if values[4] == 1 && values[5] == 3 {
			return true
		}
	}

	return false
}

/**
 * Property "Synaptics Edge Scrolling" 8 bit (BOOL), 3 values, vertical,
 * horizontal, corner. :
 *	Option "VertEdgeScroll" "boolean":
 *		Enable vertical scrolling when dragging along the right edge.
 *	Option "HorizEdgeScroll" "boolean" :
 *		Enable horizontal scrolling when dragging along
 *		the bottom edge.
 *	Option "CornerCoasting" "boolean":
 *		Enable edge scrolling to continue while the finger stays
 *		in an edge corner.
 **/
func (tpad *Touchpad) EnableEdgeScroll(enabled bool) error {
	if enabled == tpad.CanEdgeScroll() {
		return nil
	}

	if globalWayland {
		return kwayland.EnableScrollEdge(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id), enabled)
	}
	if tpad.isLibinputUsed {
		return libinputEnableScrollEdge(tpad.Id, enabled)
	}

	values, err := getInt8Prop(tpad.Id, propEdgeScroll, 3)
	if err != nil {
		return err
	}

	if enabled {
		values[0], values[1], values[2] = 1, 1, 1
	} else {
		values[0], values[1], values[2] = 0, 0, 0
	}

	return utils.SetInt8Prop(tpad.Id, propEdgeScroll, values)
}

func (tpad *Touchpad) CanEdgeScroll() bool {
	if globalWayland {
		return kwayland.CanScrollEdge(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id))
	}
	if tpad.isLibinputUsed {
		_, edge, _ := libinputCanScroll(tpad.Id)
		return edge
	}

	values, err := getInt8Prop(tpad.Id, propEdgeScroll, 3)
	if err != nil {
		return false
	}

	if values[0] != 1 || values[1] != 1 {
		return false
	}

	return true
}

/**
 * Property 'Synaptics Two-Finger Scrolling' 8 bit (BOOL),
 * 2 values, vertical, horizontal.
 *	Option "VertTwoFingerScroll" "boolean":
 *		Enable vertical scrolling when dragging with
 *		two fingers anywhere on the touchpad.
 *	Option "HorizTwoFingerScroll" "boolean" :
 *		Enable horizontal scrolling when dragging with
 *		two fingers anywhere on the touchpad.
 **/
func (tpad *Touchpad) EnableTwoFingerScroll(vert, horiz bool) error {
	if globalWayland {
		v := vert
		if !v {
			v = horiz
		}
		return kwayland.EnableScrollTwoFinger(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id), v)
	}
	oldVert, oldHoriz := tpad.CanTwoFingerScroll()
	if oldVert == vert && oldHoriz == horiz {
		return nil
	}

	if tpad.isLibinputUsed {
		err := libinputEnableScrollTwoFinger(tpad.Id, vert)
		if err != nil {
			return err
		}
		err = libinputInt8PropSet(tpad.Id, libinputPropHorizScrollEnabled, horiz)
		return err
	}

	var (
		newVert  int8 = 0
		newHoriz int8 = 0
	)
	if vert {
		newVert = 1
	}
	if horiz {
		newHoriz = 1
	}

	return utils.SetInt8Prop(tpad.Id, propTwoFingerScrol,
		[]int8{newVert, newHoriz})
}

func (tpad *Touchpad) CanTwoFingerScroll() (bool, bool) {
	if globalWayland {
		return true, kwayland.CanScrollTwoFinger(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id))
	}
	if tpad.isLibinputUsed {
		twoFinger, _, _ := libinputCanScroll(tpad.Id)
		return twoFinger, libinputInt8PropCan(tpad.Id, libinputPropHorizScrollEnabled)
	}

	values, err := getInt8Prop(tpad.Id, propTwoFingerScrol, 2)
	if err != nil {
		return false, false
	}

	return (values[0] == 1), (values[1] == 1)
}

/**
 * Property "Synaptics Scrolling Distance" 32 bit, 2 values, vert, horiz.
 *	Option "VertScrollDelta" "integer":
 *		Move distance of the finger for a scroll event.
 *	Option "HorizScrollDelta" "integer" :
 *		Move distance of the finger for a scroll event.
 *
 * if delta = 0, use value from property getting
 **/
func (tpad *Touchpad) EnableNaturalScroll(enabled bool) error {
	if enabled == tpad.CanNaturalScroll() {
		return nil
	}

	if globalWayland {
		return kwayland.EnableNaturalScroll(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id), enabled)
	}
	if tpad.isLibinputUsed {
		return libinputInt8PropSet(tpad.Id, libinputPropNaturalEnabled, enabled)
	}

	values, err := getInt32Prop(tpad.Id, propScrollDistance, 2)
	if err != nil {
		return err
	}

	if enabled {
		values[0], values[1] = -absInt32(values[0]), -absInt32(values[1])
	} else {
		values[0], values[1] = absInt32(values[0]), absInt32(values[1])
	}

	return utils.SetInt32Prop(tpad.Id, propScrollDistance, values)
}

func (tpad *Touchpad) CanNaturalScroll() bool {
	if globalWayland {
		return kwayland.CanNaturalScroll(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id))
	}
	if tpad.isLibinputUsed {
		return libinputInt8PropCan(tpad.Id, libinputPropNaturalEnabled)
	}

	values, err := getInt32Prop(tpad.Id, propScrollDistance, 2)
	if err != nil {
		return false
	}

	if values[0] < 0 || values[1] < 0 {
		return true
	}

	return false
}

func (tpad *Touchpad) SetScrollDistance(vert, horiz int32) error {
	if tpad.isLibinputUsed || globalWayland {
		return fmt.Errorf("Libinput unsupport the property")
	}

	oldVert, oldHoriz := tpad.ScrollDistance()
	if oldVert == vert && oldHoriz == horiz {
		return nil
	}

	if tpad.CanNaturalScroll() {
		vert = -vert
		horiz = -horiz
	}

	return utils.SetInt32Prop(tpad.Id, propScrollDistance,
		[]int32{vert, horiz})
}

func (tpad *Touchpad) ScrollDistance() (int32, int32) {
	if tpad.isLibinputUsed || globalWayland {
		return 0, 0
	}

	values, err := getInt32Prop(tpad.Id, propScrollDistance, 2)
	if err != nil {
		return 0, 0
	}

	return absInt32(values[0]), absInt32(values[1])
}

func (tpad *Touchpad) SetMotionAcceleration(accel float32) error {
	if globalWayland {
		return kwayland.SetPointerAccel(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id), float64(accel))
	}
	if tpad.isLibinputUsed {
		return libinputSetAccel(tpad.Id, 1-accel/1.5)
	}
	return setMotionAcceleration(tpad.Id, accel)
}

func (tpad *Touchpad) MotionAcceleration() (float32, error) {
	if globalWayland {
		v, err := kwayland.GetPointerAccel(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id))
		return float32(v), err
	}
	if tpad.isLibinputUsed {
		return libinputGetAccel(tpad.Id)
	}
	return getMotionAcceleration(tpad.Id)
}

func (tpad *Touchpad) SetMotionThreshold(thres float32) error {
	if tpad.isLibinputUsed || globalWayland {
		return fmt.Errorf("Libinput unsupport the property")
	}
	return setMotionThreshold(tpad.Id, thres)
}

func (tpad *Touchpad) MotionThreshold() (float32, error) {
	if tpad.isLibinputUsed || globalWayland {
		return 0, fmt.Errorf("Libinput unsupport the property")
	}
	return getMotionThreshold(tpad.Id)
}

func (tpad *Touchpad) SetMotionScaling(scaling float32) error {
	if tpad.isLibinputUsed || globalWayland {
		return fmt.Errorf("Libinput unsupport the property")
	}
	return setMotionScaling(tpad.Id, scaling)
}

func (tpad *Touchpad) MotionScaling() (float32, error) {
	if tpad.isLibinputUsed || globalWayland {
		return 0, fmt.Errorf("Libinput unsupport the property")
	}
	return getMotionScaling(tpad.Id)
}

func (tpad *Touchpad) SetRotation(direction uint8) error {
	return setRotation(tpad.Id, direction)
}

func (tpad *Touchpad) CanDisableWhileTyping() bool {
	if globalWayland {
		return kwayland.CanDisableWhileTyping(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id))
	}

	if !tpad.isLibinputUsed {
		return true
	}
	return libinputInt8PropCan(tpad.Id, libinputPropDisableWhileTyping)
}

func (tpad *Touchpad) EnableDisableWhileTyping(enabled bool) error {
	if globalWayland {
		return kwayland.EnableDisableWhileTyping(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id), enabled)
	}
	if !tpad.isLibinputUsed {
		return fmt.Errorf("Unsupported this prop: unused libinput as driver")
	}

	if enabled == tpad.CanDisableWhileTyping() {
		return nil
	}

	return libinputInt8PropSet(tpad.Id, libinputPropDisableWhileTyping, enabled)
}

// EnablePalmDetect set synaptics palm detect
// 'Synaptics Palm Detection' 8 bit (BOOL)
func (tpad *Touchpad) EnablePalmDetect(enabled bool) error {
	if tpad.isLibinputUsed || globalWayland {
		return fmt.Errorf("libinput unsupported palm detect setting")
	}

	if enabled == tpad.CanPalmDetect() {
		return nil
	}

	var values []int8
	if enabled {
		values = []int8{1}
	} else {
		values = []int8{0}
	}

	return utils.SetInt8Prop(tpad.Id, propPalmDetect, values)
}

func (tpad *Touchpad) CanPalmDetect() bool {
	if tpad.isLibinputUsed || globalWayland {
		// libinput enable palm detect as default
		return true
	}
	values, err := getInt8Prop(tpad.Id, propPalmDetect, 1)
	if err != nil {
		return false
	}
	return (values[0] == 1)
}

// 'Synaptics Palm Dimensions' 32 bit, 2 values, width, z
func (tpad *Touchpad) SetPalmDimensions(width, z int32) error {
	if tpad.isLibinputUsed || globalWayland {
		return fmt.Errorf("libinput unsupported palm detect setting")
	}

	oldWidth, oldZ, _ := tpad.GetPalmDimensions()
	if width == oldWidth && z == oldZ {
		return nil
	}
	return utils.SetInt32Prop(tpad.Id, propPalmDimensions, []int32{width, z})
}

func (tpad *Touchpad) GetPalmDimensions() (int32, int32, error) {
	if tpad.isLibinputUsed || globalWayland {
		return 0, 0, fmt.Errorf("libinput unsupported palm detect setting")
	}

	values, err := getInt32Prop(tpad.Id, propPalmDimensions, 2)
	if err != nil {
		return 0, 0, err
	}
	return values[0], values[1], nil
}

// 校验是否已经使能了触摸板中键模拟功能
func (tpad *Touchpad) IsMiddleButtonEnulationEnabled() bool {
	if globalWayland {
		return kwayland.CanMiddleButtonEmulation(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id))
	}

	if !tpad.isLibinputUsed {
		return false
	}
	return libinputInt8PropCan(tpad.Id, libinputPropMiddleEmulationEnabled)
}

// 触摸板中键模拟使能接口
func (tpad *Touchpad) EnableMiddleButtonEmulation(enabled bool) error {
	if globalWayland {
		return kwayland.EnableMiddleEmulation(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, tpad.Id), enabled)
	}
	if !tpad.isLibinputUsed {
		return fmt.Errorf("unsupported this prop: unused libinput as driver")
	}

	if enabled == tpad.IsMiddleButtonEnulationEnabled() {
		return nil
	}

	return libinputInt8PropSet(tpad.Id, libinputPropMiddleEmulationEnabled, enabled)
}
