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
// #include "property.h"
// #include "list.h"
// #include "type.h"
// #include "keyboard.h"
import "C"

import (
	"fmt"
	"os"
	"unsafe"

	. "pkg.deepin.io/dde/api/dxinput/common"
	"pkg.deepin.io/dde/api/dxinput/kwayland"
)

const (
	// see 'property.c' MAX_BUF_LEN
	maxBufferLen = 1000
)

func ListDevice() DeviceInfos {
	if len(os.Getenv("WAYLAND_DISPLAY")) != 0 {
		infos, _ := kwayland.ListDevice()
		if len(infos) != 0 {
			return infos
		}
	}
	var cNum C.int = 0
	cInfos := C.list_device(&cNum)
	if cNum == 0 && cInfos == nil {
		return nil
	}
	defer C.free_device_list(cInfos, cNum)

	var infos DeviceInfos
	itemLen := unsafe.Sizeof(*cInfos)
	for i := C.int(0); i < cNum; i++ {
		cInfo := (*C.DeviceInfo)(unsafe.Pointer(uintptr(unsafe.Pointer(cInfos)) + uintptr(i)*itemLen))
		infos = append(infos, &DeviceInfo{
			Id:      int32(cInfo.id),
			Type:    int32(cInfo.ty),
			Name:    C.GoString(cInfo.name),
			Enabled: (int(cInfo.enabled) == 1),
		})
	}
	return infos
}

func QueryDeviceType(id int32) int32 {
	return int32(C.query_device_type(C.int(id)))
}

func SetKeyboardRepeat(enabled bool, delay, interval uint32) error {
	var repeated int = 0
	if enabled {
		repeated = 1
	}

	ret := C.set_keyboard_repeat(C.int(repeated),
		C.uint(delay), C.uint(interval))
	if ret != 0 {
		return fmt.Errorf("Not found compatible version of the Xkb extension in the server")
	}

	return nil
}

func IsPropertyExist(id int32, prop string) bool {
	if len(prop) == 0 {
		return false
	}

	cprop := C.CString(prop)
	defer C.free(unsafe.Pointer(cprop))
	ret := C.is_property_exist(C.int(id), cprop)
	return int(ret) != 0
}

func GetProperty(id int32, prop string) ([]byte, int32) {
	if len(prop) == 0 {
		return nil, 0
	}

	cprop := C.CString(prop)
	defer C.free(unsafe.Pointer(cprop))
	nitems := C.int(0)
	cdatas := C.get_prop(C.int(id), cprop, &nitems)
	if cdatas == nil {
		return nil, 0
	}

	datas := ucharArrayToByte(cdatas, maxBufferLen)
	return datas, int32(nitems)
}

func SetInt8Prop(id int32, prop string, values []int8) error {
	cdatas := byteArrayToUChar(WriteInt8(values))
	cprop := C.CString(prop)
	defer C.free(unsafe.Pointer(cprop))
	ret := C.set_prop_int(C.int(id), cprop, &(cdatas[0]),
		C.int(len(values)), 8)
	if int(ret) == -1 {
		return fmt.Errorf("[SetPropInt8] failed for: '%v -- %s -- %v'",
			id, prop, values)
	}
	return nil
}

func SetInt16Prop(id int32, prop string, values []int16) error {
	cdatas := byteArrayToUChar(WriteInt16(values))
	cprop := C.CString(prop)
	defer C.free(unsafe.Pointer(cprop))
	ret := C.set_prop_int(C.int(id), cprop, &(cdatas[0]),
		C.int(len(values)), 16)
	if int(ret) == -1 {
		return fmt.Errorf("[SetPropInt16] failed for: '%v -- %s -- %v'",
			id, prop, values)
	}
	return nil
}

func SetInt32Prop(id int32, prop string, values []int32) error {
	cdatas := byteArrayToUChar(WriteInt32(values))
	cprop := C.CString(prop)
	defer C.free(unsafe.Pointer(cprop))
	ret := C.set_prop_int(C.int(id), cprop, &(cdatas[0]),
		C.int(len(values)), 32)
	if int(ret) == -1 {
		return fmt.Errorf("[SetPropInt32] failed for: '%v -- %s -- %v'",
			id, prop, values)
	}
	return nil
}

func SetFloat32Prop(id int32, prop string, values []float32) error {
	cdatas := byteArrayToUChar(WriteFloat32(values))
	cprop := C.CString(prop)
	defer C.free(unsafe.Pointer(cprop))
	ret := C.set_prop_float(C.int(id), cprop, &(cdatas[0]),
		C.int(len(values)))
	if int(ret) == -1 {
		return fmt.Errorf("[SetPropFloat] failed for: '%v -- %s -- %v'",
			id, prop, values)
	}
	return nil
}

func ucharArrayToByte(cData *C.uchar, length int) []byte {
	if cData == nil {
		return nil
	}
	cItemSize := unsafe.Sizeof(*cData)

	var data []byte
	for i := 0; i < length; i++ {
		cdata := (*C.uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(cData)) + uintptr(i)*cItemSize))
		if cdata == nil {
			break
		}
		data = append(data, byte(*cdata))
	}
	return data
}

func byteArrayToUChar(datas []byte) []C.uchar {
	var cdatas []C.uchar
	for i := 0; i < len(datas); i++ {
		cdatas = append(cdatas, C.uchar(datas[i]))
	}
	return cdatas
}
