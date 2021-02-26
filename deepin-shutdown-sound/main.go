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

	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/sound_effect"
)

var (
	logger      = log.NewLogger("api/shutdown-sound")
)

func main() {
	handleSignal()

	cfg, err := soundutils.GetShutdownSoundConfig()
	if err != nil {
		logger.Warning("failed to get shutdown sound config:", err)
		return
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
