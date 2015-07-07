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
	"pkg.deepin.io/lib/dbus"
	dutils "pkg.deepin.io/lib/utils"
	"sync"
)

const _FullscreenId = "d41d8cd98f00b204e9800998ecf8427e"

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

	idAreaInfoMap   map[string]*coordinateInfo
	idReferCountMap map[string]int32

	countLock sync.Mutex
}

var _manager *Manager

func GetManager() *Manager {
	if _manager == nil {
		_manager = &Manager{}
		_manager.idAreaInfoMap = make(map[string]*coordinateInfo)
		_manager.idReferCountMap = make(map[string]int32)
	}
	return _manager
}

func (m *Manager) handleCursorEvent(x, y int32, press bool) {
	press = !press

	inList, outList := m.getIdList(x, y)
	for _, id := range inList {
		array, ok := m.idAreaInfoMap[id]
		if !ok {
			continue
		}

		/* moveIntoFlag == true : mouse move in area */
		if !array.moveIntoFlag {
			if press {
				dbus.Emit(m, "CursorInto", x, y, id)
				array.moveIntoFlag = true
			}
		}

		if array.motionFlag {
			dbus.Emit(m, "CursorMove", x, y, id)
		}
	}
	for _, id := range outList {
		array, ok := m.idAreaInfoMap[id]
		if !ok {
			continue
		}

		/* moveIntoFlag == false : mouse move out area */
		if array.moveIntoFlag {
			dbus.Emit(m, "CursorOut", x, y, id)
			array.moveIntoFlag = false
		}
	}

	_, ok := m.idReferCountMap[_FullscreenId]
	if !ok {
		return
	}

	dbus.Emit(m, "CursorMove", x, y, _FullscreenId)
}

func (m *Manager) handleButtonEvent(button int32, press bool, x, y int32) {

	list, _ := m.getIdList(x, y)
	for _, id := range list {
		array, ok := m.idAreaInfoMap[id]
		if !ok || !array.buttonFlag {
			continue
		}

		if press {
			dbus.Emit(m, "ButtonPress", x, y, id)
		} else {
			dbus.Emit(m, "ButtonRelease", x, y, id)
		}
	}

	_, ok := m.idReferCountMap[_FullscreenId]
	if !ok {
		return
	}

	if press {
		dbus.Emit(m, "ButtonPress", button, x, y, _FullscreenId)
	} else {
		dbus.Emit(m, "ButtonRelease", button, x, y, _FullscreenId)
	}
}

func (m *Manager) handleKeyboardEvent(code int32, press bool, x, y int32) {
	list, _ := m.getIdList(x, y)
	for _, id := range list {
		array, ok := m.idAreaInfoMap[id]
		if !ok || !array.keyFlag {
			continue
		}

		if press {
			dbus.Emit(m, "KeyPress", keyCode2Str(code), x, y, id)
		} else {
			dbus.Emit(m, "KeyRelease", keyCode2Str(code), x, y, id)
		}
	}

	_, ok := m.idReferCountMap[_FullscreenId]
	if !ok {
		return
	}

	if press {
		dbus.Emit(m, "KeyPress", keyCode2Str(code), x, y, _FullscreenId)
	} else {
		dbus.Emit(m, "KeyRelease", keyCode2Str(code), x, y, _FullscreenId)
	}
}

func (m *Manager) cancelAllReigsterArea() {
	m.idAreaInfoMap = make(map[string]*coordinateInfo)
	m.idReferCountMap = make(map[string]int32)

	dbus.Emit(m, "CancelAllArea")
}

func (m *Manager) RegisterArea(x1, y1, x2, y2, flag int32) (string, error) {
	return m.RegisterAreas(
		[]coordinateRange{coordinateRange{x1, y1, x2, y2}},
		flag)
}

func (m *Manager) RegisterAreas(areas []coordinateRange, flag int32) (id string, err error) {
	md5Str, ok := m.sumAreasMd5(areas, flag)
	if !ok {
		err = fmt.Errorf("sumAreasMd5 failed:", areas)
		return
	}
	id = md5Str

	m.countLock.Lock()
	defer m.countLock.Unlock()
	_, ok = m.idReferCountMap[id]
	if ok {
		m.idReferCountMap[id] += 1
		return id, nil
	}

	info := &coordinateInfo{}
	info.areas = areas
	info.moveIntoFlag = false
	info.buttonFlag = hasButtonFlag(flag)
	info.keyFlag = hasKeyFlag(flag)
	info.motionFlag = hasMotionFlag(flag)

	m.idAreaInfoMap[id] = info
	m.idReferCountMap[id] = 1

	return id, nil
}

func (m *Manager) RegisterFullScreen() (id string) {
	m.countLock.Lock()
	defer m.countLock.Unlock()
	_, ok := m.idReferCountMap[_FullscreenId]
	if !ok {
		m.idReferCountMap[_FullscreenId] = 1
	} else {
		m.idReferCountMap[_FullscreenId] += 1
	}

	return _FullscreenId
}

func (m *Manager) UnregisterArea(dbusMsg dbus.DMessage, id string) {
	_, ok := m.idReferCountMap[id]
	if !ok {
		return
	}

	m.countLock.Lock()
	defer m.countLock.Unlock()
	m.idReferCountMap[id] -= 1
	if m.idReferCountMap[id] == 0 {
		delete(m.idReferCountMap, id)
		delete(m.idAreaInfoMap, id)
	}
	dbus.Emit(m, "CancelArea", fmt.Sprintf("%v", dbusMsg.GetSenderPID()))
}

func (m *Manager) getIdList(x, y int32) ([]string, []string) {
	inList := []string{}
	outList := []string{}

	for id, array := range m.idAreaInfoMap {
		inFlag := false
		for _, area := range array.areas {
			if isInArea(x, y, area) {
				inFlag = true
				if !isInIdList(id, inList) {
					inList = append(inList, id)
				}
			}
		}
		if !inFlag {
			if !isInIdList(id, outList) {
				outList = append(outList, id)
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

func (m *Manager) sumAreasMd5(areas []coordinateRange, flag int32) (md5Str string, ok bool) {
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
	content += fmt.Sprintf("-%v", flag)

	logger.Debug("areas content:", content)
	md5Str, ok = dutils.SumStrMd5(content)

	return
}
