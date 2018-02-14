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

package main

import (
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/utils"
)

const (
	deviceServiceName = "com.deepin.api.Device"
	devicePath        = "/com/deepin/api/Device"
	deviceInterface   = deviceServiceName
	rfkillBin         = "/usr/sbin/rfkill"
)

type Device struct {
	service *dbusutil.Service
	methods *struct {
		UnblockDevice func() `in:"deviceType"`
		BlockDevice   func() `in:"deviceType"`
	}
}

func (d *Device) GetDBusExportInfo() dbusutil.ExportInfo {
	return dbusutil.ExportInfo{
		Path:      devicePath,
		Interface: deviceInterface,
	}
}

// UnblockDevice unblock target devices through rfkill, the device
// type could be all, wifi, wlan, bluetooth, uwb, ultrawideband,
// wimax, wwan, gps, fm, and nfc.
func (d *Device) UnblockDevice(deviceType string) *dbus.Error {
	d.service.DelayAutoQuit()
	_, _, err := utils.ExecAndWait(5, rfkillBin, "unblock", deviceType)
	return dbusutil.ToError(err)
}

// BlockDevice block target devices through rfkill, the device
// type could be all, wifi, wlan, bluetooth, uwb, ultrawideband,
// wimax, wwan, gps, fm, and nfc. Need polkit authentication.
func (d *Device) BlockDevice(deviceType string) *dbus.Error {
	d.service.DelayAutoQuit()
	// TODO need polkit authentication
	_, _, err := utils.ExecAndWait(5, rfkillBin, "block", deviceType)
	return dbusutil.ToError(err)
}
