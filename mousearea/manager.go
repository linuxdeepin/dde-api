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
	fullScreenId  int32
	MotionInto    func(int32, int32, int32)
	MotionOut     func(int32, int32, int32)
	MotionMove    func(int32, int32, int32)
	ButtonPress   func(int32, int32, int32, int32)
	ButtonRelease func(int32, int32, int32, int32)
	KeyPress      func(string, int32, int32, int32)
	KeyRelease    func(string, int32, int32, int32)
	CancelAllArea func() //resolution changed
	idRangeMap    map[int32]*coordinateInfo
}

var _manager *Manager

func GetManager() *Manager {
	if _manager == nil {
		_manager = &Manager{}
		_manager.fullScreenId = -1
		_manager.idRangeMap = make(map[int32]*coordinateInfo)
	}
	return _manager
}

func (obj *Manager) handleMotionEvent(x, y int32, press bool) {
	press = !press
	if obj.MotionMove == nil {
		return
	}

	//fmt.Println("X:", x, "Y:", y, "Press:", press)
	inList, outList := obj.getIDList(x, y)
	for _, cookie := range inList {
		if array, ok := obj.idRangeMap[cookie]; ok {
			/* moveIntoFlag == true : mouse move in area */
			if !array.moveIntoFlag {
				if press {
					obj.MotionInto(x, y, cookie)
					array.moveIntoFlag = true
				}
			}

			if array.motionFlag {
				obj.MotionMove(x, y, cookie)
			}
		}
	}
	for _, cookie := range outList {
		if array, ok := obj.idRangeMap[cookie]; ok {
			/* moveIntoFlag == false : mouse move out area */
			if array.moveIntoFlag {
				obj.MotionOut(x, y, cookie)
				array.moveIntoFlag = false
			}
		}
	}

	if obj.fullScreenId != -1 {
		obj.MotionMove(x, y, obj.fullScreenId)
	}
}

func (obj *Manager) handleButtonEvent(button int32, press bool, x, y int32) {
	if obj.ButtonPress == nil {
		return
	}

	cookies, _ := obj.getIDList(x, y)
	for _, cookie := range cookies {
		if array, ok := obj.idRangeMap[cookie]; ok {
			if !array.buttonFlag {
				continue
			}
			if press {
				obj.ButtonPress(button, x, y, cookie)
			} else {
				obj.ButtonRelease(button, x, y, cookie)
			}
		}
	}

	if obj.fullScreenId != -1 {
		if press {
			obj.ButtonPress(button, x, y, obj.fullScreenId)
		} else {
			obj.ButtonRelease(button, x, y, obj.fullScreenId)
		}
	}
}

func (obj *Manager) handleKeyboardEvent(code int32, press bool, x, y int32) {
	if obj.KeyPress == nil {
		return
	}
	cookies, _ := obj.getIDList(x, y)
	for _, cookie := range cookies {
		if array, ok := obj.idRangeMap[cookie]; ok {
			if !array.keyFlag {
				continue
			}
			if press {
				obj.KeyPress(keyCode2Str(code), x, y, cookie)
			} else {
				obj.KeyRelease(keyCode2Str(code), x, y, cookie)
			}
		}
	}

	if obj.fullScreenId != -1 {
		if press {
			obj.KeyPress(keyCode2Str(code), x, y, obj.fullScreenId)
		} else {
			obj.KeyRelease(keyCode2Str(code), x, y, obj.fullScreenId)
		}
	}
}

func (obj *Manager) cancelAllReigsterArea() {
	obj.idRangeMap = make(map[int32]*coordinateInfo)
	obj.fullScreenId = -1

	obj.CancelAllArea()
}

func (obj *Manager) RegisterArea(x1, y1, x2, y2, flag int32) int32 {
	return obj.RegisterAreas([]coordinateRange{coordinateRange{x1, y1, x2, y2}}, flag)
}

func (obj *Manager) RegisterAreas(areas []coordinateRange, flag int32) int32 {
	cookie := genID()
	Logger.Debug("ID: ", cookie)

	info := &coordinateInfo{}
	info.areas = areas
	info.moveIntoFlag = false
	info.buttonFlag = hasButtonFlag(flag)
	info.keyFlag = hasKeyFlag(flag)
	info.motionFlag = hasMotionFlag(flag)
	obj.idRangeMap[cookie] = info

	return cookie
}

func (obj *Manager) RegisterFullScreen() int32 {
	if obj.fullScreenId == -1 {
		cookie := genID()
		obj.fullScreenId = cookie
	}
	Logger.Debug("ID: ", obj.fullScreenId)

	return obj.fullScreenId
}

func (obj *Manager) UnregisterArea(cookie int32) {
	if _, ok := obj.idRangeMap[cookie]; ok {
		delete(obj.idRangeMap, cookie)
	}
	if cookie == obj.fullScreenId {
		obj.fullScreenId = -1
	}
}

func (obj *Manager) getIDList(x, y int32) ([]int32, []int32) {
	inList := []int32{}
	outList := []int32{}

	for id, array := range obj.idRangeMap {
		inFlag := false
		for _, area := range array.areas {
			if isInArea(x, y, area) {
				inFlag = true
				if !isInIDList(id, inList) {
					inList = append(inList, id)
				}
			}
		}
		if !inFlag {
			if !isInIDList(id, outList) {
				outList = append(outList, id)
			}
		}
	}

	return inList, outList
}

func (obj *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		MouseAreaDest,
		"/com/deepin/api/XMouseArea",
		"com.deepin.api.XMouseArea",
	}
}
