/*
 * Copyright (C) 2018 ~ 2018 Deepin Technology Co., Ltd.
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
	"pkg.deepin.io/dde/api/dxinput/utils"
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
	info := utils.ListDevice().Get(id)
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
	// only supported libinput
	if !touch.isLibinputUsed {
		return fmt.Errorf("Unsupport rotation for (%d - %s)", touch.Id, touch.Name)
	}
	return setRotation(touch.Id, direction)
}

func (touch *Touchscreen) SetTransformationMatrix(m [9]float32) error {
	// only supported libinput
	if !touch.isLibinputUsed {
		return fmt.Errorf("Unsupport transformation matrix for (%d - %s)", touch.Id, touch.Name)
	}

	return setTransformationMatrix(touch.Id, m)
}
