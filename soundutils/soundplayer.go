package soundutils

import (
	"pkg.deepin.io/lib/gio-2.0"
	player "pkg.deepin.io/lib/sound"
)

const (
	KeyLogin         = "login"
	KeyShutdown      = "shutdown"
	KeyLogout        = "logout"
	KeyWakeup        = "wakeup"
	KeyNotification  = "notification"
	KeyUnableOperate = "unable-operate"
	KeyEmptyTrash    = "empty-trash"
	KeyVolumeChange  = "volume-change"
	KeyBatteryLow    = "battery-low"
	KeyPowerPlug     = "power-plug"
	KeyPowerUnplug   = "power-unplug"
	KeyDevicePlug    = "device-plug"
	KeyDeviceUnplug  = "device-unplug"
	KeyIconToDesktop = "icon-to-desktop"
	KeyScreenshot    = "screenshot"
)

const (
	soundEffectSchema = "com.deepin.dde.sound-effect"
	appearanceSchema  = "com.deepin.dde.appearance"
	keySoundTheme     = "sound-theme"
	soundThemeDeepin  = "deepin"
)

// deepin sound theme 'key - event' map
var soundEventMap = map[string]string{
	KeyLogin:         "sys-login",
	KeyShutdown:      "sys-shutdown",
	KeyLogout:        "sys-logout",
	KeyWakeup:        "suspend-resume",
	KeyNotification:  "message-out",
	KeyUnableOperate: "app-error-critical",
	KeyEmptyTrash:    "trash-empty",
	KeyVolumeChange:  "audio-volume-change",
	KeyBatteryLow:    "power-unplug-battery-low",
	KeyPowerPlug:     "power-plug",
	KeyPowerUnplug:   "power-unplug",
	KeyDevicePlug:    "device-added",
	KeyDeviceUnplug:  "device-removed",
	KeyIconToDesktop: "send-to",
	KeyScreenshot:    "screen-capture",
}

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
	event = QueryEvent(event)

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

func CanPlayEvent(event string) bool {
	s := gio.NewSettings(soundEffectSchema)
	defer s.Unref()
	if !isItemInList(event, s.ListKeys()) {
		return true
	}

	return s.GetBoolean(event)
}

func QueryEvent(key string) string {
	if GetSoundTheme() != soundThemeDeepin {
		return key
	}

	value, ok := soundEventMap[key]
	if !ok {
		return key
	}
	return value
}

func GetSoundTheme() string {
	s := gio.NewSettings(appearanceSchema)
	defer s.Unref()
	return s.GetString(keySoundTheme)
}

func isItemInList(item string, list []string) bool {
	for _, v := range list {
		if item == v {
			return true
		}
	}
	return false
}
