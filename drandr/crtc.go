// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package drandr

import (
	"github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/randr"
)

type CrtcInfo struct {
	Id       uint32 // if crtc == 0, means output closed or disconnection
	Mode     uint32
	X        int16
	Y        int16
	Width    uint16
	Height   uint16
	Rotation uint16
	Reflect  uint16

	Rotations []uint16
	Reflects  []uint16
}

func toCrtcInfo(conn *x.Conn, crtc randr.Crtc) CrtcInfo {
	reply, err := randr.GetCrtcInfo(conn, crtc, lastConfigTimestamp).Reply(conn)
	if err != nil {
		return CrtcInfo{}
	}
	var info = CrtcInfo{
		Id:        uint32(crtc),
		X:         reply.X,
		Y:         reply.Y,
		Mode:      uint32(reply.Mode),
		Width:     reply.Width,
		Height:    reply.Height,
		Rotations: getRotations(reply.Rotations),
		Reflects:  getReflects(reply.Rotations),
	}
	info.Rotation, info.Reflect = parseCrtcRotation(reply.Rotation)
	return info
}

func parseCrtcRotation(origin uint16) (rotation, reflect uint16) {
	rotation = origin & 0xf
	reflect = origin & 0xf0

	switch rotation {
	case 1, 2, 4, 8:
		break
	default:
		//Invalid rotation value
		rotation = 1
	}

	switch reflect {
	case 0, 16, 32, 48:
		break
	default:
		// Invalid reflect value
		reflect = 0
	}

	return
}

func getRotations(origin uint16) []uint16 {
	var ret []uint16

	if origin&randr.RotationRotate0 == randr.RotationRotate0 {
		ret = append(ret, randr.RotationRotate0)
	}
	if origin&randr.RotationRotate90 == randr.RotationRotate90 {
		ret = append(ret, randr.RotationRotate90)
	}
	if origin&randr.RotationRotate180 == randr.RotationRotate180 {
		ret = append(ret, randr.RotationRotate180)
	}
	if origin&randr.RotationRotate270 == randr.RotationRotate270 {
		ret = append(ret, randr.RotationRotate270)
	}
	return ret
}

func getReflects(origin uint16) []uint16 {
	var ret = []uint16{0}

	if origin&randr.RotationReflectX == randr.RotationReflectX {
		ret = append(ret, randr.RotationReflectX)
	}
	if origin&randr.RotationReflectY == randr.RotationReflectY {
		ret = append(ret, randr.RotationReflectY)
	}
	if len(ret) == 3 {
		ret = append(ret, randr.RotationReflectX|randr.RotationReflectY)
	}
	return ret
}
