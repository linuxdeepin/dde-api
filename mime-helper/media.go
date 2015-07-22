package main

import (
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gio-2.0"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	mediaPath = "/com/deepin/api/Media"
	mediaIFC  = "com.deepin.api.Media"
)

const (
	mediaSchema     = "org.gnome.desktop.media-handling"
	gsKeyIgnore     = "autorun-x-content-ignore"
	gsKeyOpenFolder = "autorun-x-content-open-folder"
	gsKeyStartSoft  = "autorun-x-content-start-app"
)

const (
	nautilusRunSoft = "nautilus-autorun-software.desktop"
)

type Media struct {
	setting *gio.Settings
}

type appIdType string

var (
	idIgnore     appIdType = "ignore"
	idOpenFolder appIdType = "open-folder"
	idRunSoft    appIdType = "run-soft"
)

func (id appIdType) GetAppInfo() *AppInfo {
	switch id {
	case idIgnore:
		return &AppInfo{
			Id:   string(id),
			Name: Tr("Nothing"),
			Exec: "",
		}
	case idOpenFolder:
		return &AppInfo{
			Id:   string(id),
			Name: Tr("Open Folder"),
			Exec: "",
		}
	case idRunSoft:
		return &AppInfo{
			Id:   string(id),
			Name: Tr("Run Software"),
			Exec: "",
		}
	}
	return nil
}

func NewMedia() (*Media, error) {
	s, err := dutils.CheckAndNewGSettings(mediaSchema)
	if err != nil {
		return nil, err
	}

	var media = Media{
		setting: s,
	}
	return &media, nil
}

func (media *Media) destroy() {
	media.setting.Unref()
}

// Reset reset media mime action
func (media *Media) Reset() {
	media.setting.Reset(gsKeyIgnore)
	media.setting.Reset(gsKeyOpenFolder)
	media.setting.Reset(gsKeyStartSoft)
}

// GetDefaultApp get the default app id for the media mime
// ty: the media mime
// ret0: the default media action or app desktop id
// ret1: error message
func (media *Media) GetDefaultApp(ty string) (string, error) {
	var info *AppInfo
	switch {
	case media.isInIgnoreList(ty):
		info = idIgnore.GetAppInfo()
	case media.isInOpenList(ty):
		info = idOpenFolder.GetAppInfo()
	case media.isInRunList(ty):
		info = idRunSoft.GetAppInfo()
	default:
		var err error
		info, err = GetAppInfo(ty)
		if err != nil {
			return "", err
		}
	}

	return marshal(info)
}

// SetDefaultApp set the default app for the media mime
// ty: the media mime
// deskId: the default media action or app desktop id
// ret0: error message
func (media *Media) SetDefaultApp(ty, id string) error {
	switch appIdType(id) {
	case idIgnore:
		return media.addToIgnoreList(ty)
	case idOpenFolder:
		return media.addToOpenList(ty)
	case idRunSoft:
		return media.addToRunList(ty)
	}

	return SetAppInfo(ty, id)
}

// ListApps list the apps that supported the media mime
// ty: the media mime
// ret0: the app desktop id list and media actions
func (media *Media) ListApps(ty string) string {
	infos := GetAppInfos(ty)
	// idRunSoft == nautilusRunSoft
	infos = infos.Delete(nautilusRunSoft)

	infos = append(infos, AppInfos{
		idIgnore.GetAppInfo(),
		idOpenFolder.GetAppInfo(),
		idRunSoft.GetAppInfo(),
	}...)

	content, _ := marshal(infos)
	return content
}

func (media *Media) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: mediaPath,
		Interface:  mediaIFC,
	}
}

func (media *Media) isInIgnoreList(ty string) bool {
	list := media.setting.GetStrv(gsKeyIgnore)
	return isStrInList(ty, list)
}

func (media *Media) isInOpenList(ty string) bool {
	list := media.setting.GetStrv(gsKeyOpenFolder)
	return isStrInList(ty, list)
}

func (media *Media) isInRunList(ty string) bool {
	list := media.setting.GetStrv(gsKeyStartSoft)
	return isStrInList(ty, list)
}

func (media *Media) addToIgnoreList(ty string) error {
	if media.isInIgnoreList(ty) {
		return nil
	}
	list := media.setting.GetStrv(gsKeyIgnore)
	list = append(list, ty)
	media.setting.SetStrv(gsKeyIgnore, list)

	media.delFromOpenList(ty)
	media.delFromRunList(ty)
	return nil
}

func (media *Media) addToOpenList(ty string) error {
	if media.isInOpenList(ty) {
		return nil
	}
	list := media.setting.GetStrv(gsKeyOpenFolder)
	list = append(list, ty)
	media.setting.SetStrv(gsKeyOpenFolder, list)

	media.delFromIgnoreList(ty)
	media.delFromRunList(ty)
	return nil
}

func (media *Media) addToRunList(ty string) error {
	if media.isInRunList(ty) {
		return nil
	}
	list := media.setting.GetStrv(gsKeyStartSoft)
	list = append(list, ty)
	media.setting.SetStrv(gsKeyStartSoft, list)

	media.delFromIgnoreList(ty)
	media.delFromOpenList(ty)
	return nil
}

func (media *Media) delFromIgnoreList(ty string) {
	list := media.setting.GetStrv(gsKeyIgnore)
	newList, deleted := delStrFromList(ty, list)
	if !deleted {
		return
	}
	media.setting.SetStrv(gsKeyIgnore, newList)
}

func (media *Media) delFromOpenList(ty string) {
	list := media.setting.GetStrv(gsKeyOpenFolder)
	newList, deleted := delStrFromList(ty, list)
	if !deleted {
		return
	}
	media.setting.SetStrv(gsKeyOpenFolder, newList)
}

func (media *Media) delFromRunList(ty string) {
	list := media.setting.GetStrv(gsKeyStartSoft)
	newList, deleted := delStrFromList(ty, list)
	if !deleted {
		return
	}
	media.setting.SetStrv(gsKeyStartSoft, newList)
}

func delStrFromList(s string, list []string) ([]string, bool) {
	var (
		ret     []string
		deleted bool
	)
	for _, v := range list {
		if s == v {
			deleted = true
			continue
		}
		ret = append(ret, v)
	}
	return ret, deleted
}
