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
	"os"
	"strings"
	"sync"
	"time"

	"github.com/godbus/dbus"
	"pkg.deepin.io/dde/api/themes"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
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
	dbusServiceName = "com.deepin.api.CursorHelper"
	dbusPath        = "/com/deepin/api/CursorHelper"
	dbusInterface   = "com.deepin.api.CursorHelper"
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
