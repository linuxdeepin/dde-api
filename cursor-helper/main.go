// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/linuxdeepin/dde-api/themes"
	"github.com/linuxdeepin/go-lib/dbusutil"
	"github.com/linuxdeepin/go-lib/log"
)

//go:generate dbusutil-gen em -type Manager

type Manager struct {
	service    *dbusutil.Service
	runningMu  sync.Mutex
	running    bool
	setThemeMu sync.Mutex
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

const (
	dbusServiceName = "org.deepin.dde.CursorHelper1"
	dbusPath        = "/org/deepin/dde/CursorHelper1"
	dbusInterface   = "org.deepin.dde.CursorHelper1"
)

var logger = log.NewLogger("api/cursor-helper")

func (m *Manager) Set(name string) *dbus.Error {
	m.service.DelayAutoQuit()

	m.runningMu.Lock()
	m.running = true
	m.runningMu.Unlock()

	go func() {
		m.setThemeMu.Lock()
		err := setTheme(name)
		if err != nil {
			logger.Warning(err)
		}
		m.setThemeMu.Unlock()

		m.runningMu.Lock()
		m.running = false
		m.runningMu.Unlock()
	}()
	return nil
}

func main() {
	var name string
	if len(os.Args) == 2 {
		tmp := strings.ToLower(os.Args[1])
		if tmp == "-h" || tmp == "--help" {
			fmt.Println("Usage: cursor-theme-helper <Cursor theme>")
			return
		}
		name = os.Args[1]
	}

	if name != "" {
		setTheme(name)
		return
	}

	// start DBus service
	service, err := dbusutil.NewSessionService()
	if err != nil {
		logger.Fatal("failed to new session service", err)
	}

	hasOwner, err := service.NameHasOwner(dbusServiceName)
	if err != nil {
		logger.Fatal("failed to call NameHasOwner:", err)
	}
	if hasOwner {
		logger.Fatalf("name %q already has the owner", dbusServiceName)
	}

	m := &Manager{
		service: service,
	}
	err = service.Export(dbusPath, m)
	if err != nil {
		logger.Fatal("failed to export:", err)
	}
	err = service.RequestName(dbusServiceName)
	if err != nil {
		logger.Fatal("failed to request name:", err)
	}
	service.SetAutoQuitHandler(time.Second*5, func() bool {
		m.runningMu.Lock()
		r := m.running
		m.runningMu.Unlock()
		return !r
	})
	service.Wait()
}

func setTheme(name string) error {
	if name == themes.GetCursorTheme() {
		return nil
	}

	err := themes.SetCursorTheme(name)
	if err != nil {
		logger.Warning("Set failed:", err)
		return err
	}
	return nil
}
