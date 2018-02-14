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
	"time"

	"sync"

	"errors"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/sound_effect"
)

const (
	dbusServiceName = "com.deepin.api.SoundThemePlayer"
	dbusPath        = "/com/deepin/api/SoundThemePlayer"
	dbusInterface   = dbusServiceName
)

var (
	logger = log.NewLogger("sound-theme-player")
)

type Manager struct {
	playing bool
	mu      sync.Mutex
	player  *sound_effect.Player
	service *dbusutil.Service

	methods *struct {
		Play func() `in:"theme,event,device"`
	}
}

func (m *Manager) Play(theme, event, device string) *dbus.Error {
	if theme == "" || event == "" {
		return dbusutil.ToError(errors.New("invalid theme or event"))
	}
	go m.doPlaySound(theme, event, device)
	return nil
}

func (m *Manager) GetDBusExportInfo() dbusutil.ExportInfo {
	return dbusutil.ExportInfo{
		Path:      dbusPath,
		Interface: dbusInterface,
	}
}

func (m *Manager) doPlaySound(theme, event, device string) {
	m.mu.Lock()
	m.playing = true
	m.mu.Unlock()

	err := m.player.Play(theme, event, device)

	m.mu.Lock()
	m.playing = false
	m.mu.Unlock()

	if err != nil {
		logger.Warning("failed to play:", err)
	}
}

func (m *Manager) canQuit() bool {
	m.mu.Lock()
	playing := m.playing
	m.mu.Unlock()
	return !playing
}

func newManager(service *dbusutil.Service) *Manager {
	player := sound_effect.NewPlayer(false, sound_effect.PlayBackendALSA)
	return &Manager{
		player:  player,
		service: service,
	}
}

func main() {
	logger.Info("start sound-theme-player")
	service, err := dbusutil.NewSystemService()
	if err != nil {
		logger.Fatal("failed to new system service", err)
	}

	hasOwner, err := service.NameHasOwner(dbusServiceName)
	if err != nil {
		logger.Fatal("failed to call NameHasOwner:", err)
	}
	if hasOwner {
		logger.Fatalf("name %q already has the owner", dbusServiceName)
	}

	m := newManager(service)
	err = service.Export(m)
	if err != nil {
		logger.Fatal("failed to export:", err)
	}
	err = service.RequestName(dbusServiceName)
	if err != nil {
		logger.Fatal("failed to request name:", err)
	}
	service.SetAutoQuitHandler(time.Second*10, m.canQuit)
	service.Wait()
}
