/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
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
