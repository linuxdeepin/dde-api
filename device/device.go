/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/utils"
)

const (
	deviceDest = "com.deepin.api.Device"
	devicePath = "/com/deepin/api/Device"
	deviceObj  = "com.deepin.api.Device"
	rfkillBin  = "/usr/sbin/rfkill"
)

type Device struct{}

func (d *Device) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       deviceDest,
		ObjectPath: devicePath,
		Interface:  deviceObj,
	}
}

// UnblockDevice unblock target devices through rfkill, the device
// type could be all, wifi, wlan, bluetooth, uwb, ultrawideband,
// wimax, wwan, gps, fm, and nfc.
func (d *Device) UnblockDevice(deviceType string) (err error) {
	_, _, err = utils.ExecAndWait(5, rfkillBin, "unblock", deviceType)
	return
}

// BlockDevice block target devices through rfkill, the device
// type could be all, wifi, wlan, bluetooth, uwb, ultrawideband,
// wimax, wwan, gps, fm, and nfc. Need polkit authentication.
func (d *Device) BlockDevice(deviceType string) (err error) {
	// TODO need polkit authentication
	_, _, err = utils.ExecAndWait(5, rfkillBin, "block", deviceType)
	return
}
