// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"encoding/json"
	"path/filepath"

	"os"

	"io/ioutil"

	"github.com/linuxdeepin/dde-api/huangli"
	"github.com/linuxdeepin/go-lib/calendar"
	"github.com/linuxdeepin/go-lib/utils"
	"github.com/linuxdeepin/go-lib/xdg/basedir"
)

// HuangLiInfo huang li
type HuangLiInfo struct {
	calendar.LunarDayInfo
	Avoid string
	Suit  string
}

// HuangLiInfoList huang li list
type HuangLiInfoList []*HuangLiInfo

// HuangLiMonthInfo huang li month info
type HuangLiMonthInfo struct {
	FirstDayWeek int32
	Days         int32
	Datas        HuangLiInfoList
}

const (
	defaultHuangLiDBFile  = "/usr/share/dde-api/data/huangli.db"
	defaultHuangLiVerFile = "/usr/share/dde-api/data/huangli.version"
)

var (
	_hasHuangLi bool
)

func initHuangLi() {
	err := huangli.Init(getDBFile())
	if err != nil {
		logger.Error("Failed to open huangli db:", err)
		_hasHuangLi = false
		return
	}
	_hasHuangLi = true
}

func finalizeHuangLi() {
	huangli.Finalize()
}

// String json marshal
func (info *HuangLiInfo) String() string {
	data, _ := json.Marshal(info)
	return string(data)
}

// String json marshal
func (info *HuangLiMonthInfo) String() string {
	data, _ := json.Marshal(info)
	return string(data)
}

func newHuangLiInfoList(lunarDays []calendar.LunarDayInfo, days DayInfoList) (list HuangLiInfoList) {
	var infos huangli.HuangLiList
	if _hasHuangLi {
		infos, _ = huangli.NewHuangLiList(days.GetIDList())
	} else {
		for i := 0; i < len(lunarDays); i++ {
			infos = append(infos, &huangli.HuangLi{})
		}
	}
	for i := 0; i < len(lunarDays); i++ {
		list = append(list, &HuangLiInfo{
			LunarDayInfo: lunarDays[i],
			Avoid:        infos[i].Avoid,
			Suit:         infos[i].Suit,
		})
	}
	return
}

func newFestivalList(year, month int) (huangli.FestivalList, error) {
	return huangli.NewFestivalList(year, month)
}

func getDBFile() string {
	filename := filepath.Join(basedir.GetUserConfigDir(), "deepin", "dde-api", "huangli.db")
	if utils.IsFileExist(filename) && checkDBVersion() {
		return filename
	}

	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		logger.Warning("Failed to mkdir for huangli db:", err)
	} else {
		err := utils.CopyFile(defaultHuangLiDBFile, filename)
		if err != nil {
			logger.Warning("Failed to copy huangli db file:", err)
		}
		versionFile := filepath.Join(basedir.GetUserConfigDir(), "deepin", "dde-api", "huangli.version")
		err = utils.CopyFile(defaultHuangLiVerFile, versionFile)
		if err != nil {
			logger.Warning("Failed to copy huangli version file:", err)
		}
	}
	return filename

}

func checkDBVersion() bool {
	filename := filepath.Join(basedir.GetUserConfigDir(), "deepin", "dde-api", "huangli.version")
	if !utils.IsFileExist(filename) {
		return false
	}
	src, err := ioutil.ReadFile(defaultHuangLiVerFile)
	if err != nil {
		return false
	}
	dest, err := ioutil.ReadFile(filename)
	if err != nil {
		return false
	}
	return string(src) == string(dest)
}
