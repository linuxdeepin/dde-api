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

	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/sound_effect"
)

var logger = log.NewLogger("api/shutdown-sound")

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

	err = doPlayShutdownSound(cfg.Theme, cfg.Event, cfg.Device)
	if err != nil {
		logger.Error("failed to play shutdown sound:", err)
	}
}

func handleSignal() {
	var sigChan = make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Kill, os.Interrupt) //nolint
	go func() {
		sig := <-sigChan
		switch sig {
		case os.Kill, os.Interrupt:
			// Nothing to do
			logger.Info("receive signal:", sig.String())
		}
	}()
}

func doPlayShutdownSound(theme, event, device string) error {
	logger.Infof("play theme: %s, event: %s, device: %s", theme, event, device)
	player := sound_effect.NewPlayer(false, sound_effect.PlayBackendALSA)
	duration, _ := player.GetDuration(theme, event)
	logger.Info("duration:", duration)
	if duration > 0 {
		time.AfterFunc(duration, func() {
			os.Exit(0)
		})
	}

	err := player.Play(theme, event, device)
	return err
}
