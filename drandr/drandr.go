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
	"sync"

	"github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/randr"
)

var (
	infoLocker sync.Mutex

	lastConfigTimestamp x.Timestamp
)

type ScreenInfo struct {
	Outputs OutputInfos
	Modes   ModeInfos

	conn   *x.Conn
	window x.Window
}

func GetScreenInfo(conn *x.Conn) (*ScreenInfo, error) {
	infoLocker.Lock()
	defer infoLocker.Unlock()

	root := conn.GetDefaultScreen().Root
	resource, err := randr.GetScreenResources(conn, root).Reply(conn)
	if err != nil {
		return nil, err
	}

	var outputInfos OutputInfos
	lastConfigTimestamp = resource.ConfigTimestamp
	for _, output := range resource.Outputs {
		info := toOutputInfo(conn, output)
		if len(info.Name) == 0 {
			continue
		}

		outputInfos = append(outputInfos, info)
	}

	var modeInfos ModeInfos
	for _, mode := range resource.Modes {
		modeInfos = append(modeInfos, toModeInfo(mode))
	}
	return &ScreenInfo{
		Outputs: outputInfos,
		Modes:   modeInfos,
		conn:    conn,
		window:  root,
	}, nil
}

func (info *ScreenInfo) GetPrimary() (*OutputInfo, error) {
	reply, err := randr.GetOutputPrimary(info.conn, info.window).Reply(info.conn)
	if err != nil {
		return nil, err
	}

	output := info.Outputs.Query(uint32(reply.Output))
	if len(output.Name) == 0 {
		return nil, fmt.Errorf("No primary found for %v", reply.Output)
	}
	return &output, nil
}

func (info *ScreenInfo) GetScreenSize() (uint16, uint16) {
	screen := info.conn.GetDefaultScreen()
	return screen.WidthInPixels, screen.HeightInPixels
}
