/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package soundutils

import (
	"gir/gio-2.0"
	splayer "pkg.deepin.io/lib/sound"
	"pkg.deepin.io/lib/strv"
	"pkg.deepin.io/lib/utils"
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
	soundThemeDeepin  = "deepin"
)

func PlaySystemSound(event, device string, sync bool) error {
	return PlayThemeSound(GetSoundTheme(), event, device, sync)
}

func PlayThemeSound(theme, event, device string, sync bool) error {
	if len(theme) == 0 {
		theme = soundThemeDeepin
	}

	if !CanPlayEvent(event) {
		return nil
	}

	if sync {
		return splayer.PlayThemeSound(theme, event, device, "", GetSoundPlayer())
	}

	go splayer.PlayThemeSound(theme, event, device, "", GetSoundPlayer())
	return nil
}

func PlaySoundFile(file, device string, sync bool) error {
	if sync {
		return splayer.PlaySoundFile(file, device, "", GetSoundPlayer())
	}

	go splayer.PlaySoundFile(file, device, "", GetSoundPlayer())
	return nil
}

var setting *gio.Settings

func CanPlayEvent(event string) bool {
	if event == keyEnabled || event == keyPlayer {
		return false
	}

	if setting == nil {
		s, err := utils.CheckAndNewGSettings(soundEffectSchema)
		if err != nil {
			return true
		}
		setting = s
	}

	// check main switch
	if !setting.GetBoolean(keyEnabled) {
		return false
	}

	keys := strv.Strv(setting.ListKeys())
	if keys.Contains(event) {
		// has key
		return setting.GetBoolean(event)
	}
	return true
}

func GetSoundPlayer() string {
	if setting == nil {
		s, err := utils.CheckAndNewGSettings(soundEffectSchema)
		if err != nil {
			return ""
		}
		setting = s
	}
	return setting.GetString(keyPlayer)
}

func GetSoundTheme() string {
	s := gio.NewSettings(appearanceSchema)
	defer s.Unref()
	return s.GetString(keySoundTheme)
}
