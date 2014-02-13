/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
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
	"dlib/logger"
	"io/ioutil"
	"os"
)

const (
	ETC_TIMEZONE  = "/etc/timezone"
	ETC_LOCALTIME = "/etc/localtime"
	ZONE_INFO_DIR = "/usr/share/zoneinfo/"
	ETC_PERM      = 0644
)

func getTimezone() (string, bool) {
	contents, err := ioutil.ReadFile(ETC_TIMEZONE)
	if err != nil {
		logger.Printf("ReadFile '%s' failed: %s\n",
			ETC_TIMEZONE, err)
		return "", false
	}

	return string(contents), true
}

func setTimezone(tz string) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("Recover Error:", err)
		}
	}()

	if !fileIsRegular(ETC_LOCALTIME) {
		return false
	}

	tzFile := ZONE_INFO_DIR + tz
	if !fileIsRegular(tzFile) {
		return false
	}

	/* Modify /etc/localtime */
	if fileIsSymlink(ETC_LOCALTIME) {
		err := os.Remove(ETC_LOCALTIME)
		if err != nil {
			logger.Printf("Remove '%s' failed: %s\n",
				ETC_TIMEZONE, err)
			return false
		}

		err = os.Symlink(tzFile, ETC_TIMEZONE)
		if err != nil {
			logger.Printf("Symlink '%s' to '%s' failed: %s\n",
				tzFile, ETC_TIMEZONE, err)
			return false
		}
	} else {
		if !copyFile(tzFile, ETC_LOCALTIME, ETC_PERM) {
			return false
		}
	}

	/* Modify /etc/timezone */
	if !fileIsRegular(ETC_TIMEZONE) {
		return false
	}
	err := ioutil.WriteFile(ETC_TIMEZONE, []byte(tz),
		os.FileMode(ETC_PERM))
	if err != nil {
		logger.Printf("WriteFile '%s' failed: %s\n", ETC_TIMEZONE, err)
		return false
	}

	return true
}

func copyFile(src, dest string, perm os.FileMode) bool {
	contents, err := ioutil.ReadFile(src)
	if err != nil {
		logger.Printf("ReadFile '%s' failed: %s\n", src, err)
		return false
	}

	err = ioutil.WriteFile(dest, contents, os.FileMode(ETC_PERM))
	if err != nil {
		logger.Printf("WriteFile '%s' failed: %s\n", dest, err)
		return false
	}

	return true
}

func getFileMode(file string) os.FileMode {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		logger.Printf("Open '%s' failed: %s\n", file, err)
		panic(err)
	}

	info, err1 := f.Stat()
	if err1 != nil {
		logger.Printf("Stat '%s' failed: %s\n", file, err1)
		panic(err1)
	}

	return info.Mode()
}

func fileIsRegular(file string) bool {
	ok := getFileMode(file).IsRegular()
	if !ok {
		logger.Printf("'%s' is not regular\n", file)
		return false
	}

	return true
}

func fileIsDir(file string) bool {
	ok := getFileMode(file).IsDir()
	if !ok {
		logger.Printf("'%s' is not dir\n", file)
		return false
	}

	return true
}

func fileIsSymlink(file string) bool {
	mode := getFileMode(file)
	if mode == os.ModeSymlink {
		logger.Printf("'%s' is symlink\n", file)
		return true
	}

	return false
}

func zoneFileIsExist(file string) bool {
	if _, err := os.Stat(file); os.IsExist(err) {
		return true
	}
	logger.Printf("'%s' is not exist\n", file)

	return false
}
