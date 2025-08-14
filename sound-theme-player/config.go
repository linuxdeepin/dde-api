// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
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

func loadConfig(filename string, cfg *config) error {
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
