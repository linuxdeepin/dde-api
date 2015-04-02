/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
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
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

const (
	defaultLocaleFile    = "/etc/default/locale"
	defaultLocaleGenFile = "/etc/locale.gen"
)

func (h *Helper) SetLocale(locale string) error {
	if !IsLocaleValid(locale) {
		return fmt.Errorf("Invalid locale:", locale)
	}

	return writeContentToFile(defaultLocaleFile,
		fmt.Sprintf("LANG=%s", locale))
}

func (h *Helper) GenerateLocale(locale string) error {
	h.running = true
	defer func() {
		h.running = false
	}()

	if !IsLocaleValid(locale) {
		dbus.Emit(h, "Success", false,
			fmt.Sprintf("Invalid locale: %v", locale))
		return fmt.Errorf("Invalid locale:", locale)
	}

	// locales version <= 2.13
	if !dutils.IsFileExist(defaultLocaleGenFile) {
		err := h.doGenLocaleWithParam(locale)
		if err != nil {
			dbus.Emit(h, "Success", false, err.Error())
			return err
		}

		dbus.Emit(h, "Success", true, "")
		return nil
	}

	err := enableLocaleInFile(locale, defaultLocaleGenFile)
	if err != nil {
		dbus.Emit(h, "Success", false, err.Error())
		return err
	}

	err = h.doGenLocale()
	if err != nil {
		dbus.Emit(h, "Success", false, err.Error())
		return err
	}

	dbus.Emit(h, "Success", true, "")
	return nil
}

func enableLocaleInFile(locale, file string) error {
	finfo, err := NewLocaleFileInfo(file)
	if err != nil {
		return err
	}

	if finfo.IsLocaleEnabled(locale) {
		return nil
	}

	finfo.EnableLocale(locale)
	err = finfo.Save(defaultLocaleGenFile)
	if err != nil {
		return err
	}

	return nil
}
