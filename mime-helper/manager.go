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
	"path"

	"gir/glib-2.0"
	"pkg.deepin.io/lib/dbus"
	dutils "pkg.deepin.io/lib/utils"
	"sync"
)

const (
	AppMimeTerminal = "application/x-terminal"
)

const (
	dbusDest = "com.deepin.api.Mime"
	dbusPath = "/com/deepin/api/Manager"
	dbusIFC  = "com.deepin.api.Manager"
)

const (
	stateResetStart int = iota + 1
	stateResetFinished
)

type Manager struct {
	Change func()

	media *Media

	userManager *userAppManager
	stateLocker sync.Mutex
	resetState  int
}

func NewManager() *Manager {
	m := new(Manager)
	m.resetState = stateResetFinished

	userManager, err := newUserAppManager(userAppFile)
	if err != nil {
		userManager = &userAppManager{
			filename: userAppFile,
		}
	}
	m.userManager = userManager
	return m
}

func (m *Manager) initConfigData() {
	if dutils.IsFileExist(path.Join(glib.GetUserConfigDir(),
		"mimeapps.list")) {
		return
	}

	err := m.doInitConfigData()
	if err != nil {
		logger.Warning("Init mime config file failed", err)
	}
}

func (m *Manager) doInitConfigData() error {
	return genMimeAppsFile(
		findFilePath(path.Join("dde-api", "mime", "data.json")))
}

// Reset reset mimes default app
func (m *Manager) Reset() {
	m.stateLocker.Lock()
	if m.resetState == stateResetStart {
		m.stateLocker.Unlock()
		return
	}

	m.resetState = stateResetStart
	go func() {
		err := m.doInitConfigData()
		if err != nil {
			logger.Warning("Init mime config file failed", err)
		}
		m.resetState = stateResetFinished
		m.stateLocker.Unlock()
		dbus.Emit(m, "Change")
	}()

	resetTerminal()
}

// GetDefaultApp get the default app id for the special mime
// ty: the special mime
// ret0: the default app info
// ret1: error message
func (m *Manager) GetDefaultApp(ty string) (string, error) {
	var (
		info *AppInfo
		err  error
	)
	if ty == AppMimeTerminal {
		info, err = getDefaultTerminal()
	} else {
		info, err = GetAppInfo(ty)
	}
	if err != nil {
		return "", err
	}

	return marshal(info)
}

// SetDefaultApp set the default app for the special mime
// ty: the special mime
// deskId: the default app desktop id
// ret0: error message
func (m *Manager) SetDefaultApp(ty string, deskId string) error {
	if ty == AppMimeTerminal {
		return setDefaultTerminal(deskId)
	}
	return SetAppInfo(ty, deskId)
}

// ListApps list the apps that supported the special mime
// ty: the special mime
// ret0: the app infos
func (m *Manager) ListApps(ty string) string {
	var infos AppInfos
	if ty == AppMimeTerminal {
		infos = getTerminalInfos()
	} else {
		infos = GetAppInfos(ty)
	}
	content, _ := marshal(infos)
	return content
}

func (m *Manager) ListUserApp(ty string) string {
	apps := m.userManager.Get(ty)
	if len(apps) == 0 {
		return ""
	}
	var infos AppInfos
	for _, app := range apps {
		info, err := newAppInfoById(app.DesktopId)
		if err != nil {
			logger.Warningf("New '%s' failed: %v", app.DesktopId, err)
			continue
		}
		infos = append(infos, info)
	}
	content, _ := marshal(infos)
	return content
}

func (m *Manager) AddUserApp(desktopId string, mimes []string) error {
	// check app validity
	_, err := newAppInfoById(desktopId)
	if err != nil {
		logger.Error("Invalid desktop id:", desktopId)
		return err
	}
	if !m.userManager.Add(desktopId, mimes) {
		return nil
	}
	return m.userManager.Write()
}

func (m *Manager) DeleteUserApp(desktopId string) error {
	err := m.userManager.Delete(desktopId)
	if err != nil {
		logger.Errorf("Delete '%s' failed: %v", desktopId, err)
		return err
	}
	return m.userManager.Write()
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}
