/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

// #cgo pkg-config: x11 xi
// #include <stdlib.h>
// #include "button_map.h"
import "C"

import (
	"fmt"
	"unsafe"
)

func GetButtonMap(xid uint32, devName string) ([]byte, error) {
	if len(devName) == 0 {
		return nil, fmt.Errorf("Device name is empty")
	}

	var cbtnNum C.int = 0
	cname := C.CString(devName)
	defer C.free(unsafe.Pointer(cname))
	cbtnMap := C.get_button_map(C.ulong(xid), cname, &cbtnNum)
	if cbtnMap == nil {
		return nil, fmt.Errorf("Can not get button mapping for %s", devName)
	}
	defer C.free(unsafe.Pointer(cbtnMap))

	return ucharArrayToByte(cbtnMap, int(cbtnNum)), nil
}

func SetButtonMap(xid uint32, devName string, btnMap []byte) error {
	if len(devName) == 0 || len(btnMap) == 0 {
		return fmt.Errorf("Device name or map value is empty")
	}

	cbtnMap := byteArrayToUChar(btnMap)
	cname := C.CString(devName)
	defer C.free(unsafe.Pointer(cname))
	ret := C.set_button_map(C.ulong(xid), cname, &(cbtnMap[0]), C.int(len(btnMap)))
	if ret == -1 {
		return fmt.Errorf("Set button mapping failed")
	}

	return nil
}

func SetLeftHanded(xid uint32, devName string, useLeft bool) error {
	btnMap, err := GetButtonMap(xid, devName)
	if err != nil {
		return err
	}

	if len(btnMap) < 3 {
		return fmt.Errorf("Invalid device: button mapping number < 3")
	}

	if useLeft {
		if btnMap[0] == 3 && btnMap[2] == 1 {
			return nil
		}
		btnMap[0], btnMap[2] = 3, 1
	} else {
		if btnMap[0] == 1 && btnMap[2] == 3 {
			return nil
		}
		btnMap[0], btnMap[2] = 1, 3
	}

	return SetButtonMap(xid, devName, btnMap)
}

func CanLeftHanded(xid uint32, devName string) bool {
	btnMap, err := GetButtonMap(xid, devName)
	if err != nil {
		return false
	}

	if len(btnMap) < 3 {
		return false
	}

	if btnMap[0] == 3 && btnMap[2] == 1 {
		return true
	}

	return false
}
