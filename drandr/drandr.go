// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
