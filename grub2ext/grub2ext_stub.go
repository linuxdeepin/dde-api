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

func (grub *Grub2Ext) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.api.Grub2",
		"/com/deepin/api/Grub2",
		"com.deepin.api.Grub2",
	}
}

func (grub *Grub2Ext) DoWriteSettings(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(_GRUB_CONFIG_FILE, []byte(fileContent), 0664)
	if err != nil {
		_LOGGER.Error(err.Error())
		return false, err
	}
	return true, nil
}

func (grub *Grub2Ext) DoWriteCacheConfig(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(_GRUB_CACHE_FILE, []byte(fileContent), 0644)
	if err != nil {
		_LOGGER.Error(err.Error())
		return false, err
	}
	return true, nil
}

func (grub *Grub2Ext) DoGenerateGrubConfig() (ok bool, err error) {
	_LOGGER.Info("start to generate a new grub configuration file")
	_, stderr, err := execAndWait(30, _GRUB_UPDATE_EXE)
	_LOGGER.Info("process output: %s", stderr)
	if err != nil {
		_LOGGER.Error("generate grub configuration failed")
		return false, err
	}
	_LOGGER.Info("generate grub configuration finished")
	return true, nil
}

func (grub *Grub2Ext) DoSetThemeBackgroundSourceFile(imageFile string, screenWidth, screenHeight uint16) (ok bool, err error) {
	// backup background source file
	_, err = copyFile(imageFile, _THEME_BG_SRC_FILE)
	if err != nil {
		return false, err
	}

	// generate a new background
	return grub.DoGenerateThemeBackground(screenWidth, screenHeight)
}

func (grub *Grub2Ext) DoGenerateThemeBackground(screenWidth, screenHeight uint16) (ok bool, err error) {
	imgWidth, imgHeight, err := graphic.GetImageSize(_THEME_BG_SRC_FILE)
	if err != nil {
		_LOGGER.Error(err.Error())
		return false, err
	}
	_LOGGER.Info("source background size %dx%d", imgWidth, imgHeight)

	w, h := getImgClipSizeByResolution(screenWidth, screenHeight, imgWidth, imgHeight)
	_LOGGER.Info("background size %dx%d", w, h)
	err = graphic.ClipPNG(_THEME_BG_SRC_FILE, _THEME_BG_FILE, 0, 0, w, h)
	if err != nil {
		_LOGGER.Error(err.Error())
		return false, err
	}
	return true, nil
}

func (grub *Grub2Ext) DoCustomTheme(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(_THEME_MAIN_FILE, []byte(fileContent), 0664)
	if err != nil {
		_LOGGER.Error(err.Error())
		return false, err
	}
	return true, nil
}

func (grub *Grub2Ext) DoWriteThemeJson(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(_THEME_JSON_FILE, []byte(fileContent), 0664)
	if err != nil {
		_LOGGER.Error(err.Error())
		return false, err
	}
	return true, nil
}
