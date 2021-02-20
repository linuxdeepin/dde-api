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
	"os"
	"os/signal"
	"time"
	"os/exec"
	"strings"
	"os/user"
	"encoding/json"
	"path/filepath"
	"fmt"
	"io/ioutil"
	"strconv"

	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/sound_effect"
)

const (
	defaultHomeDir  = "/var/lib/deepin-shut-player"
)

var (
	logger      = log.NewLogger("api/shutdown-sound")
	homeDir     string
	curUid      int
)

func init() {
	u, err := user.Current()
	if err != nil {
		logger.Warning(err)
	} else {
		_, err := os.Stat(defaultHomeDir)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.MkdirAll(defaultHomeDir, os.ModePerm)
				if err != nil {
					logger.Warning(err)
				}
			}
		}
		homeDir = defaultHomeDir
		curUid, _ = strconv.Atoi(u.Uid)
	}

	logger.Info("home:", homeDir)
	logger.Info("Uid:", curUid)
}

func main() {
	handleSignal()

	cfg, err := soundutils.GetShutdownSoundConfig()
	if err != nil {
		logger.Warning("failed to get shutdown sound config:", err)

		//加载已保存的配置，如果加载失败则退出
		cfg, err = loadExistShutdownConfig()
		if err != nil {
			logger.Warning("load shutdown config failed: ", err)
			return
		}
	} else {
		logger.Info("save shutdown config.")
		err = saveUpdateShutdownConfig(cfg)
		if err != nil {
			logger.Warning("save shutdown failed: ", err)
			return
		}
	}

	if !cfg.CanPlay {
		return
	}

	err = doPlayShutdownSound(cfg.Theme, cfg.Event, cfg.Device, cfg.Volume)
	if err != nil {
		logger.Error("failed to play shutdown sound:", err)
	} else {
		logger.Info("PlayShutdownSound ok.")
	}
}

func handleSignal() {
	var sigChan = make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Kill, os.Interrupt)
	go func() {
		sig := <-sigChan
		switch sig {
		case os.Kill, os.Interrupt:
			// Nothing to do
			logger.Info("receive signal:", sig.String())
		}
	}()
}

func doPlayShutdownSound(theme, event, device string, volume float32) error {
	logger.Infof("play theme: %s, event: %s, device: %s, volume: %f", theme, event, device, volume)
	player := sound_effect.NewPlayer(false, sound_effect.PlayBackendALSA)
	player.Volume = volume
	duration, _ := player.GetDuration(theme, event)
	logger.Info("duration:", duration)
	if duration > 0 {
		PGresult, err := checkCurrentMachineIsPangu()
		if true == PGresult && nil == err {
			logger.Debug("this is panguV.")
			time.Sleep(time.Millisecond * 3000)
		}
		time.AfterFunc(duration, func() {
			os.Exit(0)
		})
	}

	err := player.Play(theme, event, device)
	if err != nil {
		logger.Warning("shutdown play failed: ", err)
	}

	return err
}

func checkCurrentMachineIsPangu() (bool, error) {
	var value string
	var result bool = true

	out, err := exec.Command("/sbin/dmidecode", "-t", "1").CombinedOutput()
	if err != nil {
		logger.Warning("fetch dmidecode information failed: ", err)
		return false, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		kv := strings.SplitN(line, ":", 2)

		if len(kv) < 2 {
			continue
		}

		val := strings.TrimSpace(kv[1])
		switch strings.TrimSpace(kv[0]) {
		case "Product Name":
			value = val
		}
	}
	// value = string(out)
	if strings.Contains(value, "PGUV-WBY0") {
		value = "panguV"
	} else if strings.Contains(value, "PGUV-WBX0") {
		value = "panguV"
	} else if strings.Contains(value, "PGU-WBX0") {
		value = "pangu"
	} else {
		value = ""
		result = false
	}

	if true == result {
		logger.Info("print fetch machine product name: ", value)
	}

	return result, nil
}

func loadExistShutdownConfig() (*soundutils.ShutdownSoundConfig, error) {
	logger.Info("loadUserConfig check user: ", curUid)

	var cfg soundutils.ShutdownSoundConfig
	err := loadUserConfig(curUid, &cfg)
	if err != nil {
		logger.Warning("loadUserConfig failed: ", err)
		return nil, err
	}

	logger.Info("print cfg.CanPlay: ", cfg.CanPlay)

	return &cfg, nil
}

func saveUpdateShutdownConfig(cfg *soundutils.ShutdownSoundConfig) error {
	logger.Info("saveUserConfig check user: ", curUid)

	err := saveUserConfig(curUid, cfg)
	if err != nil {
		logger.Warning("save user shutdown config faild: ", err)
		return err
	}

	return nil
}

func getConfigFile(uid int) string {
	return filepath.Join(homeDir, fmt.Sprintf("config-%d.json", uid))
}

func loadUserConfig(uid int, cfg *soundutils.ShutdownSoundConfig) error {
	filename := getConfigFile(uid)
	logger.Info("check filename: ", filename)
	return loadConfig(filename, cfg)
}

func loadConfig(filename string, cfg *soundutils.ShutdownSoundConfig) error {
	cfg.Theme = "deepin"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, cfg)
	logger.Debugf("load config file %q: %#v", filename, cfg)
	return err
}


func saveUserConfig(uid int, cfg *soundutils.ShutdownSoundConfig) error {
	filename := getConfigFile(uid)
	logger.Info("check filename: ", filename)
	return saveConfig(filename, cfg)
}

func saveConfig(filename string, cfg *soundutils.ShutdownSoundConfig) error {
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
