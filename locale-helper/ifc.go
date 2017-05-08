/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/polkit"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	polkitManageLocale = "com.deepin.api.locale-helper.manage-locale"
	polkitAuthMsg      = "Authentication is required to switch language"

	defaultLocaleFile    = "/etc/default/locale"
	defaultLocaleGenFile = "/etc/locale.gen"
)

var errAuthFailed = fmt.Errorf("Authentication failed")

func (h *Helper) SetLocale(dmessage dbus.DMessage, locale string) error {
	ok, err := checkAuth(dmessage)
	logger.Debug("---Auth ret:", ok, err)
	if !ok || err != nil {
		return errAuthFailed
	}

	if !IsLocaleValid(locale) {
		return fmt.Errorf("Invalid locale: %s", locale)
	}

	return writeContentToFile(defaultLocaleFile,
		fmt.Sprintf("LANG=%s", locale))
}

func (h *Helper) GenerateLocale(dmessage dbus.DMessage, locale string) error {
	ok, err := checkAuth(dmessage)
	logger.Debug("---Auth ret:", ok, err)
	if !ok || err != nil {
		return errAuthFailed
	}

	h.running = true
	defer func() {
		h.running = false
	}()

	if !IsLocaleValid(locale) {
		dbus.Emit(h, "Success", false,
			fmt.Sprintf("Invalid locale: %v", locale))
		return fmt.Errorf("Invalid locale: %s", locale)
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
	details["polkit.gettext_domain"] = "dde-daemon"
	details["polkit.message"] = Tr(polkitAuthMsg)
	result, err := polkit.CheckAuthorization(subject, polkitManageLocale,
		details,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}

	return result.IsAuthorized, nil
}
