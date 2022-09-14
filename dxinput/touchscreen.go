// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package dxinput

import (
	"fmt"
	"errors"

	. "github.com/linuxdeepin/dde-api/dxinput/common"
	"github.com/linuxdeepin/dde-api/dxinput/utils"
)

const (
	libinputPropCalibrationMatrix = "libinput Calibration Matrix"
)

type Touchscreen struct {
	Id   int32
	Name string

	// -1: unknown, 0: not used, 1: used
	isLibinputUsed bool
}

func NewTouchscreen(id int32) (*Touchscreen, error) {
	infos := utils.ListDevice()
	if infos == nil {
		return nil, errors.New("No device")
	}

	info := infos.Get(id)

	if info == nil {
		return nil, fmt.Errorf("Invalid device id: %d", id)
	}
	return NewTouchscreenFromDevInfo(info)
}

func NewTouchscreenFromDevInfo(dev *DeviceInfo) (*Touchscreen, error) {
	if dev == nil || dev.Type != DevTypeTouchscreen {
		return nil, fmt.Errorf("Not a touchscreen device(%d - %s)",
			dev.Id, dev.Name)
	}
	return &Touchscreen{
		Id:             dev.Id,
		Name:           dev.Name,
		isLibinputUsed: utils.IsPropertyExist(dev.Id, libinputPropCalibrationMatrix),
	}, nil
}

func (touch *Touchscreen) Enable(enabled bool) error {
	return enableDevice(touch.Id, enabled)
}

func (touch *Touchscreen) IsEnabled() bool {
	return isDeviceEnabled(touch.Id)
}

func (touch *Touchscreen) SetRotation(direction uint8) error {
	return setRotation(touch.Id, direction)
}

func (touch *Touchscreen) SetTransformationMatrix(m [9]float32) error {
	return setTransformationMatrix(touch.Id, m)
}
