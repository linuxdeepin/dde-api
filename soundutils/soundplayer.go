// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package soundutils

import (
	"errors"
	"sync"

	"github.com/godbus/dbus/v5"
	configManager "github.com/linuxdeepin/go-dbus-factory/org.desktopspec.ConfigManager"
	"github.com/linuxdeepin/go-lib/log"
	"github.com/linuxdeepin/go-lib/sound_effect"
	"github.com/linuxdeepin/go-lib/strv"
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
	dconfigDaemonAppId     = "org.deepin.dde.daemon"
	dconfigSoundEffectId   = "org.deepin.dde.daemon.soundeffect"
	dconfigAppearanceAppId = "org.deepin.dde.appearance"
	dconfigAppearanceId    = dconfigAppearanceAppId

	keySoundTheme     = "Sound_Theme"
	keyEnabled        = "enabled"
	keyPlayer         = "player"
	defaultSoundTheme = "deepin"
)

var logger = log.NewLogger("soundplayer")

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

// TODO: 后续此部分dconfig逻辑封装成库到go-lib中
func makeDConfigManager(appID string, id string) (configManager.Manager, error) {
	bus, err := dbus.SystemBus()
	if err != nil {
		return nil, errors.New("bus cannot be nil")
	}
	dsMgr := configManager.NewConfigManager(bus)

	if dsMgr == nil {
		return nil, errors.New("dsManager cannot be nil")
	}

	dsPath, err := dsMgr.AcquireManager(0, appID, id, "")
	if err != nil {
		logger.Warning(err)
	}

	dconfigMgr, err := configManager.NewManager(bus, dsPath)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}
	return dconfigMgr, nil
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

	soundeffectDconfig, err := makeDConfigManager(dconfigDaemonAppId, dconfigSoundEffectId)
	if err != nil {
		logger.Warning(err)
		return false
	}

	soundeffectEnabledValue, err := soundeffectDconfig.Value(0, keyEnabled)

	if err != nil {
		logger.Warning(err)
		return false
	}

	// check main switch
	if !soundeffectEnabledValue.Value().(bool) {
		return false
	}

	keyList, _ := soundeffectDconfig.KeyList().Get(0)

	keys := strv.Strv(keyList)
	if keys.Contains(event) {
		// has key
		soundEnabled, err := soundeffectDconfig.Value(0, event)
		if err != nil {
			return false
		}
		return soundEnabled.Value().(bool)
	}
	return true
}

func GetSoundTheme() string {
	appeearanceDconfig, err := makeDConfigManager(dconfigAppearanceAppId, dconfigAppearanceId)
	if err != nil {
		logger.Warning(err)
		return defaultSoundTheme
	}

	soundeffectEnabledValue, err := appeearanceDconfig.Value(0, keySoundTheme)
	if err != nil {
		logger.Warning(err)
		return defaultSoundTheme
	}
	return soundeffectEnabledValue.Value().(string)
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
