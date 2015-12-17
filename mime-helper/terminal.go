package main

import (
	"fmt"
	"path"
	"strings"

	"gir/gio-2.0"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	terminalSchema = "com.deepin.desktop.default-applications.terminal"
	gsKeyExec      = "exec"

	kfKeyCategory    = "Categories"
	cateKeyTerminal  = "TerminalEmulator"
	execKeyXTerminal = "x-terminal-emulator"
)

var termBlackList = []string{
	"guake.desktop",
}

func resetTerminal() {
	s := gio.NewSettings(terminalSchema)
	defer s.Unref()

	s.Reset(gsKeyExec)
}

func setDefaultTerminal(id string) error {
	s := gio.NewSettings(terminalSchema)
	defer s.Unref()

	for _, info := range getTerminalInfos() {
		if info.Id == id {
			s.SetString(gsKeyExec, strings.Split(info.Exec, " ")[0])
			return nil
		}
	}
	return fmt.Errorf("Invalid terminal id '%s'", id)
}

func getDefaultTerminal() (*AppInfo, error) {
	s := gio.NewSettings(terminalSchema)
	defer s.Unref()

	exec := s.GetString(gsKeyExec)
	for _, info := range getTerminalInfos() {
		if exec == strings.Split(info.Exec, " ")[0] {
			return info, nil
		}
	}

	return nil, fmt.Errorf("Not found app id for '%s'", exec)
}

func getTerminalInfos() AppInfos {
	infos := gio.AppInfoGetAll()
	defer unrefAppInfos(infos)

	var list AppInfos
	for _, info := range infos {
		if !isTerminalApp(info.GetId()) {
			continue
		}

		list = append(list, &AppInfo{
			Id:   info.GetId(),
			Name: info.GetDisplayName(),
			Exec: info.GetCommandline(),
		})
	}
	return list
}

func isTerminalApp(id string) bool {
	kfile, err := dutils.NewKeyFileFromFile(
		findFilePath(path.Join("applications", id)))
	if err != nil {
		return false
	}
	defer kfile.Free()

	cates, _ := kfile.GetString(kfGroupDesktop, kfKeyCategory)
	if err != nil {
		return false
	}

	if !strings.Contains(cates, cateKeyTerminal) {
		return false
	}

	exec, _ := kfile.GetString(kfGroupDesktop, kfKeyExec)
	if strings.Contains(exec, execKeyXTerminal) {
		return false
	}

	return (isStrInList(id, termBlackList) == false)
}

func isStrInList(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}

	return false
}

func unrefAppInfos(infos []*gio.AppInfo) {
	for _, info := range infos {
		info.Unref()
	}
}
