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
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/godbus/dbus"
	"github.com/linuxdeepin/dde-api/soundutils"
	"github.com/linuxdeepin/go-lib/dbusutil"
	"github.com/linuxdeepin/go-lib/log"
	"github.com/linuxdeepin/go-lib/sound_effect"
	"github.com/linuxdeepin/go-lib/strv"
)

//go:generate dbusutil-gen em -type Manager

const (
	dbusServiceName = "com.deepin.api.SoundThemePlayer"
	dbusPath        = "/com/deepin/api/SoundThemePlayer"
	dbusInterface   = dbusServiceName
	defaultHomeDir  = "/var/lib/deepin-sound-player"
	alsaCtlBin      = "/usr/sbin/alsactl"
)

var (
	logger      = log.NewLogger("sound-theme-player")
	optAutoQuit bool
	homeDir     string
)

func init() {
	flag.BoolVar(&optAutoQuit, "auto-quit", true, "auto quit")
	u, err := user.Current()
	if err != nil {
		logger.Warning(err)
	} else {
		homeDir = u.HomeDir
	}
	if homeDir == "" {
		homeDir = defaultHomeDir
	}
	logger.Debug("home:", homeDir)
}

type Manager struct {
	playing bool
	mu      sync.Mutex
	player  *sound_effect.Player
	service *dbusutil.Service

	configCache map[int]*config
}

func (m *Manager) PlaySoundDesktopLogin(sender dbus.Sender) *dbus.Error {
	m.service.DelayAutoQuit()
	autoLoginUser, err := getLightDMAutoLoginUser()
	if err != nil {
		logger.Warning(err)
	}
	if autoLoginUser != "" {
		logger.Debug("autoLoginUser is not empty")
		return nil
	}

	uid, err := getLastUser()
	if err != nil {
		return dbusutil.ToError(err)
	}

	var cfg config
	err = loadUserConfig(uid, &cfg)
	if err != nil && !os.IsNotExist(err) {
		logger.Warning(err)
	}

	if cfg.Enabled && cfg.DesktopLoginEnabled && !cfg.Mute {
		err = runALSARestore(uid)
		if err != nil && !os.IsNotExist(err) {
			logger.Warning("failed to restore ALSA state:", err)
			return dbusutil.ToError(err)
		}

		device := "default"
		if cfg.Card != "" && cfg.Device != "" {
			device = fmt.Sprintf("plughw:CARD=%s,DEV=%s", cfg.Card, cfg.Device)
		}
		go func() {
			m.player.Volume = cfg.Volume
			logger.Info("volume: ", m.player.Volume)
			m.doPlaySound(cfg.Theme, "desktop-login", device)
			os.Exit(0)
		}()
	}
	return nil
}

func (m *Manager) Play(theme, event, device string) *dbus.Error {
	m.service.DelayAutoQuit()

	if theme == "" || event == "" {
		return dbusutil.ToError(errors.New("invalid theme or event"))
	}
	uid, err := getLastUser()
	if err != nil {
		return dbusutil.ToError(err)
	}
	var cfg config
	err = loadUserConfig(int(uid), &cfg)
	if err != nil && !os.IsNotExist(err) {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}

	go func() {
		m.player.Volume = cfg.Volume
		logger.Info("volume: ", m.player.Volume)
		m.doPlaySound(theme, event, device)
		os.Exit(0)
	}()
	return nil
}

