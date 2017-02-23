/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dxinput

import (
	"fmt"
	"pkg.deepin.io/dde/api/dxinput/utils"
	"strings"
)

const (
	propMidBtnEmulation        string = "Evdev Middle Button Emulation"
	propMidBtnEmulationTimeout        = "Evdev Middle Button Timeout"
	propWheelEmulation                = "Evdev Wheel Emulation"
	propWheelEmulationButton          = "Evdev Wheel Emulation Button"
	propWheelEmulationTimeout         = "Evdev Wheel Emulation Timeout"
	propWheelEmulationAxes            = "Evdev Wheel Emulation Axes"
	propEvdevScrollDistance           = "Evdev Scrolling Distance"
)

type Mouse struct {
	Id         int32
	Name       string
	TrackPoint bool

	// -1: unknow, 0: not used, 1: used
	isLibinputUsed bool
}

func NewMouse(id int32) (*Mouse, error) {
	info := utils.ListDevice().Get(id)
	if info == nil {
		return nil, fmt.Errorf("Invalid device id: %v", id)
	}
	return NewMouseFromDeviceInfo(info)
}

func NewMouseFromDeviceInfo(dev *utils.DeviceInfo) (*Mouse, error) {
	if dev == nil || dev.Type != utils.DevTypeMouse {
		return nil, fmt.Errorf("Not a mouse device(%d - %s)", dev.Id, dev.Name)
	}
	return &Mouse{
		Id:             dev.Id,
		Name:           dev.Name,
		TrackPoint:     isTrackPoint(dev),
		isLibinputUsed: utils.IsPropertyExist(dev.Id, libinputPropButtonScrollingButton),
	}, nil
}

func (m *Mouse) Enable(enabled bool) error {
	return enableDevice(m.Id, enabled)
}

func (m *Mouse) IsEnabled() bool {
	return isDeviceEnabled(m.Id)
}

func (m *Mouse) EnableLeftHanded(enabled bool) error {
	if enabled == m.CanLeftHanded() {
		return nil
	}

	return utils.SetLeftHanded(uint32(m.Id), m.Name, enabled)
}

func (m *Mouse) CanLeftHanded() bool {
	return utils.CanLeftHanded(uint32(m.Id), m.Name)
}

// EnableMiddleButtonEmulation enable mouse middle button emulation
// "Evdev Middle Button Emulation"
//     1 boolean value (8 bit, 0 or 1).
func (m *Mouse) EnableMiddleButtonEmulation(enabled bool) error {
	if enabled == m.CanMiddleButtonEmulation() {
		return nil
	}

	if m.isLibinputUsed {
		return libinputInt8PropSet(m.Id, libinputPropMiddleEmulationEnabled, enabled)
	}

	var values []int8
	if enabled {
		values = []int8{1}
	} else {
		values = []int8{0}
	}

	return utils.SetInt8Prop(m.Id, propMidBtnEmulation, values)
}

func (m *Mouse) CanMiddleButtonEmulation() bool {
	if m.isLibinputUsed {
		return libinputInt8PropCan(m.Id, libinputPropMiddleEmulationEnabled)
	}

	values, err := getInt32Prop(m.Id, propMidBtnEmulation, 1)
	if err != nil {
		return false
	}

	return (values[0] == 1)
}

// SetMiddleButtonEmulationTimeout set middle button emulation timeout
// "Evdev Middle Button Timeout"
//     1 16-bit positive value.
func (m *Mouse) SetMiddleButtonEmulationTimeout(timeout int16) error {
	if m.isLibinputUsed {
		return fmt.Errorf("Libinput unsupport the property")
	}

	old, err := m.MiddleButtonEmulationTimeout()
	if err == nil && timeout == old {
		return nil
	}
	return utils.SetInt16Prop(m.Id, propMidBtnEmulationTimeout,
		[]int16{timeout})
}

func (m *Mouse) MiddleButtonEmulationTimeout() (int16, error) {
	if m.isLibinputUsed {
		return 0, fmt.Errorf("Libinput unsupport the property")
	}

	values, err := getInt16Prop(m.Id, propMidBtnEmulationTimeout, 1)
	if err != nil {
		return 0, err
	}
	return values[0], nil
}

