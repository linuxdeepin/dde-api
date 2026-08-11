// SPDX-FileCopyrightText: 2018-2026 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
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

	authorized, err := h.authorizeOrPolkit(sender)
	if err != nil {
		logger.Warning("SetLocale access denied:", err)
		return dbusutil.ToError(err)
	}
	if !authorized {
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
	authorized, err := h.authorizeOrPolkit(sender)
	if err != nil {
		logger.Warning("GenerateLocale access denied:", err)
		return err
	}
	if !authorized {
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

func (h *Helper) SetAllowCaller(sender dbus.Sender, uniqueName string) *dbus.Error {
	h.service.DelayAutoQuit()
	err := h.allowCallers.addCaller(sender, uniqueName)
	if err != nil {
		logger.Warningf("SetAllowCaller rejected sender %s for target %s: %v", sender, uniqueName, err)
	}
	return dbusutil.ToError(err)
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

// authorizeOrPolkit checks the allow-caller registry first. If the registry is
// not enabled (the service was not launched through deepin-security-loader),
// fall back to the Polkit authorization dialog.
func (h *Helper) authorizeOrPolkit(sender dbus.Sender) (bool, error) {
	err := h.allowCallers.authorize(sender)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, errAllowCallerNotEnabled) {
		// Not launched via deepin-security-loader: fall back to Polkit.
		return h.checkAuth(sender)
	}
	return false, err
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
