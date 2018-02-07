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
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/sound"
	dutils "pkg.deepin.io/lib/utils"
)

var logger = log.NewLogger("api/shutdown-sound")

func main() {
	logger.Info("[DEEPIN SHUTDOWN SOUND] play shutdown sound")
	handleSignal()

	canPlay, theme, event, err := getShutdownSound()
	if err != nil {
		logger.Warning("[DEEPIN SHUTDOWN SOUND] get shutdown sound info failed:", err)
		return
	}
	logger.Info("[DEEPIN SHUTDOWN SOUND] can play:", canPlay, theme, event)

	if !canPlay {
		return
	}

	err = doPlayShutdwonSound(theme, event)
	if err != nil {
		logger.Error("[DEEPIN SHUTDOWN SOUND] play shutdown sound failed:", theme, event, err)
	}
}

func handleSignal() {
	var sigs = make(chan os.Signal, 2)
	signal.Notify(sigs, os.Kill, os.Interrupt)
	go func() {
		sig := <-sigs
		switch sig {
		case os.Kill, os.Interrupt:
			// Nothing to do
			logger.Info("[DEEPIN SHUTDOWN SOUND] receive signal:", sig.String())
		}
	}()
}

func doPlayShutdwonSound(theme, event string) error {
	logger.Info("[DEEPIN SHUTDOWN SOUND] do play:", theme, event)
	err := sound.PlayThemeSound(theme, event, "", "alsa", "")
	if err != nil {
		logger.Error("[DEEPIN SHUTDOWN SOUND] do play failed:", theme, event, err)
		return err
	}
	return nil
}

// fixed compile failure when soundutils api changed
const (
	shutdownFile    = "/tmp/deepin-shutdown-sound.ini"
	kfGroupShutdown = "Shutdown"
	kfKeyCanPlay    = "CanPlay"
	kfKeySoundTheme = "SoundTheme"
	kfKeySoundEvent = "SoundEvent"
)

func getShutdownSound() (bool, string, string, error) {
	kf, err := dutils.NewKeyFileFromFile(shutdownFile)
	if err != nil {
		return false, "", "", err
	}
	defer kf.Free()

	canPlay, err := kf.GetBoolean(kfGroupShutdown, kfKeyCanPlay)
	if err != nil {
		return false, "", "", err
	}

	theme, err := kf.GetString(kfGroupShutdown, kfKeySoundTheme)
	if err != nil {
		return false, "", "", err
	}

	event, err := kf.GetString(kfGroupShutdown, kfKeySoundEvent)
	if err != nil {
		return false, "", "", err
	}

	return canPlay, theme, event, nil
}
