// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package dxinput

import (
	"errors"
	"fmt"

	. "github.com/linuxdeepin/dde-api/dxinput/common"
	"github.com/linuxdeepin/dde-api/dxinput/kwayland"
	"github.com/linuxdeepin/dde-api/dxinput/utils"
)

func SetKeyboardRepeat(enabled bool, delay, interval uint32) error {
	return utils.SetKeyboardRepeat(enabled, delay, interval)
}

type Keyboard struct {
	Id   int32
	Name string
}

func NewKeyboard(id int32) (*Keyboard, error) {
	infos := utils.ListDevice()
	if infos == nil {
		return nil, errors.New("no device")
	}

	info := infos.Get(id)

	if info == nil {
		return nil, fmt.Errorf("invalid device id: %v", id)
	}
	return NewKeyboardDevInfo(info)
}

func NewKeyboardDevInfo(dev *DeviceInfo) (*Keyboard, error) {
	if dev == nil || dev.Type != DevTypeKeyboard {
		return nil, fmt.Errorf("not a keyboard device(%d - %s)", dev.Id, dev.Name)
	}

	return &Keyboard{
		Id:   dev.Id,
		Name: dev.Name,
	}, nil
}

func (m *Keyboard) Enable(enabled bool) error {
	if globalWayland {
		return kwayland.Enable(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, m.Id), enabled)
	}

	return enableDevice(m.Id, enabled)
}

func (m *Keyboard) IsEnabled() bool {
	if globalWayland {
		return kwayland.CanEnabled(fmt.Sprintf("%s%d", kwayland.SysNamePrefix, m.Id))
	}

	return isDeviceEnabled(m.Id)
}
