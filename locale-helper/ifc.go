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

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/polkit"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	polkitManageLocale = "com.deepin.api.locale-helper.manage-locale"

	defaultLocaleFile    = "/etc/default/locale"
	defaultLocaleGenFile = "/etc/locale.gen"
)

var errAuthFailed = fmt.Errorf("authentication failed")

func (h *Helper) SetLocale(sender dbus.Sender, locale string) *dbus.Error {
	h.service.DelayAutoQuit()

	ok, err := h.checkAuth(sender)
	logger.Debug("---Auth ret:", ok, err)
	if !ok || err != nil {
		return dbusutil.ToError(errAuthFailed)
	}

	if !IsLocaleValid(locale) {
		return dbusutil.ToError(fmt.Errorf("invalid locale: %s", locale))
	}

	err = writeContentToFile(defaultLocaleFile, "LANG="+locale)
	return dbusutil.ToError(err)
}

func (h *Helper) emitFailed(err error) {
	h.service.Emit(h, "Success", false, err.Error())
}

func (h *Helper) emitRealSuccess() {
	h.service.Emit(h, "Success", true, "")
}

func (h *Helper) generateLocale(sender dbus.Sender, locale string) error {
	ok, err := h.checkAuth(sender)
	logger.Debug("---Auth ret:", ok, err)
	if !ok || err != nil {
		return errAuthFailed
	}

	if !IsLocaleValid(locale) {
		return fmt.Errorf("invalid locale: %s", locale)
	}

	// locales version <= 2.13
	if !dutils.IsFileExist(defaultLocaleGenFile) {
		err := h.doGenLocaleWithParam(locale)
		if err != nil {
			return err
		}
		return nil
	}

	err = enableLocaleInFile(locale, defaultLocaleGenFile)
	if err != nil {
		return err
	}

	err = h.doGenLocale()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) GenerateLocale(sender dbus.Sender, locale string) *dbus.Error {
	h.service.DelayAutoQuit()

	h.mu.Lock()
	h.running = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		h.running = false
		h.mu.Unlock()
	}()

	err := h.generateLocale(sender, locale)
	if err != nil {
		h.emitFailed(err)
	} else {
		h.emitRealSuccess()
	}

	return dbusutil.ToError(err)
}

func enableLocaleInFile(locale, file string) error {
	info, err := NewLocaleFileInfo(file)
	if err != nil {
		return err
	}

	if info.IsLocaleEnabled(locale) {
		return nil
	}

	info.EnableLocale(locale)
	err = info.Save(defaultLocaleGenFile)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	polkit.Init()
}

func (h *Helper) checkAuth(sender dbus.Sender) (bool, error) {
	pid, err := h.service.GetConnPID(string(sender))
	if err != nil {
		return false, err
	}

	subject := polkit.NewSubject(polkit.SubjectKindUnixProcess)
	subject.SetDetail("pid", pid)
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
