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

package main

import (
	"dlib/dbus"
	liblogger "dlib/logger"
	"os"
)

const (
	grubConfigFile = "/etc/default/grub"
	grubUpdateExe  = "/usr/sbin/update-grub"
	grubCacheFile  = "/var/cache/dde-daemon/grub2.json"

	themePath      = "/boot/grub/themes/deepin"
	themeMainFile  = themePath + "/theme.txt"
	themeJSONFile  = themePath + "/theme_tpl.json"
	themeBgSrcFile = themePath + "/background_source"
	themeBgFile    = themePath + "/background.png"
)

var logger = liblogger.NewLogger("dde-api/grub2ext")

// Grub2Ext is a dbus object, and is split from dde-daemon/grub2 to
// fix launch issue through dbus-daemon for that system bus in root
// couldn't access session bus interface.
type Grub2Ext struct{}

// NewGrub2Ext create a Grub2Ext object.
func NewGrub2Ext() *Grub2Ext {
	grub := &Grub2Ext{}
	return grub
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Fatal("%v", err)
		}
	}()

	// configure logger
	logger.SetRestartCommand("/usr/lib/deepin-api/grub2ext", "--debug")
	if stringInSlice("-d", os.Args) || stringInSlice("--debug", os.Args) {
		logger.SetLogLevel(liblogger.LEVEL_DEBUG)
	}

	grub := NewGrub2Ext()
	err := dbus.InstallOnSystem(grub)
	if err != nil {
		panic(err)
	}

	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); err != nil {
		liblogger.Printf("lost dbus session: %v\n", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
