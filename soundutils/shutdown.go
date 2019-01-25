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

package soundutils

import (
	"encoding/json"
	"io/ioutil"
)

const (
	shutdownSoundFile = "/tmp/deepin-shutdown-sound.json"
)

type ShutdownSoundConfig struct {
	CanPlay bool
	Theme   string
	Event   string
	Device  string
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

func SetShutdownSoundConfig(v *ShutdownSoundConfig) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(shutdownSoundFile, data, 0644)
	return err
}
