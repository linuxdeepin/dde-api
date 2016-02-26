/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
	"unsafe"
)

const (
	DevTypeUnknown int32 = iota
	DevTypeKeyboard
	DevTypeMouse
	DevTypeTouchpad
	DevTypeWacom
	DevTypeTouchscreen
)

const (
	// see 'property.c' MAX_BUF_LEN
	maxBufferLen = 1000
)

type DeviceInfo struct {
	Id      int32
	Type    int32
	Name    string
	Enabled bool
}
type DeviceInfos []*DeviceInfo

func ListDevice() DeviceInfos {
	var cnum C.int = 0
	cinfos := C.list_device(&cnum)
	if cnum == 0 && cinfos == nil {
		return nil
	}
	defer C.free_device_list(cinfos, cnum)

	var infos DeviceInfos
	clist := uintptr(unsafe.Pointer(cinfos))
	itemLen := unsafe.Sizeof(*cinfos)
	for i := C.int(0); i < cnum; i++ {
		cinfo := (*C.DeviceInfo)(unsafe.Pointer(
			clist + uintptr(i)*itemLen))
		infos = append(infos, &DeviceInfo{
			Id:      int32(cinfo.id),
			Type:    int32(cinfo.ty),
			Name:    C.GoString(cinfo.name),
			Enabled: (int(cinfo.enabled) == 1),
		})
	}
	return infos
}

func (infos DeviceInfos) Get(id int32) *DeviceInfo {
	for _, info := range infos {
		if info.Id == id {
			return info
		}
	}
	return nil
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
	if int(ret) == 0 {
		return false
	}
	return true
}

func GetProperty(id int32, prop string, nitems int32) []byte {
	if len(prop) == 0 {
		return nil
	}

	cprop := C.CString(prop)
	defer C.free(unsafe.Pointer(cprop))
	cdatas := C.get_prop(C.int(id), cprop, C.int(nitems))
	if cdatas == nil {
		return nil
	}

	datas := ucharArrayToByte(cdatas, maxBufferLen)
	return datas
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

func ucharArrayToByte(cdatas *C.uchar, length int) []byte {
	clist := uintptr(unsafe.Pointer(cdatas))
	citemLen := unsafe.Sizeof(*cdatas)

	var datas []byte
	for i := 0; i < length; i++ {
		cdata := (*C.uchar)(unsafe.Pointer(clist + uintptr(i)*citemLen))
		datas = append(datas, byte(*cdata))
	}
	return datas
}

func byteArrayToUChar(datas []byte) []C.uchar {
	var cdatas []C.uchar
	for i := 0; i < len(datas); i++ {
		cdatas = append(cdatas, C.uchar(datas[i]))
	}
	return cdatas
}