// PrepareShutdownSound 准备播放关闭音效需要的配置文件，用于 greeter UI 界面调用。
// uid 参数是 greeter UI 界面中准备登录的用户的 uid。
func (m *Manager) PrepareShutdownSound(uid int) *dbus.Error {
	logger.Debug("PrepareShutdownSound", uid)
	m.service.DelayAutoQuit()

	var cfg config
	err := loadUserConfig(uid, &cfg)
	if err != nil && !os.IsNotExist(err) {
		logger.Warning(err)
	}
	shutdownCfg := &soundutils.ShutdownSoundConfig{}
	if cfg.Enabled && cfg.SystemShutdownEnabled && !cfg.Mute {
		err = runALSARestore(uid)
		if err != nil && !os.IsNotExist(err) {
			logger.Warning("failed to restore ALSA state:", err)
			return dbusutil.ToError(err)
		}

		device := "default"
		if cfg.Card != "" && cfg.Device != "" {
			device = fmt.Sprintf("plughw:CARD=%s,DEV=%s", cfg.Card, cfg.Device)
		}
		shutdownCfg.CanPlay = true
		shutdownCfg.Theme = cfg.Theme
		shutdownCfg.Event = soundutils.EventSystemShutdown
		shutdownCfg.Device = device
		shutdownCfg.Volume = cfg.Volume
	}
	logger.Debugf("set shutdown sound config %#v", shutdownCfg)
	err = soundutils.SetShutdownSoundConfig(shutdownCfg)
	return dbusutil.ToError(err)
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) doPlaySound(theme, event, device string) {
	m.mu.Lock()
	m.playing = true
	m.mu.Unlock()

	logger.Debug("doPlaySound", theme, event, device)
	err := m.player.Play(theme, event, device)

	m.mu.Lock()
	m.playing = false
	m.mu.Unlock()

	if err != nil {
		logger.Warning("failed to play:", err)
	}
}

func (m *Manager) saveAudioState(uid int, activePlayback map[string]dbus.Variant) error {
	cfg := m.getUserConfig(uid)

	var ok bool
	cfg.Card, ok = activePlayback["card"].Value().(string)
	if !ok {
		return errors.New("type of field card is not string")
	}
	cfg.Device, ok = activePlayback["device"].Value().(string)
	if !ok {
		return errors.New("type of field device is not string")
	}
	cfg.Mute, ok = activePlayback["mute"].Value().(bool)
	if !ok {
		return errors.New("type of field mute is not bool")
	}

	var volume64 float64
	volume64, ok = activePlayback["volume"].Value().(float64) //防止dde-daemon 与dde-api 之间的版本差异导致问题
	if !ok {
		logger.Warning("There is no volume type")
		volume64 = 50.0
	}
	cfg.Volume = float32(volume64)

	err := m.saveUserConfig(uid)
	if err != nil {
		return err
	}

	err = runAlsaCtlStore(uid)
	return err
}

func (m *Manager) SaveAudioState(sender dbus.Sender,
	activePlayback map[string]dbus.Variant) *dbus.Error {
	m.service.DelayAutoQuit()

	uid, err := m.service.GetConnUID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	err = m.saveAudioState(int(uid), activePlayback)
	return dbusutil.ToError(err)
}

func (m *Manager) getUserConfig(uid int) *config {
	m.mu.Lock()
	defer m.mu.Unlock()

	cfg, ok := m.configCache[uid]
	if ok {
		return cfg
	}
	var cfg0 config
	err := loadUserConfig(int(uid), &cfg0)
	if err != nil && !os.IsNotExist(err) {
		logger.Warning(err)
	}
	m.configCache[uid] = &cfg0
	return &cfg0
}

func (m *Manager) saveUserConfig(uid int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cfg, ok := m.configCache[uid]
	if !ok {
		logger.Warningf("config for uid %d not loaded", uid)
		return nil
	}
	return saveUserConfig(uid, cfg)
}

func (m *Manager) EnableSoundDesktopLogin(sender dbus.Sender, enabled bool) *dbus.Error {
	uid, err := m.service.GetConnUID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	err = m.enableSoundDesktopLogin(int(uid), enabled)
	return dbusutil.ToError(err)
}

func (m *Manager) enableSoundDesktopLogin(uid int, enabled bool) error {
	cfg := m.getUserConfig(uid)
	if cfg.DesktopLoginEnabled == enabled {
		return nil
	}

	cfg.DesktopLoginEnabled = enabled
	return m.saveUserConfig(uid)
}

