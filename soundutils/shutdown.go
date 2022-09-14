// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package soundutils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	shutdownSoundFile = "/tmp/deepin-shutdown-sound.json"
)

type ShutdownSoundConfig struct {
	CanPlay bool
	Theme   string
	Event   string
	Device  string
	Volume  float32
}

func GetShutdownSoundConfig() (*ShutdownSoundConfig, error) {
	data, err := ioutil.ReadFile(shutdownSoundFile)
	if err != nil {
		return nil, err
	}
	var v ShutdownSoundConfig
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func SetShutdownSoundConfig(v *ShutdownSoundConfig) (err error) {
	if !v.CanPlay {
		// 当不需要播放关机音效时，删除配置文件。
		err = os.Remove(shutdownSoundFile)
		if err != nil && os.IsNotExist(err) {
			err = nil
		}
		return
	}
	data, err := json.Marshal(v)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(shutdownSoundFile, data, 0644)
	return
}
