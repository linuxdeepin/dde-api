/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
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

// This module is split from dde-daemon/grub2 to fix launch issue
// through dbus-daemon for that system bus in root couldn't access
// session bus interface.

package main

import (
	"dlib/dbus"
	"dlib/logger"
)

const (
	_GRUB_CONFIG_FILE = "/etc/default/grub"
	_GRUB_UPDATE_EXE  = "/usr/sbin/update-grub"
	_GRUB_CACHE_FILE  = "/var/cache/dde-daemon/grub2.json"

	_THEME_PATH        = "/boot/grub/themes/deepin"
	_THEME_MAIN_FILE   = _THEME_PATH + "/theme.txt"
	_THEME_JSON_FILE   = _THEME_PATH + "/theme_tpl.json"
	_THEME_BG_SRC_FILE = _THEME_PATH + "/background_source"
	_THEME_BG_FILE     = _THEME_PATH + "/background.png"
)

var _LOGGER, _ = logger.New("dde-api/grub2ext")

type Grub2Ext struct{}

func NewGrub2Ext() *Grub2Ext {
	grub := &Grub2Ext{}
	return grub
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			_LOGGER.Fatal("%v", err)
		}
	}()

	grub := NewGrub2Ext()
	err := dbus.InstallOnSystem(grub)
	if err != nil {
		panic(err)
	}

	dbus.DealWithUnhandledMessage()

	select {}
}
