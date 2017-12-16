/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
	"time"

	"sync"

	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/sound_effect"
)

type Manager struct {
	playing bool
	mu      sync.Mutex
	player  *sound_effect.Player
}

var (
	logger = log.NewLogger("sound-theme-player")
)

func (m *Manager) Play(theme, event, device string) error {
	if theme == "" || event == "" {
		return fmt.Errorf("invalid theme or event")
	}
	go m.doPlaySound(theme, event, device)
	return nil
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.api.SoundThemePlayer",
		ObjectPath: "/com/deepin/api/SoundThemePlayer",
		Interface:  "com.deepin.api.SoundThemePlayer",
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

func main() {
	logger.Info("start sound-theme-player")
	var m = new(Manager)
	m.player = sound_effect.NewPlayer(false, sound_effect.PlayBackendALSA)

	err := dbus.InstallOnSystem(m)
	if err != nil {
		logger.Error("Install sound player bus failed:", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*10, m.canQuit)

	err = dbus.Wait()
	if err != nil {
		logger.Error("Lost system bus:", err)
	}
}
