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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommonModes(t *testing.T) {
	var infos1 = ModeInfos{
		{
			Id:     71,
			Width:  1920,
			Height: 1080,
			Rate:   60.1,
		},
		{
			Id:     70,
			Width:  1440,
			Height: 900,
			Rate:   60.1,
		},
		{
			Id:     72,
			Width:  1366,
			Height: 768,
			Rate:   60.1,
		},
		{
			Id:     74,
			Width:  1366,
			Height: 768,
			Rate:   59.0,
		},
		{
			Id:     75,
			Width:  800,
			Height: 600,
			Rate:   60.1,
		},
	}
	var infos2 = ModeInfos{
		{
			Id:     71,
			Width:  1440,
			Height: 900,
			Rate:   60.1,
		},
		{
			Id:     72,
			Width:  1366,
			Height: 768,
			Rate:   60.1,
		},
		{
			Id:     73,
			Width:  1366,
			Height: 768,
			Rate:   59.0,
		},
		{
			Id:     75,
			Width:  800,
			Height: 600,
			Rate:   60.1,
		},
	}
	var result = ModeInfos{
		{
			Id:     71,
			Width:  1440,
			Height: 900,
			Rate:   60.1,
		},
		{
			Id:     72,
			Width:  1366,
			Height: 768,
			Rate:   60.1,
		},
		{
			Id:     75,
			Width:  800,
			Height: 600,
			Rate:   60.1,
		},
	}

	matches := doFoundCommonModes(infos1, infos2)
	for i := 0; i < len(matches); i++ {
		assert.Equal(t, matches[i].Width, result[i].Width)
		assert.Equal(t, matches[i].Height, result[i].Height)
		assert.Equal(t, matches[i].Rate, result[i].Rate)
	}
}

var infos = ModeInfos{
	{
		0,
		1920,
		1080,
		59.9,
	},
	{
		1,
		1440,
		720,
		60.1,
	},
	{
		2,
		1600,
		900,
		75,
	},
}

func TestQuery(t *testing.T) {
	tests := []struct {
		Id       uint32
		expected ModeInfo
	}{
		{
			0,
			ModeInfo{
				0,
				1920,
				1080,
				59.9,
			},
		},
		{
			1,
			ModeInfo{
				1,
				1440,
				720,
				60.1,
			},
		},
		{
			2,
			ModeInfo{
				2,
				1600,
				900,
				75,
			},
		},
	}
	for _, data := range tests {
		modeInfo := infos.Query(data.Id)
		assert.Equal(t, data.expected, modeInfo)
	}
}

func TestQueryBySize(t *testing.T) {
	tests := []struct {
		width    uint16
		height   uint16
		expected ModeInfos
	}{
		{
			1920,
			1080,
			ModeInfos{
				{
					0,
					1920,
					1080,
					59.9,
				},
			},
		},
		{
			1440,
			720,
			ModeInfos{
				{
					1,
					1440,
					720,
					60.1,
				},
			},
		},
		{
			1600,
			900,
			ModeInfos{
				{
					2,
					1600,
					900,
					75},
			},
		},
	}
	for _, data := range tests {
		modeInfos := infos.QueryBySize(data.width, data.height)
		assert.True(t, sliceModeInfosEq(data.expected, modeInfos))
	}
}

func sliceModeInfosEq(a, b []ModeInfo) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
