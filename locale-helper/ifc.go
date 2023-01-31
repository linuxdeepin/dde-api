// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"

	"github.com/godbus/dbus"
	polkit "github.com/linuxdeepin/go-dbus-factory/system/org.freedesktop.policykit1"
	"github.com/linuxdeepin/go-lib/dbusutil"
	dutils "github.com/linuxdeepin/go-lib/utils"
)

const (
	polkitManageLocale = "org.deepin.dde.locale-helper.manage-locale"

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
	erro := h.service.Emit(h, "Success", false, err.Error())
	if erro != nil {
		logger.Warning(erro)
	}
}

func (h *Helper) emitRealSuccess() {
	err := h.service.Emit(h, "Success", true, "")
	if err != nil {
		logger.Warning(err)
	}
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

func (h *Helper) checkAuth(sender dbus.Sender) (bool, error) {
	systemBus := h.service.Conn()
	authority := polkit.NewAuthority(systemBus)
	subject := polkit.MakeSubject(polkit.SubjectKindSystemBusName)
	subject.SetDetail("name", string(sender))
	result, err := authority.CheckAuthorization(0, subject, polkitManageLocale,
		nil,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}
	return result.IsAuthorized, nil
}
