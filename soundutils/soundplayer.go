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
	EventLogin               = "sys-login"
	EventLogout              = "sys-logout"
	EventShutdown            = "sys-shutdown"
	EventWakeup              = "suspend-resume"
	EventNotification        = "message-out"
	EventUnableOperate       = "app-error"
	EventEmptyTrash          = "trash-empty"
	EventVolumeChanged       = "audio-volume-change"
	EventBatteryLow          = "power-unplug-battery-low"
	EventPowerPlug           = "power-plug"
	EventPowerUnplug         = "power-unplug"
	EventDevicePlug          = "device-added"
	EventDeviceUnplug        = "device-removed"
	EventIconToDesktop       = "send-to"
	EventCameraShutter       = "camera-shutter"
	EventScreenCapture       = "screen-capture"
	EventScreenCaptureFinish = "screen-capture-complete"
)

// map sound file name -> key in gsettings
var soundFileKeyMap = map[string]string{
	EventLogin:         "login",
	EventLogout:        "logout",
	EventShutdown:      "shutdown",
	EventWakeup:        "wakeup",
	EventNotification:  "notification",
	EventUnableOperate: "unable-operate",
	EventEmptyTrash:    "empty-trash",
	EventVolumeChanged: "volume-change",
	EventBatteryLow:    "battery-low",
	// power-plug
	// power-unplug
	EventDevicePlug:    "device-plug",
	EventDeviceUnplug:  "device-unplug",
	EventIconToDesktop: "icon-to-desktop",
	// camera-shutter
	EventScreenCapture:       "screenshot",
	EventScreenCaptureFinish: "screenshot",
}

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

	key, ok := soundFileKeyMap[event]
	if !ok {
		key = event
	}

	keys := strv.Strv(setting.ListKeys())
	if keys.Contains(key) {
		// has key
		return setting.GetBoolean(key)
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
