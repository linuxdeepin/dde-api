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

package soundutils

import (
	"sync"

	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/sound_effect"
	"pkg.deepin.io/lib/strv"
)

const (
	EventPowerPlug     = "power-plug"
	EventPowerUnplug   = "power-unplug"
	EventBatteryLow    = "power-unplug-battery-low"
	EventVolumeChanged = "audio-volume-change"
	EventIconToDesktop = "x-deepin-app-sent-to-desktop"
	EventLogin         = "desktop-login"
	EventLogout        = "desktop-logout"
	EventShutdown      = "system-shutdown"
	EventWakeup        = "suspend-resume"

	EventPowerUnplugBatteryLow   = "power-unplug-battery-low"
	EventAudioVolumeChanged      = "audio-volume-change"
	EventXDeepinAppSentToDesktop = "x-deepin-app-sent-to-desktop"
	EventDesktopLogin            = "desktop-login"
	EventDesktopLogout           = "desktop-logout"
	EventSystemShutdown          = "system-shutdown"
	EventSuspendResume           = "suspend-resume"

	EventDeviceAdded   = "device-added"
	EventDeviceRemoved = "device-removed"
)

const (
	soundEffectSchema = "com.deepin.dde.sound-effect"
	appearanceSchema  = "com.deepin.dde.appearance"
	keySoundTheme     = "sound-theme"
	keyEnabled        = "enabled"
	keyPlayer         = "player"
	defaultSoundTheme = "deepin"
)

func PlaySystemSound(event, device string) error {
	return PlayThemeSound(GetSoundTheme(), event, device)
}

var UseCache = true

var player *sound_effect.Player
var playerOnce sync.Once

func initPlayer() {
	playerOnce.Do(func() {
		player = sound_effect.NewPlayer(UseCache, sound_effect.PlayBackendPulseAudio)
	})
}

func PlayThemeSound(theme, event, device string) error {
	if theme == "" {
		theme = defaultSoundTheme
	}

	if !CanPlayEvent(event) {
		return nil
	}

	initPlayer()
	return player.Play(theme, event, device)
}

func CanPlayEvent(event string) bool {
	if event == keyEnabled || event == keyPlayer {
		return false
	}

	settings := gio.NewSettings(soundEffectSchema)
	defer settings.Unref()

	// check main switch
	if !settings.GetBoolean(keyEnabled) {
		return false
	}

	keys := strv.Strv(settings.ListKeys())
	if keys.Contains(event) {
		// has key
		return settings.GetBoolean(event)
	}
	return true
}

func GetSoundTheme() string {
	s := gio.NewSettings(appearanceSchema)
	defer s.Unref()
	return s.GetString(keySoundTheme)
}

func GetThemeSoundFile(theme, event string) string {
	if theme == "" {
		theme = defaultSoundTheme
	}

	initPlayer()
	return player.Finder().Find(theme, "stereo", event)
}

func GetSystemSoundFile(event string) string {
	return GetThemeSoundFile(GetSoundTheme(), event)
}
