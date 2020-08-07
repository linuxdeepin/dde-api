/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package drandr

import (
	"fmt"

	"github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/randr"
)

type randrIdList []uint32

type OutputInfo struct {
	Name       string
	Id         uint32
	MmWidth    uint32
	MmHeight   uint32
	Crtc       CrtcInfo
	Connection bool
	Timestamp  x.Timestamp

	EDID []byte

	Clones         randrIdList
	Crtcs          randrIdList
	Modes          randrIdList
	PreferredModes randrIdList
}
type OutputInfos []OutputInfo

func (infos OutputInfos) Query(id uint32) OutputInfo {
	return infos.query("id", fmt.Sprintf("%v", id))
}

func (infos OutputInfos) QueryByName(name string) OutputInfo {
	return infos.query("name", name)
}

func (infos OutputInfos) ListNames() []string {
	var names []string
	for _, info := range infos {
		names = append(names, info.Name)
	}
	return names
}

func (infos OutputInfos) ListConnectionOutputs() OutputInfos {
	var ret OutputInfos
	for _, info := range infos {
		if !info.Connection {
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func (infos OutputInfos) query(key, value string) OutputInfo {
	for _, info := range infos {
		if key == "id" {
			if fmt.Sprintf("%d", info.Id) == value {
				return info
			}
		} else if key == "name" {
			if info.Name == value {
				return info
			}
		}
	}
	return OutputInfo{}
}

func toOutputInfo(conn *x.Conn, output randr.Output) OutputInfo {
	reply, err := randr.GetOutputInfo(conn, output, lastConfigTimestamp).Reply(conn)
	if err != nil {
		return OutputInfo{}
	}
	var info = OutputInfo{
		Name:       string(reply.Name),
		Id:         uint32(output),
		MmWidth:    reply.MmWidth,
		MmHeight:   reply.MmHeight,
		Connection: reply.Connection == randr.ConnectionConnected,
		Timestamp:  reply.Timestamp,
		Clones:     outputsToRandrIdList(reply.Clones),
		Crtcs:      crtcsToRandrIdList(reply.Crtcs),
		Modes:      modesToRandrIdList(reply.Modes),
	}
	info.PreferredModes = getOutputPreferredModes(info.Modes, reply.NumPreferred)
	info.EDID, _ = getOutputEDID(conn, output)
	info.Crtc = toCrtcInfo(conn, reply.Crtc)

	return info
}

func outputsToRandrIdList(outputs []randr.Output) randrIdList {
	var list randrIdList
	for _, output := range outputs {
		list = append(list, uint32(output))
	}
	return list
}

func crtcsToRandrIdList(crtcs []randr.Crtc) randrIdList {
	var list randrIdList
	for _, crtc := range crtcs {
		list = append(list, uint32(crtc))
	}
	return list
}

func modesToRandrIdList(modes []randr.Mode) randrIdList {
	var list randrIdList
	for _, mode := range modes {
		list = append(list, uint32(mode))
	}
	return list
}

func getOutputEDID(conn *x.Conn, output randr.Output) ([]byte, error) {
	atomEDID, err := conn.GetAtom("EDID")
	if err != nil {
		return nil, err
	}

	reply, err := randr.GetOutputProperty(conn, output,
		atomEDID, x.AtomInteger,
		0, 32, false, false).Reply(conn)
	if err != nil {
		return nil, err
	}
	return reply.Value, nil
}

func getOutputPreferredModes(modes randrIdList, nPreferred uint16) randrIdList {
	if nPreferred == 0 || uint16(len(modes)) < nPreferred {
		return nil
	}
	return modes[:nPreferred]
}