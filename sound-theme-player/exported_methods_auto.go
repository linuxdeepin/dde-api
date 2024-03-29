// Code generated by "dbusutil-gen em -type Manager"; DO NOT EDIT.

package main

import (
	"github.com/linuxdeepin/go-lib/dbusutil"
)

func (v *Manager) GetExportedMethods() dbusutil.ExportedMethods {
	return dbusutil.ExportedMethods{
		{
			Name:   "EnableSound",
			Fn:     v.EnableSound,
			InArgs: []string{"name", "enabled"},
		},
		{
			Name:   "EnableSoundDesktopLogin",
			Fn:     v.EnableSoundDesktopLogin,
			InArgs: []string{"enabled"},
		},
		{
			Name:   "Play",
			Fn:     v.Play,
			InArgs: []string{"theme", "event", "device"},
		},
		{
			Name: "PlaySoundDesktopLogin",
			Fn:   v.PlaySoundDesktopLogin,
		},
		{
			Name:   "PrepareShutdownSound",
			Fn:     v.PrepareShutdownSound,
			InArgs: []string{"uid"},
		},
		{
			Name:   "SaveAudioState",
			Fn:     v.SaveAudioState,
			InArgs: []string{"activePlayback"},
		},
		{
			Name:   "SetSoundTheme",
			Fn:     v.SetSoundTheme,
			InArgs: []string{"theme"},
		},
	}
}
