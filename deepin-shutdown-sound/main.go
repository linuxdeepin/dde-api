// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/linuxdeepin/dde-api/soundutils"
	"github.com/linuxdeepin/go-lib/log"
	"github.com/linuxdeepin/go-lib/sound_effect"
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

	err = doPlayShutdownSound(cfg.Theme, cfg.Event, cfg.Device, cfg.Volume)
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

func doPlayShutdownSound(theme, event, device string, volume float32) error {
	logger.Infof("play theme: %s, event: %s, device: %s,volume %f", theme, event, device, volume)
	player := sound_effect.NewPlayer(false, sound_effect.PlayBackendALSA)
	player.Volume = volume
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
