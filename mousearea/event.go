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

// #cgo CFLAGS: -g -Wall
// #cgo pkg-config: x11 xtst glib-2.0
// #include "record.h"
import "C"

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/keybind"
	"strings"
)

type coordinateRange struct {
	X1 int32
	Y1 int32
	X2 int32
	Y2 int32
}

//export parseMotionEvent
func parseMotionEvent(_x, _y, _press int32) {
	coorX := int32(_x)
	coorY := int32(_y)
	pressFlag := int32(_press)

	inList, outList := getIDList(coorX, coorY)
	for _, cookie := range inList {
		if array, ok := idRangeMap[cookie]; ok {
			/* moveIntoFlag == true : mouse move in area */
			if !array.moveIntoFlag {
				array.moveIntoFlag = true
				if pressFlag != 1 {
					opMouse.MotionInto(coorX, coorY, cookie)
				}
			}

			if array.motionFlag {
				opMouse.MotionMove(coorX, coorY, cookie)
			}
		}
	}
	for _, cookie := range outList {
		if array, ok := idRangeMap[cookie]; ok {
			/* moveIntoFlag == false : mouse move out area */
			if array.moveIntoFlag {
				array.moveIntoFlag = false
				opMouse.MotionOut(coorX, coorY, cookie)
			}
		}
	}

	if opMouse.FullScreenId != -1 {
		opMouse.MotionMove(coorX, coorY, opMouse.FullScreenId)
	}
}

//export parseButtonEvent
func parseButtonEvent(_code, _type, _x, _y int32) {
	btnCode := int32(_code)
	coorX := int32(_x)
	coorY := int32(_y)
	tmp := int32(_type)
	coorType := false
	if tmp == C.BUTTON_PRESS {
		coorType = true
	} else {
		coorType = false
	}

	btnStr := ""
	if btnCode == 1 {
		btnStr = "LeftButton"
	} else if btnCode == 2 {
		btnStr = "Middlebutton"
	} else if btnCode == 3 {
		btnStr = "Rightbutton"
	} else if btnCode == 4 {
		btnStr = "RollForward"
	} else if btnCode == 5 {
		btnStr = "RollBack"
	}

	cookies, _ := getIDList(coorX, coorY)
	//logger.Info("Button Cookies: ", cookies)
	//logger.Infof("\tX: %d, Y: %d\n\n", coorX, coorY)
	for _, cookie := range cookies {
		if array, ok := idRangeMap[cookie]; ok {
			if !array.buttonFlag {
				continue
			}
			if coorType {
				opMouse.ButtonPress(btnStr, coorX, coorY, cookie)
			} else {
				opMouse.ButtonRelease(btnStr, coorX, coorY, cookie)
			}
		}
	}

	if opMouse.FullScreenId != -1 {
		if coorType {
			opMouse.ButtonPress(btnStr, coorX, coorY, opMouse.FullScreenId)
		} else {
			opMouse.ButtonRelease(btnStr, coorX, coorY, opMouse.FullScreenId)
		}
	}
}

//export parseKeyboardEvent
func parseKeyboardEvent(_code, _type, _x, _y int32) {
	keyCode := int32(_code)
	coorX := int32(_x)
	coorY := int32(_y)
	tmp := int32(_type)
	coorType := false
	if tmp == C.KEY_PRESS {
		coorType = true
	} else {
		coorType = false
	}

	keyStr := keybind.LookupString(X, 0, xproto.Keycode(keyCode))
	if keyStr == " " {
		keyStr = "space"
	}
	keyStr = strings.ToLower(keyStr)
	//logger.Info("KeyStr: ", keyStr)

	cookies, _ := getIDList(coorX, coorY)
	//logger.Info("Keyboard Cookies: ", cookies)
	//logger.Infof("\tX: %d, Y: %d\n\n", coorX, coorY)
	for _, cookie := range cookies {
		if array, ok := idRangeMap[cookie]; ok {
			if !array.keyFlag {
				continue
			}
			if coorType {
				opMouse.KeyPress(keyStr, coorX, coorY, cookie)
			} else {
				opMouse.KeyRelease(keyStr, coorX, coorY, cookie)
			}
		}
	}

	if opMouse.FullScreenId != -1 {
		if coorType {
			opMouse.KeyPress(keyStr, coorX, coorY, opMouse.FullScreenId)
		} else {
			opMouse.KeyRelease(keyStr, coorX, coorY, opMouse.FullScreenId)
		}
	}
}

func cancelAllReigsterArea() {
	list := []int32{}

	for id, _ := range idRangeMap {
		list = append(list, id)
		delete(idRangeMap, id)
	}

	opMouse.FullScreenId = -1

	println("map len:", len(idRangeMap))
	for _, cookie := range list {
		opMouse.CancelAllArea(1365, 767, cookie)
	}
}

func getIDList(x, y int32) ([]int32, []int32) {
	inList := []int32{}
	outList := []int32{}

	for id, array := range idRangeMap {
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

func isInArea(x, y int32, area coordinateRange) bool {
	if (x >= area.X1 && x <= area.X2) &&
		(y >= area.Y1 && y <= area.Y2) {
		return true
	}

	return false
}

func isInIDList(id int32, list []int32) bool {
	for _, v := range list {
		if id == v {
			return true
		}
	}

	return false
}
