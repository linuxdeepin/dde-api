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
	"dlib/graphic"
	"io/ioutil"
)

// GetDBusInfo implement interface of dbus.DBusObject
func (grub *Grub2Ext) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.api.Grub2",
		"/com/deepin/api/Grub2",
		"com.deepin.api.Grub2",
	}
}

// DoWriteSettings write file content to "/etc/default/grub".
func (grub *Grub2Ext) DoWriteSettings(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(grubConfigFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoWriteCacheConfig write file content to "/var/cache/dde-daemon/grub2.json".
func (grub *Grub2Ext) DoWriteCacheConfig(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(grubCacheFile, []byte(fileContent), 0644)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoGenerateGrubConfig execute command "/usr/sbin/update-grub" to
// generate a new grub configuration.
func (grub *Grub2Ext) DoGenerateGrubConfig() (ok bool, err error) {
	logger.Info("start to generate a new grub configuration file")
	_, stderr, err := execAndWait(30, grubUpdateExe)
	logger.Info("process output: %s", stderr)
	if err != nil {
		logger.Error("generate grub configuration failed: %v", err)
		return false, err
	}
	logger.Info("generate grub configuration successful")
	return true, nil
}

// DoSetThemeBackgroundSourceFile setup a new background source file
// for deepin grub2 theme, and then generate the background depends on
// screen resolution.
func (grub *Grub2Ext) DoSetThemeBackgroundSourceFile(imageFile string, screenWidth, screenHeight uint16) (ok bool, err error) {
	// backup background source file
	_, err = copyFile(imageFile, themeBgSrcFile)
	if err != nil {
		return false, err
	}

	// generate a new background
	return grub.DoGenerateThemeBackground(screenWidth, screenHeight)
}

// DoGenerateThemeBackground generate the background for deepin grub2
// theme depends on screen resolution.
func (grub *Grub2Ext) DoGenerateThemeBackground(screenWidth, screenHeight uint16) (ok bool, err error) {
	imgWidth, imgHeight, err := graphic.GetImageSize(themeBgSrcFile)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	logger.Info("source background size %dx%d", imgWidth, imgHeight)

	w, h := getImgClipSizeByResolution(screenWidth, screenHeight, imgWidth, imgHeight)
	logger.Info("background size %dx%d", w, h)
	err = graphic.ClipPNG(themeBgSrcFile, themeBgFile, 0, 0, w, h)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoCustomTheme write file content to "/boot/grub/themes/deepin/theme.txt".
func (grub *Grub2Ext) DoCustomTheme(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(themeMainFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoWriteThemeJson write file content to "/boot/grub/themes/deepin/theme_tpl.json".
func (grub *Grub2Ext) DoWriteThemeJson(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(themeJSONFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}
