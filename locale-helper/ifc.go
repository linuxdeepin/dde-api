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
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/polkit"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	polkitManageLocale = "com.deepin.api.locale-helper.manage-locale"

	defaultLocaleFile    = "/etc/default/locale"
	defaultLocaleGenFile = "/etc/locale.gen"
)

var errAuthFailed = fmt.Errorf("authentication failed")

func (h *Helper) SetLocale(dmessage dbus.DMessage, locale string) error {
	ok, err := checkAuth(dmessage)
	logger.Debug("---Auth ret:", ok, err)
	if !ok || err != nil {
		return errAuthFailed
	}

	if !IsLocaleValid(locale) {
		return fmt.Errorf("invalid locale: %s", locale)
	}

	return writeContentToFile(defaultLocaleFile,
		fmt.Sprintf("LANG=%s", locale))
}

func (h *Helper) GenerateLocale(dmessage dbus.DMessage, locale string) error {
	h.running = true
	defer func() {
		h.running = false
	}()

	ok, err := checkAuth(dmessage)
	logger.Debug("---Auth ret:", ok, err)
	if !ok || err != nil {
		dbus.Emit(h, "Success", false, errAuthFailed.Error())
		return errAuthFailed
	}

	if !IsLocaleValid(locale) {
		dbus.Emit(h, "Success", false,
			fmt.Sprintf("Invalid locale: %v", locale))
		return fmt.Errorf("invalid locale: %s", locale)
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

	err = enableLocaleInFile(locale, defaultLocaleGenFile)
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

func init() {
	polkit.Init()
}

func checkAuth(dmessage dbus.DMessage) (bool, error) {
	subject := polkit.NewSubject(polkit.SubjectKindUnixProcess)
	subject.SetDetail("pid", dmessage.GetSenderPID())
	subject.SetDetail("start-time", uint64(0))
	details := make(map[string]string)
	result, err := polkit.CheckAuthorization(subject, polkitManageLocale,
		details,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}

	return result.IsAuthorized, nil
}