// EnableWheelEmulation enable mouse wheel emulation
// "Evdev Wheel Emulation"
//    1 boolean value (8 bit, 0 or 1).
func (m *Mouse) EnableWheelEmulation(enabled bool) error {
	if enabled == m.CanWheelEmulation() {
		return nil
	}

	if m.isLibinputUsed {
		return libinputEnableScrollButton(m.Id, enabled)
	}

	var values []int8
	if enabled {
		values = []int8{1}
	} else {
		values = []int8{0}
	}
	return utils.SetInt8Prop(m.Id, propWheelEmulation, values)
}

func (m *Mouse) SetRotation(direction uint8) error {
	if m.isLibinputUsed {
		return fmt.Errorf("libinput unsupported transformation matrix")
	}
	return setRotation(m.Id, direction)
}

func (m *Mouse) CanWheelEmulation() bool {
	if m.isLibinputUsed {
		_, _, v := libinputCanScroll(m.Id)
		return v
	}

	values, err := getInt8Prop(m.Id, propWheelEmulation, 1)
	if err != nil {
		return false
	}
	return (values[0] == 1)
}

// SetWheelEmulationButton set wheel emulation button
// "Evdev Wheel Emulation Button"
//    1 8-bit value, allowed range 0-32, 0 disables the button.
func (m *Mouse) SetWheelEmulationButton(btnNum int8) error {
	old, _ := m.WheelEmulationButton()
	if btnNum == old {
		return nil
	}

	if m.isLibinputUsed {
		return libinputSetScrollButton(m.Id, int32(btnNum))
	}

	return utils.SetInt8Prop(m.Id, propWheelEmulationButton,
		[]int8{btnNum})
}

func (m *Mouse) WheelEmulationButton() (int8, error) {
	if m.isLibinputUsed {
		v, err := libinputGetScrollButton(m.Id)
		return int8(v), err
	}

	values, err := getInt8Prop(m.Id, propWheelEmulationButton, 1)
	if err != nil {
		return -1, err
	}
	return values[0], nil
}

// SetWheelEmulationTimeout set wheel emulation timeout
// "Evdev Wheel Emulation Timeout"
//     1 16-bit positive value.
func (m *Mouse) SetWheelEmulationTimeout(timeout int16) error {
	if m.isLibinputUsed {
		return fmt.Errorf("Libinput unsupport the property")
	}

	old, err := m.WheelEmulationTimeout()
	if err == nil && timeout == old {
		return nil
	}
	return utils.SetInt16Prop(m.Id, propWheelEmulationTimeout,
		[]int16{timeout})
}

func (m *Mouse) WheelEmulationTimeout() (int16, error) {
	if m.isLibinputUsed {
		return 0, fmt.Errorf("Libinput unsupport the property")
	}

	values, err := getInt16Prop(m.Id, propWheelEmulationTimeout, 1)
	if err != nil {
		return 0, err
	}
	return values[0], nil
}

func (m *Mouse) EnableWheelHorizScroll(enabled bool) error {
	if enabled == m.CanWheelHorizScroll() {
		return nil
	}

	if m.isLibinputUsed {
		return libinputInt8PropSet(m.Id, libinputPropHorizScrollEnabled, enabled)
	}
	return m.enableWheelHorizScroll(enabled, false)
}

func (m *Mouse) CanWheelHorizScroll() bool {
	if m.isLibinputUsed {
		return libinputInt8PropCan(m.Id, libinputPropHorizScrollEnabled)
	}
	return m.canWheelHorizScroll(false)
}

func (m *Mouse) EnableWheelHorizNaturalScroll(enabled bool) error {
	if enabled == m.CanWheelHorizNaturalScroll() {
		return nil
	}
	if m.isLibinputUsed {
		return libinputInt8PropSet(m.Id, libinputPropNaturalEnabled, enabled)
	}
	return m.enableWheelHorizScroll(enabled, true)
}

func (m *Mouse) CanWheelHorizNaturalScroll() bool {
	if m.isLibinputUsed {
		return libinputInt8PropCan(m.Id, libinputPropNaturalEnabled)
	}
	return m.canWheelHorizScroll(true)
}

