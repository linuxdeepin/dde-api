/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package utils

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

func ReadInt8(datas []byte, nitems int32) []int8 {
	reader := bytes.NewReader(datas)
	var array []int8
	for i := int32(0); i < nitems; i++ {
		var tmp int8
		binary.Read(reader, machineEndian(), &tmp)
		array = append(array, tmp)
		tmp = 0
	}
	return array
}

func ReadInt16(datas []byte, nitems int32) []int16 {
	reader := bytes.NewReader(datas)
	var array []int16
	for i := int32(0); i < nitems; i++ {
		var tmp int16
		binary.Read(reader, machineEndian(), &tmp)
		array = append(array, tmp)
		tmp = 0
	}
	return array
}

func ReadInt32(datas []byte, nitems int32) []int32 {
	reader := bytes.NewReader(datas)
	var array []int32
	for i := int32(0); i < nitems; i++ {
		var tmp int32
		binary.Read(reader, machineEndian(), &tmp)
		array = append(array, tmp)
		tmp = 0
	}
	return array
}

func ReadFloat32(datas []byte, nitems int32) []float32 {
	reader := bytes.NewReader(datas)
	var array []float32
	for i := int32(0); i < nitems; i++ {
		var tmp float32
		binary.Read(reader, machineEndian(), &tmp)
		array = append(array, tmp)
		tmp = 0
	}
	return array
}

func WriteInt8(values []int8) []byte {
	var writer = new(bytes.Buffer)
	for i := 0; i < len(values); i++ {
		binary.Write(writer, machineEndian(), values[i])
	}
	return writer.Bytes()
}

func WriteInt16(values []int16) []byte {
	var writer = new(bytes.Buffer)
	for i := 0; i < len(values); i++ {
		binary.Write(writer, machineEndian(), values[i])
	}
	return writer.Bytes()
}

func WriteInt32(values []int32) []byte {
	var writer = new(bytes.Buffer)
	for i := 0; i < len(values); i++ {
		binary.Write(writer, machineEndian(), values[i])
	}
	return writer.Bytes()
}

func WriteFloat32(values []float32) []byte {
	var writer = new(bytes.Buffer)
	for i := 0; i < len(values); i++ {
		binary.Write(writer, machineEndian(), values[i])
	}
	return writer.Bytes()
}

func machineEndian() binary.ByteOrder {
	var x uint32 = 0x012345
	var ptr unsafe.Pointer = unsafe.Pointer(&x)

	if 0x01 == *((*byte)(ptr)) {
		return binary.BigEndian
	}
	return binary.LittleEndian
}
