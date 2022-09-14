// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
		err := binary.Read(reader, machineEndian(), &tmp)
		if err != nil {
			return nil
		}
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
		err := binary.Read(reader, machineEndian(), &tmp)
		if err != nil {
			return nil
		}
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
		err := binary.Read(reader, machineEndian(), &tmp)
		if err != nil {
			return nil
		}
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
		err := binary.Read(reader, machineEndian(), &tmp)
		if err != nil {
			return nil
		}
		array = append(array, tmp)
		tmp = 0
	}
	return array
}

func WriteInt8(values []int8) []byte {
	var writer = new(bytes.Buffer)
	for i := 0; i < len(values); i++ {
		err := binary.Write(writer, machineEndian(), values[i])
		if err != nil {
			return nil
		}
	}
	return writer.Bytes()
}

func WriteInt16(values []int16) []byte {
	var writer = new(bytes.Buffer)
	for i := 0; i < len(values); i++ {
		err := binary.Write(writer, machineEndian(), values[i])
		if err != nil {
			return nil
		}
	}
	return writer.Bytes()
}

func WriteInt32(values []int32) []byte {
	var writer = new(bytes.Buffer)
	for i := 0; i < len(values); i++ {
		err := binary.Write(writer, machineEndian(), values[i])
		if err != nil {
			return nil
		}
	}
	return writer.Bytes()
}

func WriteFloat32(values []float32) []byte {
	var writer = new(bytes.Buffer)
	for i := 0; i < len(values); i++ {
		err := binary.Write(writer, machineEndian(), values[i])
		if err != nil {
			return nil
		}
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
