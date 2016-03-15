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
	player "pkg.deepin.io/lib/sound"
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
		return player.PlayThemeSound(theme, event, device, "")
	}

	go player.PlayThemeSound(theme, event, device, "")
	return nil
}

func PlaySoundFile(file, device string, sync bool) error {
	if sync {
		return player.PlaySoundFile(file, device, "")
	}

	go player.PlaySoundFile(file, device, "")
	return nil
}

func CanPlayEvent() bool {
	setting := gio.NewSettings(soundEffectSchema)
	defer setting.Unref()
	return setting.GetBoolean(keyEnabled)
}

func GetSoundTheme() string {
	s := gio.NewSettings(appearanceSchema)
	defer s.Unref()
	return s.GetString(keySoundTheme)
}