func (m *Mouse) EnableNaturalScroll(enabled bool) error {
	if enabled == m.CanNaturalScroll() {
		return nil
	}

	if m.isLibinputUsed {
		return libinputInt8PropSet(m.Id, libinputPropNaturalEnabled, enabled)
	}
	values, err := getInt32Prop(m.Id, propEvdevScrollDistance, 3)
	if err != nil {
		return err
	}

	if enabled {
		values[0], values[1], values[2] = -absInt32(values[0]), -absInt32(values[1]), -absInt32(values[2])
	} else {
		values[0], values[1], values[2] = absInt32(values[0]), absInt32(values[1]), absInt32(values[2])
	}

	return utils.SetInt32Prop(m.Id, propEvdevScrollDistance, values)
}

func (m *Mouse) CanNaturalScroll() bool {
	if m.isLibinputUsed {
		return libinputInt8PropCan(m.Id, libinputPropNaturalEnabled)
	}
	values, err := getInt32Prop(m.Id, propEvdevScrollDistance, 3)
	if err != nil {
		return false
	}

	if values[0] < 0 || values[1] < 0 || values[2] < 0 {
		return true
	}
	return false
}

func (m *Mouse) SetMotionAcceleration(accel float32) error {
	if m.isLibinputUsed {
		return libinputSetAccel(m.Id, 1-accel/1.5)
	}
	return setMotionAcceleration(m.Id, accel)
}

func (m *Mouse) MotionAcceleration() (float32, error) {
	if m.isLibinputUsed {
		return libinputGetAccel(m.Id)
	}
	return getMotionAcceleration(m.Id)
}

func (m *Mouse) SetMotionThreshold(thres float32) error {
	if m.isLibinputUsed {
		return fmt.Errorf("Libinput unsupport the property")
	}

	return setMotionThreshold(m.Id, thres)
}

func (m *Mouse) MotionThreshold() (float32, error) {
	if m.isLibinputUsed {
		return 0.0, fmt.Errorf("Libinput unsupport the property")
	}

	return getMotionThreshold(m.Id)
}

func (m *Mouse) SetMotionScaling(scaling float32) error {
	if m.isLibinputUsed {
		return fmt.Errorf("Libinput unsupport the property")
	}

	return setMotionScaling(m.Id, scaling)
}

func (m *Mouse) MotionScaling() (float32, error) {
	if m.isLibinputUsed {
		return 0.0, fmt.Errorf("Libinput unsupport the property")
	}

	return getMotionScaling(m.Id)
}

// setWheelEmulationAxes set wheel horizontal scrolling
// "Evdev Wheel Emulation Axes"
//     4 8-bit values, order X up, X down, Y up, Y down. 0 disables a value.
// set to "6 7 4 5", enable horizontal scrolling
// default: "0 0 4 5"
func (m *Mouse) setWheelEmulationAxes(values []int8) error {
	old, err := m.wheelEmulationAxes()
	if err == nil && isInt8ArrayEqual(values, old) {
		return nil
	}
	return utils.SetInt8Prop(m.Id, propWheelEmulationAxes, values)
}

func (m *Mouse) wheelEmulationAxes() ([]int8, error) {
	values, err := getInt8Prop(m.Id, propWheelEmulationAxes, 4)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (m *Mouse) enableWheelHorizScroll(enabled, natural bool) error {
	var values []int8
	if enabled {
		if natural {
			values = []int8{7, 6, 5, 4}
		} else {
			values = []int8{6, 7, 4, 5}
		}
	} else {
		values = []int8{0, 0, 4, 5}
	}
	return utils.SetInt8Prop(m.Id, propWheelEmulationAxes, values)
}

func (m *Mouse) canWheelHorizScroll(natural bool) bool {
	values, err := m.wheelEmulationAxes()
	if err != nil {
		return false
	}

	if natural {
		return isInt8ArrayEqual(values, []int8{7, 6, 5, 4})
	}
	return isInt8ArrayEqual(values, []int8{6, 7, 4, 5})
}

func isTrackPoint(info *utils.DeviceInfo) bool {
	name := strings.ToLower(info.Name)
	return strings.Contains(name, "trackpoint") ||
		strings.Contains(name, "dualpoint stick")
}

func isInt8ArrayEqual(a1, a2 []int8) bool {
	l1 := len(a1)
	l2 := len(a2)
	if l1 != l2 {
		return false
	}
	for i := 0; i < l1; i++ {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}
