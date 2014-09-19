/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
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

import (
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

type coordinateInfo struct {
	areas        []coordinateRange
	moveIntoFlag bool
	motionFlag   bool
	buttonFlag   bool
	keyFlag      bool
}

type coordinateRange struct {
	X1 int32
	Y1 int32
	X2 int32
	Y2 int32
}

type Manager struct {
	CursorInto    func(int32, int32, string)
	CursorOut     func(int32, int32, string)
	CursorMove    func(int32, int32, string)
	ButtonPress   func(int32, int32, int32, string)
	ButtonRelease func(int32, int32, int32, string)
	KeyPress      func(string, int32, int32, string)
	KeyRelease    func(string, int32, int32, string)

	CancelArea    func(string)
	CancelAllArea func() //resolution changed

	fullScreenId string
	md5RangeMap  map[string]*coordinateInfo
}

var _manager *Manager

func GetManager() *Manager {
	if _manager == nil {
		_manager = &Manager{}
		_manager.fullScreenId = ""
		_manager.md5RangeMap = make(map[string]*coordinateInfo)
	}
	return _manager
}

func (m *Manager) handleCursorEvent(x, y int32, press bool) {
	press = !press
	if m.CursorMove == nil {
		return
	}

	//fmt.Println("X:", x, "Y:", y, "Press:", press)
	inList, outList := m.getMd5List(x, y)
	for _, md5Str := range inList {
		if array, ok := m.md5RangeMap[md5Str]; ok {
			/* moveIntoFlag == true : mouse move in area */
			if !array.moveIntoFlag {
				if press {
					m.CursorInto(x, y, md5Str)
					array.moveIntoFlag = true
				}
			}

			if array.motionFlag {
				m.CursorMove(x, y, md5Str)
			}
		}
	}
	for _, md5Str := range outList {
		if array, ok := m.md5RangeMap[md5Str]; ok {
			/* moveIntoFlag == false : mouse move out area */
			if array.moveIntoFlag {
				m.CursorOut(x, y, md5Str)
				array.moveIntoFlag = false
			}
		}
	}

	if len(m.fullScreenId) > 0 {
		m.CursorMove(x, y, m.fullScreenId)
	}
}

func (m *Manager) handleButtonEvent(button int32, press bool, x, y int32) {
	if m.ButtonPress == nil {
		return
	}

	list, _ := m.getMd5List(x, y)
	for _, md5Str := range list {
		if array, ok := m.md5RangeMap[md5Str]; ok {
			if !array.buttonFlag {
				continue
			}
			if press {
				m.ButtonPress(button, x, y, md5Str)
			} else {
				m.ButtonRelease(button, x, y, md5Str)
			}
		}
	}

	if len(m.fullScreenId) > 0 {
		if press {
			m.ButtonPress(button, x, y, m.fullScreenId)
		} else {
			m.ButtonRelease(button, x, y, m.fullScreenId)
		}
	}
}

func (m *Manager) handleKeyboardEvent(code int32, press bool, x, y int32) {
	if m.KeyPress == nil {
		return
	}

	list, _ := m.getMd5List(x, y)
	for _, md5Str := range list {
		if array, ok := m.md5RangeMap[md5Str]; ok {
			if !array.keyFlag {
				continue
			}
			if press {
				m.KeyPress(keyCode2Str(code), x, y, md5Str)
			} else {
				m.KeyRelease(keyCode2Str(code), x, y, md5Str)
			}
		}
	}

	if len(m.fullScreenId) > 0 {
		if press {
			m.KeyPress(keyCode2Str(code), x, y, m.fullScreenId)
		} else {
			m.KeyRelease(keyCode2Str(code), x, y, m.fullScreenId)
		}
	}
}

func (m *Manager) cancelAllReigsterArea() {
	m.md5RangeMap = make(map[string]*coordinateInfo)
	m.fullScreenId = ""

	m.CancelAllArea()
}

func (m *Manager) RegisterArea(dbusMsg dbus.DMessage, x1, y1, x2, y2, flag int32) (string, error) {
	return m.RegisterAreas(dbusMsg,
		[]coordinateRange{coordinateRange{x1, y1, x2, y2}}, flag)
}

func (m *Manager) RegisterAreas(dbusMsg dbus.DMessage, areas []coordinateRange, flag int32) (md5Str string, err error) {
	var ok bool
	if md5Str, ok = m.sumAreasMd5(areas, flag, dbusMsg.GetSenderPID()); !ok {
		err = fmt.Errorf("sumAreasMd5 failed:", areas)
		return
	}

	logger.Debug("md5Str:", md5Str)
	if _, ok := m.md5RangeMap[md5Str]; ok {
		return md5Str, nil
	}

	info := &coordinateInfo{}
	info.areas = areas
	info.moveIntoFlag = false
	info.buttonFlag = hasButtonFlag(flag)
	info.keyFlag = hasKeyFlag(flag)
	info.motionFlag = hasMotionFlag(flag)
	m.md5RangeMap[md5Str] = info

	return md5Str, nil
}

func (m *Manager) RegisterFullScreen() (md5Str string) {
	if len(m.fullScreenId) < 1 {
		m.fullScreenId, _ = dutils.SumStrMd5("")
	}
	logger.Debug("fullScreenId: ", m.fullScreenId)

	return m.fullScreenId
}

func (m *Manager) UnregisterArea(md5Str string) {
	if _, ok := m.md5RangeMap[md5Str]; ok {
		delete(m.md5RangeMap, md5Str)
		m.CancelArea(md5Str)
	}

	if md5Str == m.fullScreenId {
		m.fullScreenId = ""
		m.CancelArea(md5Str)
	}
}

func (m *Manager) getMd5List(x, y int32) ([]string, []string) {
	inList := []string{}
	outList := []string{}

	for md5Str, array := range m.md5RangeMap {
		inFlag := false
		for _, area := range array.areas {
			if isInArea(x, y, area) {
				inFlag = true
				if !isInMd5List(md5Str, inList) {
					inList = append(inList, md5Str)
				}
			}
		}
		if !inFlag {
			if !isInMd5List(md5Str, outList) {
				outList = append(outList, md5Str)
			}
		}
	}

	return inList, outList
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		MouseAreaDest,
		"/com/deepin/api/XMouseArea",
		"com.deepin.api.XMouseArea",
	}
}

func (m *Manager) sumAreasMd5(areas []coordinateRange, flag int32,
	pid uint32) (md5Str string, ok bool) {
	if len(areas) < 1 {
		return
	}

	content := ""
	for _, area := range areas {
		if len(content) > 1 {
			content += "-"
		}
		content += fmt.Sprintf("%v-%v-%v-%v", area.X1, area.Y1, area.X2, area.Y2)
	}
	content += fmt.Sprintf("-%v-%v", flag, pid)

	logger.Debug("areas content:", content)
	md5Str, ok = dutils.SumStrMd5(content)

	return
}
