package main

import (
	"os"
	"path"

	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	kfGroupDesktop = "Desktop Entry"
	kfKeyName      = "Name"
	kfKeyExec      = "Exec"
)

type AppInfo struct {
	// Desktop id
	Id string
	// Display name
	Name string
	// Commandline
	Exec string
}
type AppInfos []*AppInfo

func GetAppInfo(ty string) (*AppInfo, error) {
	id, err := mime.GetDefaultApp(ty, false)
	if err != nil {
		return nil, err
	}

	return newAppInfoById(id)
}

func (infos AppInfos) Delete(id string) AppInfos {
	var ret AppInfos
	for _, info := range infos {
		if info.Id == id {
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func SetAppInfo(ty, id string) error {
	return mime.SetDefaultApp(ty, id)
}

func GetAppInfos(ty string) AppInfos {
	var infos AppInfos
	for _, id := range mime.GetAppList(ty) {
		info, err := newAppInfoById(id)
		if err != nil {
			continue
		}
		infos = append(infos, info)
	}
	return infos
}

func newAppInfoById(id string) (*AppInfo, error) {
	kfile, err := dutils.NewKeyFileFromFile(
		findFilePath(path.Join("applications", id)))
	if err != nil {
		return nil, err
	}
	defer kfile.Free()

	var info AppInfo
	info.Id = id
	info.Name, _ = kfile.GetLocaleString(kfGroupDesktop, kfKeyName, "\x00")
	info.Exec, _ = kfile.GetString(kfGroupDesktop, kfKeyExec)
	return &info, nil
}

func findFilePath(file string) string {
	data := path.Join(os.Getenv("HOME"), ".local/share", file)
	if dutils.IsFileExist(data) {
		return data
	}

	data = path.Join("/usr/local/share", file)
	if dutils.IsFileExist(data) {
		return data
	}

	return path.Join("/usr/share", file)
}
