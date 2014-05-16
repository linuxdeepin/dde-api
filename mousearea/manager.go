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
	"dlib/dbus"
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
	FullScreenId  int32
	MotionInto    func(int32, int32, int32)
	MotionOut     func(int32, int32, int32)
	MotionMove    func(int32, int32, int32)
	ButtonPress   func(string, int32, int32, int32)
	ButtonRelease func(string, int32, int32, int32)
	KeyPress      func(string, int32, int32, int32)
	KeyRelease    func(string, int32, int32, int32)
	CancelAllArea func(int32, int32, int32) //resolution changed
}

var (
	idRangeMap = make(map[int32]*coordinateInfo)
)

var _manager *Manager

func GetManager() *Manager {
	if _manager == nil {
		_manager = &Manager{}
		_manager.FullScreenId = -1
	}
	return _manager
}

func (m *Manager) handleMotionEvent(x, y int32, press bool) {
	press = !press
	if m.MotionMove == nil {
		return
	}

	//fmt.Println("X:", x, "Y:", y, "Press:", press)
	inList, outList := getIDList(x, y)
	for _, cookie := range inList {
		if array, ok := idRangeMap[cookie]; ok {
			/* moveIntoFlag == true : mouse move in area */
			if !array.moveIntoFlag {
				array.moveIntoFlag = true
				if press {
					m.MotionInto(x, y, cookie)
				}
			}

			if array.motionFlag {
				m.MotionMove(x, y, cookie)
			}
		}
	}
	for _, cookie := range outList {
		if array, ok := idRangeMap[cookie]; ok {
			/* moveIntoFlag == false : mouse move out area */
			if array.moveIntoFlag {
				array.moveIntoFlag = false
				m.MotionOut(x, y, cookie)
			}
		}
	}

	if m.FullScreenId != -1 {
		m.MotionMove(x, y, m.FullScreenId)
	}
}

func (m *Manager) handleButtonEvent(button int32, press bool, x, y int32) {
	if m.ButtonPress == nil {
		return
	}
	btnStr := buttonCode2str(button)

	cookies, _ := getIDList(x, y)
	for _, cookie := range cookies {
		if array, ok := idRangeMap[cookie]; ok {
			if !array.buttonFlag {
				continue
			}
			if press {
				m.ButtonPress(btnStr, x, y, cookie)
			} else {
				m.ButtonRelease(btnStr, x, y, cookie)
			}
		}
	}

	if m.FullScreenId != -1 {
		if press {
			m.ButtonPress(btnStr, x, y, m.FullScreenId)
		} else {
			m.ButtonRelease(btnStr, x, y, m.FullScreenId)
		}
	}
}

func (m *Manager) handleKeyboardEvent(code int32, press bool, x, y int32) {
	if m.KeyPress == nil {
		return
	}
	cookies, _ := getIDList(x, y)
	for _, cookie := range cookies {
		if array, ok := idRangeMap[cookie]; ok {
			if !array.keyFlag {
				continue
			}
			if press {
				m.KeyPress(keyCode2Str(code), x, y, cookie)
			} else {
				m.KeyRelease(keyCode2Str(code), x, y, cookie)
			}
		}
	}

	if m.FullScreenId != -1 {
		if press {
			m.KeyPress(keyCode2Str(code), x, y, m.FullScreenId)
		} else {
			m.KeyRelease(keyCode2Str(code), x, y, m.FullScreenId)
		}
	}
}

func (m *Manager) cancelAllReigsterArea() {
	list := []int32{}

	for id, _ := range idRangeMap {
		list = append(list, id)
		delete(idRangeMap, id)
	}

	m.FullScreenId = -1

	println("map len:", len(idRangeMap))
	for _, cookie := range list {
		m.CancelAllArea(1365, 767, cookie)
	}
}

/*
 * flags:
 *      motionFlag: 001
 *      buttonFlag: 010
 *      keyFlag:    100
 *      allFlag:    111
 */
func (op *Manager) RegisterArea(x1, y1, x2, y2, flag int32) int32 {
	return op.RegisterAreas([]coordinateRange{coordinateRange{x1, y1, x2, y2}}, flag)
}

func (op *Manager) RegisterAreas(areas []coordinateRange, flag int32) int32 {
	cookie := genID()
	Logger.Debug("ID: ", cookie)

	info := &coordinateInfo{}
	info.areas = areas
	info.moveIntoFlag = false
	info.buttonFlag = false
	info.keyFlag = false
	info.motionFlag = false
	if flag >= 0 && flag <= 7 {
		if flag%2 == 1 {
			info.motionFlag = true
		}

		flag = flag >> 1
		if flag%2 == 1 {
			info.buttonFlag = true
		}

		flag = flag >> 1
		if flag%2 == 1 {
			info.keyFlag = true
		}
	}
	idRangeMap[cookie] = info

	return cookie
}

func (op *Manager) RegisterFullScreen() int32 {
	if op.FullScreenId == -1 {
		cookie := genID()
		op.FullScreenId = cookie
	}
	Logger.Debug("ID: ", op.FullScreenId)

	return op.FullScreenId
}

func (op *Manager) UnregisterArea(cookie int32) {
	if _, ok := idRangeMap[cookie]; ok {
		delete(idRangeMap, cookie)
	}
	if cookie == op.FullScreenId {
		op.FullScreenId = -1
	}
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		MouseAreaDest,
		"/com/deepin/api/XMouseArea",
		"com.deepin.api.XMouseArea",
	}
}
