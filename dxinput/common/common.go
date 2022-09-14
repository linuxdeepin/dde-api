// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package common

const (
	DevTypeUnknown int32 = iota
	DevTypeKeyboard
	DevTypeMouse
	DevTypeTouchpad
	DevTypeWacom
	DevTypeTouchscreen
)

type DeviceInfo struct {
	Id      int32
	Type    int32
	Name    string
	Enabled bool
}
type DeviceInfos []*DeviceInfo

func (infos DeviceInfos) Get(id int32) *DeviceInfo {
	for _, info := range infos {
		if info.Id == id {
			return info
		}
	}
	return nil
}