func (m *Manager) EnableSound(sender dbus.Sender, name string, enabled bool) *dbus.Error {
	uid, err := m.service.GetConnUID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	err = m.enableSound(int(uid), name, enabled)
	if err != nil {
		return dbusutil.ToError(err)
	}

	return m.PrepareShutdownSound(int(uid))
}

func (m *Manager) enableSound(uid int, name string, enabled bool) error {
	cfg := m.getUserConfig(uid)
	switch name {
	case "":
		if cfg.Enabled == enabled {
			return nil
		}
		cfg.Enabled = enabled

	case soundutils.EventDesktopLogin:
		if cfg.DesktopLoginEnabled == enabled {
			return nil
		}
		cfg.DesktopLoginEnabled = enabled
	case soundutils.EventSystemShutdown:
		if cfg.SystemShutdownEnabled == enabled {
			return nil
		}
		cfg.SystemShutdownEnabled = enabled
	default:
		return fmt.Errorf("invalid sound name: %s", name)
	}
	return m.saveUserConfig(uid)
}

func (m *Manager) SetSoundTheme(sender dbus.Sender, theme string) *dbus.Error {
	uid, err := m.service.GetConnUID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	err = m.setSoundTheme(int(uid), theme)
	return dbusutil.ToError(err)
}

func (m *Manager) setSoundTheme(uid int, theme string) error {
	cfg := m.getUserConfig(uid)
	if cfg.Theme == theme {
		return nil
	}
	cfg.Theme = theme
	return m.saveUserConfig(uid)
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
		player:      player,
		service:     service,
		configCache: make(map[int]*config),
	}
}

func main() {
	flag.Parse()
	logger.Debug("start sound-theme-player")
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

	// 实际运行时才从 gsettings 加载默认设置，测试环境下没安装相关 gsettings schema 会导致崩溃。
	_loadDefaultCfgFromGSettings = true
	m := newManager(service)
	err = service.Export(dbusPath, m)
	if err != nil {
		logger.Fatal("failed to export:", err)
	}
	err = service.RequestName(dbusServiceName)
	if err != nil {
		logger.Fatal("failed to request name:", err)
	}
	if optAutoQuit {
		service.SetAutoQuitHandler(time.Second*2, m.canQuit)
	}

	time.AfterFunc(8*time.Second, func() {
		err := cleanUpConfig()
		if err != nil {
			logger.Warning(err)
		}
	})

	service.Wait()
}

// clean up redundant configuration
func cleanUpConfig() error {
	fileInfos, err := ioutil.ReadDir(homeDir)
	if err != nil {
		return err
	}

	regAsoundState, err := regexp.Compile(`asound-state-(\d)\.gz`)
	if err != nil {
		return err
	}

	regConfig, err := regexp.Compile(`config-(\d+)\.json`)
	if err != nil {
		return err
	}

	var uidStrv strv.Strv
	for _, fileInfo := range fileInfos {
		match := regAsoundState.FindStringSubmatch(fileInfo.Name())
		if match != nil {
			uidStrv, _ = uidStrv.Add(match[1])
			continue
		}

		match = regConfig.FindStringSubmatch(fileInfo.Name())
		if match != nil {
			uidStrv, _ = uidStrv.Add(match[1])
			continue
		}
	}
	logger.Debug("cleanupConfig uidStrv:", uidStrv)

	for _, uid := range uidStrv {
		_, err = user.LookupId(uid)
		if err == nil {
			// uid is ok, skip
			continue
		}

		uidInt, err := strconv.Atoi(uid)
		if err != nil {
			logger.Warning(err)
			continue
		}

		logger.Debug("clean up config for uid:", uid)

		err = os.Remove(getAsoundStateFile(uidInt))
		if err != nil && !os.IsNotExist(err) {
			logger.Warning(err)
		}

		err = os.Remove(getConfigFile(uidInt))
		if err != nil && !os.IsNotExist(err) {
			logger.Warning(err)
		}
	}
	return nil
}
