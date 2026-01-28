// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package utils

// #cgo pkg-config: x11 xi
// #cgo CFLAGS: -W -Wall -fPIC -fstack-protector-all
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

	return C.GoBytes(unsafe.Pointer(cbtnMap), cbtnNum), nil
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
