// Copyright (c) 2015 Deepin Ltd. All rights reserved.
// Use of this source is govered by General Public License that can be found
// in the LICENSE file.
package main

import (
	"pkg.deepin.io/lib/dbus"
)

type Validator struct{}

// GetDBusInfo implements dbus.DBusObject interface
func (validator *Validator) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DBusName,
		DBusPath,
		DBusInterface,
	}
}
