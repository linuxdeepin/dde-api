/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package soundutils

import (
	"io/ioutil"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	shutdownFile    = "/tmp/deepin-shutdown-sound.ini"
	kfGroupShutdown = "Shutdown"
	kfKeyCanPlay    = "CanPlay"
	kfKeySoundTheme = "SoundTheme"
	kfKeySoundEvent = "SoundEvent"
)

func GetShutdownSound() (bool, string, string, error) {
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

func SetShutdownSound(canPlay bool, theme, event string) error {
	if !dutils.IsFileExist(shutdownFile) {
		err := dutils.CreateFile(shutdownFile)
		if err != nil {
			return err
		}
	}

	kf, err := dutils.NewKeyFileFromFile(shutdownFile)
	if err != nil {
		return err
	}
	defer kf.Free()

	kf.SetBoolean(kfGroupShutdown, kfKeyCanPlay, canPlay)
	kf.SetString(kfGroupShutdown, kfKeySoundTheme, theme)
	kf.SetString(kfGroupShutdown, kfKeySoundEvent, event)

	_, content, err := kf.ToData()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(shutdownFile, []byte(content), 0644)
}
