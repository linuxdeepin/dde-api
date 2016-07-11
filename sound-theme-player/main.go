/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/sound"
	"time"
)

type Manager struct{}

var (
	playing bool
	logger  = log.NewLogger("sound-theme-player")
)

func (*Manager) Play(theme, event string) error {
	if len(theme) == 0 || len(event) == 0 {
		return fmt.Errorf("Invalid theme or event")
	}
	go doPlaySound(theme, event)
	return nil
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.api.SoundThemePlayer",
		ObjectPath: "/com/deepin/api/SoundThemePlayer",
		Interface:  "com.deepin.api.SoundThemePlayer",
	}
}

func doPlaySound(theme, event string) error {
	playing = true
	defer func() {
		playing = false
	}()

	out, err := exec.Command("/usr/bin/pulseaudio", "--start").CombinedOutput()
	logger.Debug("Launch pulseaudio output ", string(out))
	if err != nil {
		logger.Warningf("Launch pulseaudio failed: error: %v, output: %v", err, string(out))
	}

	err = sound.PlayThemeSound(theme, event, "", "pulse")
	if err != nil {
		logger.Errorf("Play '%s' '%s' failed: %v", theme, event, err)
	}
	return err
}

func (*Manager) Quit() {
	logger.Info("Quit")
	out, err := exec.Command("/usr/bin/pulseaudio", "--kill").CombinedOutput()
	if err != nil {
		logger.Warningf("quit pulseaudio failed: error: %v, output: %v", err, string(out))
	}
	os.Exit(0)
}

func main() {
	logger.Info("^^^^^^^^^^^^^^^^^^^Start sound player")
	if len(os.Args) == 3 {
		logger.Info("^^^^^^^^^^^^^^^^^Play cmd:", os.Args)
		doPlaySound(os.Args[1], os.Args[2])
		return
	}

	var m = new(Manager)
	err := dbus.InstallOnSystem(m)
	if err != nil {
		logger.Error("Install sound player bus failed:", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*1, func() bool {
		if playing {
			return false
		}
		return true
	})

	err = dbus.Wait()
	if err != nil {
		logger.Error("Lost system bus:", err)
	}
}
