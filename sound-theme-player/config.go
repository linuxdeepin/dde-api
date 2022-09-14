// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/linuxdeepin/dde-api/soundutils"
	"github.com/linuxdeepin/go-gir/gio-2.0"
)

type config struct {
	Enabled               bool // 音效总开关
	DesktopLoginEnabled   bool
	SystemShutdownEnabled bool
	Theme                 string
	Card                  string
	Device                string
	Mute                  bool
	Volume                float32
}

func getConfigFile(uid int) string {
	return filepath.Join(homeDir, fmt.Sprintf("config-%d.json", uid))
}

func loadUserConfig(uid int, cfg *config) error {
	filename := getConfigFile(uid)
	return loadConfig(filename, cfg)
}

func saveUserConfig(uid int, cfg *config) error {
	filename := getConfigFile(uid)
	return saveConfig(filename, cfg)
}

var _loadDefaultCfgFromGSettings bool = false

func loadConfig(filename string, cfg *config) error {
	if _loadDefaultCfgFromGSettings {
		// 从 gsettings 获取默认值
		soundEffectGs := gio.NewSettings("com.deepin.dde.sound-effect")
		defer soundEffectGs.Unref()
		appearanceGs := gio.NewSettings("com.deepin.dde.appearance")
		defer appearanceGs.Unref()

		cfg.Enabled = soundEffectGs.GetBoolean("enabled")
		cfg.DesktopLoginEnabled = soundEffectGs.GetBoolean(soundutils.EventDesktopLogin)
		cfg.SystemShutdownEnabled = soundEffectGs.GetBoolean(soundutils.EventSystemShutdown)
		cfg.Theme = appearanceGs.GetString("sound-theme")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, cfg)
	logger.Debugf("load config file %q: %#v", filename, cfg)
	return err
}

func saveConfig(filename string, cfg *config) error {
	logger.Debugf("save config file %q: %#v", filename, cfg)
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
