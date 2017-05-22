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

	if !CanPlayEvent() {
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

func CanPlayEvent() bool {
	if setting == nil {
		s, err := utils.CheckAndNewGSettings(soundEffectSchema)
		if err != nil {
			return true
		}
		setting = s
	}
	return setting.GetBoolean(keyEnabled)
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
