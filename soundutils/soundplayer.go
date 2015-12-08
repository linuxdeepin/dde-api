package soundutils

import (
	"pkg.deepin.io/lib/gio-2.0"
	player "pkg.deepin.io/lib/sound"
)

const (
	EventLogin         = "sys-login"
	EventLogout        = "sys-logout"
	EventShutdown      = "sys-shutdown"
	EventWakeup        = "suspend-resume"
	EventNotification  = "message-out"
	EventUnableOperate = "app-error-critical"
	EventEmptyTrash    = "trash-empty"
	EventVolumeChanged = "audio-volume-change"
	EventBatteryLow    = "power-unplug-battery-low"
	EventPowerPlug     = "power-plug"
	EventPowerUnplug   = "power-unplug"
	EventDevicePlug    = "device-added"
	EventDeviceUnplug  = "device-removed"
	EventIconToDesktop = "send-to"
	EventScreenshot    = "screen-capture"
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

// deepin sound theme 'event - key' map
var soundEventMap = map[string]string{
	EventLogin:         KeyLogin,
	EventLogout:        KeyLogout,
	EventShutdown:      KeyShutdown,
	EventWakeup:        KeyWakeup,
	EventNotification:  KeyNotification,
	EventUnableOperate: KeyUnableOperate,
	EventEmptyTrash:    KeyEmptyTrash,
	EventVolumeChanged: KeyVolumeChange,
	EventBatteryLow:    KeyBatteryLow,
	EventPowerPlug:     KeyPowerPlug,
	EventPowerUnplug:   KeyPowerUnplug,
	EventDevicePlug:    KeyDevicePlug,
	EventDeviceUnplug:  KeyDeviceUnplug,
	EventIconToDesktop: KeyIconToDesktop,
	EventScreenshot:    KeyScreenshot,
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
	key, ok := soundEventMap[event]
	if !ok {
		return true
	}

	s := gio.NewSettings(soundEffectSchema)
	defer s.Unref()
	if !isItemInList(key, s.ListKeys()) {
		return true
	}

	return s.GetBoolean(key)
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
