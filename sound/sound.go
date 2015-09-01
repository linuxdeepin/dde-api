/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

// #cgo pkg-config: glib-2.0 libcanberra
// #include <stdlib.h>
// #include "sound.h"
import "C"
import "unsafe"

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
)

const (
	soundDest = "com.deepin.api.Sound"
	soundPath = "/com/deepin/api/Sound"
	soundObj  = "com.deepin.api.Sound"

	appearanceId   = "com.deepin.dde.appearance"
	gkeySoundTheme = "sound-theme"
)

type Sound struct{}

func (s *Sound) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		soundDest,
		soundPath,
		soundObj,
	}
}

// PlaySystemSound play a target event sound, such as "bell".
func (s *Sound) PlaySystemSound(event string) (err error) {
	return s.PlayThemeSound(s.getCurrentSoundTheme(), event)
}

// PlaySystemSound play a target event sound, such as "bell".
func (s *Sound) PlaySystemSoundWithDevice(event, device string) (err error) {
	return s.PlayThemeSoundWithDevice(s.getCurrentSoundTheme(), event, device)
}

func (s *Sound) getCurrentSoundTheme() string {
	var themeSettings = gio.NewSettings(appearanceId)
	defer themeSettings.Unref()

	return themeSettings.GetString(gkeySoundTheme)
}

// PlayThemeSound play a target theme's event sound.
func (s *Sound) PlayThemeSound(theme, event string) (err error) {
	return s.PlayThemeSoundWithDevice(theme, event, "")
}

// PlayThemeSound play a target theme's event sound.
func (s *Sound) PlayThemeSoundWithDevice(theme, event, device string) (err error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.doPlayThemeSound(theme, event, device)
	}()
	return
}

func (s *Sound) doPlayThemeSound(theme, event, device string) (err error) {
	ctheme := C.CString(theme)
	defer C.free(unsafe.Pointer(ctheme))
	cevent := C.CString(event)
	defer C.free(unsafe.Pointer(cevent))
	cdevice := C.CString(device)
	defer C.free(unsafe.Pointer(cdevice))
	ret := C.canberra_play_system_sound(ctheme, cevent, cdevice)
	if ret != 0 {
		err = fmt.Errorf("Play sound theme failed: theme: %s, event: %s, device: %s, error: %s",
			theme, event, device, C.GoString(C.ca_strerror(ret)))
		logger.Error(err)
	}
	return
}

// PlaySoundFile play a target sound file.
func (s *Sound) PlaySoundFile(file string) (err error) {
	return s.PlaySoundFileWithDevice(file, "")
}

// PlaySoundFile play a target sound file.
func (s *Sound) PlaySoundFileWithDevice(file, device string) (err error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.doPlaySoundFile(file, device)
	}()
	return
}

func (s *Sound) doPlaySoundFile(file, device string) (err error) {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))
	cdevice := C.CString(device)
	defer C.free(unsafe.Pointer(cdevice))
	ret := C.canberra_play_sound_file(cfile, cdevice)
	if ret != 0 {
		err = fmt.Errorf("Play sound file: %s failed: %s",
			file, C.GoString(C.ca_strerror(ret)))
		logger.Error(err)
	}
	return
}
