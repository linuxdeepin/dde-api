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
	"encoding/json"
	"strconv"
	"testing"

	"github.com/linuxdeepin/go-x11-client/ext/randr"
	"github.com/stretchr/testify/assert"
)

func Test_doFoundCommonModes(t *testing.T) {
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
	t.Run("Test_doFoundCommonModes", func(t *testing.T) {
		matches := doFoundCommonModes(infos2, infos1)
		for i := 0; i < len(matches); i++ {
			assert.Equal(t, matches[i].Width, result[i].Width)
			assert.Equal(t, matches[i].Height, result[i].Height)
			assert.Equal(t, matches[i].Rate, result[i].Rate)
		}
	})
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

func Test_Query(t *testing.T) {
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
	for i, data := range tests {
		t.Run("Test_Query"+strconv.Itoa(i), func(t *testing.T) {
			modeInfo := infos.Query(data.Id)
			assert.Equal(t, data.expected, modeInfo)
		})

	}
}

func Test_QueryBySize(t *testing.T) {
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
	for i, data := range tests {
		t.Run("Test_QueryBySize"+strconv.Itoa(i), func(t *testing.T) {
			modeInfos := infos.QueryBySize(data.width, data.height)
			assert.True(t, sliceModeInfosEq(data.expected, modeInfos))
		})
	}
}

func Test_FindCommonModes(t *testing.T) {
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
	t.Run("Test_FindCommonModes", func(t *testing.T) {
		matches := FindCommonModes(infos1, infos2)
		for i := 0; i < len(matches); i++ {
			assert.Equal(t, matches[i].Width, result[i].Width)
			assert.Equal(t, matches[i].Height, result[i].Height)
			assert.Equal(t, matches[i].Rate, result[i].Rate)
		}
		matches = FindCommonModes(infos1)
		assert.True(t, sliceModeInfosEq(matches, infos1))

		matches = FindCommonModes()
		assert.True(t, sliceModeInfosEq(matches, ModeInfos{}))
	})
}

func Test_Max(t *testing.T) {
	var modeInfos = ModeInfos{
		{
			Id:     70,
			Width:  1440,
			Height: 900,
			Rate:   60.1,
		},
		{
			Id:     74,
			Width:  1366,
			Height: 768,
			Rate:   59.0,
		},
		{
			Id:     71,
			Width:  1920,
			Height: 1080,
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

	info := modeInfos.Max()
	assert.Equal(t, ModeInfo{
		Id:     71,
		Width:  1920,
		Height: 1080,
		Rate:   60.1,
	}, info)

	modeInfos = make(ModeInfos, 0)
	info = modeInfos.Max()
	assert.Equal(t, uint16(0), info.Width)
	assert.Equal(t, uint16(0), info.Height)
	assert.Equal(t, float64(0), info.Rate)

	modeInfos = ModeInfos{
		{
			Id:     71,
			Width:  1920,
			Height: 1080,
			Rate:   60.1,
		},
	}
	info = modeInfos.Max()
	assert.Equal(t, ModeInfo{
		Id:     71,
		Width:  1920,
		Height: 1080,
		Rate:   60.1,
	}, info)

}

func Test_Equal(t *testing.T) {
	var modeInfos1 = ModeInfos{
		{
			Id:     70,
			Width:  1440,
			Height: 900,
			Rate:   60.1,
		},
		{
			Id:     74,
			Width:  1366,
			Height: 768,
			Rate:   59.0,
		},
		{
			Id:     71,
			Width:  1920,
			Height: 1080,
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
	var modeInfos2 = ModeInfos{
		{
			Id:     70,
			Width:  1440,
			Height: 900,
			Rate:   60.1,
		},
		{
			Id:     74,
			Width:  1366,
			Height: 768,
			Rate:   59.0,
		},
		{
			Id:     71,
			Width:  1920,
			Height: 1080,
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
	var modeInfos3 = ModeInfos{
		{
			Id:     74,
			Width:  1366,
			Height: 768,
			Rate:   59.0,
		},
		{
			Id:     71,
			Width:  1920,
			Height: 1080,
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
	var modeInfos4 = ModeInfos{
		{
			Id:     70,
			Width:  1440,
			Height: 900,
			Rate:   60.1,
		},
		{
			Id:     74,
			Width:  1366,
			Height: 720,
			Rate:   59.0,
		},
		{
			Id:     76,
			Width:  1920,
			Height: 1080,
			Rate:   60.1,
		},
		{
			Id:     72,
			Width:  1366,
			Height: 768,
			Rate:   60.2,
		},
		{
			Id:     75,
			Width:  800,
			Height: 600,
			Rate:   60.1,
		},
	}

	assert.True(t, modeInfos1.Equal(modeInfos2))
	assert.False(t, modeInfos1.Equal(modeInfos3))
	assert.False(t, modeInfos1.Equal(modeInfos4))

}

func Test_FilterBySize(t *testing.T) {
	var modeInfos = ModeInfos{
		{
			Id:     70,
			Width:  1440,
			Height: 900,
			Rate:   60.1,
		},
		{
			Id:     74,
			Width:  1366,
			Height: 768,
			Rate:   59.0,
		},
		{
			Id:     71,
			Width:  1920,
			Height: 1080,
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
		{
			Id:     76,
			Width:  800,
			Height: 600,
			Rate:   60.1,
		},
	}

	filterInfo := modeInfos.FilterBySize()
	assert.Equal(t, len(filterInfo), 4)
	assert.Equal(t, ModeInfo{
		Id:     70,
		Width:  1440,
		Height: 900,
		Rate:   60.1,
	}, filterInfo.Query(70))
	assert.Equal(t, ModeInfo{
		Id:     74,
		Width:  1366,
		Height: 768,
		Rate:   59.0,
	}, filterInfo.Query(74))
	assert.Equal(t, ModeInfo{
		Id:     71,
		Width:  1920,
		Height: 1080,
		Rate:   60.1,
	}, filterInfo.Query(71))
	assert.Equal(t, ModeInfo{
		Id:     75,
		Width:  800,
		Height: 600,
		Rate:   60.1,
	}, filterInfo.Query(75))

}

func Test_HasRefreshRate(t *testing.T) {
	var modeInfos = ModeInfos{
		{
			Id:     70,
			Width:  1440,
			Height: 900,
			Rate:   60.1,
		},
		{
			Id:     74,
			Width:  1366,
			Height: 768,
			Rate:   59.0,
		},
		{
			Id:     71,
			Width:  1920,
			Height: 1080,
			Rate:   60,
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
		{
			Id:     76,
			Width:  800,
			Height: 600,
			Rate:   60.1,
		},
	}

	assert.True(t, modeInfos.HasRefreshRate(60))
	assert.True(t, modeInfos.HasRefreshRate(60.1))
	assert.False(t, modeInfos.HasRefreshRate(75))

}

func Test_CalcModeRate(t *testing.T) {
	modeInfo1 := randr.ModeInfo{
		Id:         85,
		Width:      1920,
		Height:     1080,
		DotClock:   14850000,
		HSyncStart: 2448,
		HSyncEnd:   2492,
		HTotal:     2640,
		HSkew:      0,
		VSyncStart: 1084,
		VSyncEnd:   1089,
		VTotal:     1125,
		Name:       "1920x1080",
		ModeFlags:  1 << 4,
	}
	modeInfo2 := randr.ModeInfo{
		Id:         130,
		Width:      720,
		Height:     240,
		DotClock:   13514000,
		HSyncStart: 739,
		HSyncEnd:   801,
		HTotal:     858,
		HSkew:      0,
		VSyncStart: 244,
		VSyncEnd:   247,
		VTotal:     262,
		Name:       "720x240",
		ModeFlags:  1 << 5,
	}
	modeInfo3 := randr.ModeInfo{
		Id:         135,
		Width:      720,
		Height:     240,
		DotClock:   13514000,
		HSyncStart: 739,
		HSyncEnd:   801,
		HTotal:     0,
		HSkew:      0,
		VSyncStart: 244,
		VSyncEnd:   247,
		VTotal:     262,
		Name:       "720x240",
		ModeFlags:  1 << 5,
	}

	rate := calcModeRate(modeInfo1)
	assert.Equal(t, float64(10), rate)
	rate = calcModeRate(modeInfo2)
	assert.Equal(t, 30.058364027829676, rate)
	rate = calcModeRate(modeInfo3)
	assert.Equal(t, float64(0), rate)

}

func Test_String(t *testing.T) {
	modeInfo1 := ModeInfo{
		Id:     85,
		Width:  1920,
		Height: 1080,
		Rate:   60.1,
	}
	modeInfo2 := ModeInfo{
		Id:     130,
		Width:  720,
		Height: 240,
		Rate:   59.9,
	}
	modeInfo3 := ModeInfo{
		Id:     135,
		Width:  720,
		Height: 240,
		Rate:   60.0,
	}
	var infos ModeInfos
	infos = append(infos, modeInfo1)
	infos = append(infos, modeInfo2)
	infos = append(infos, modeInfo3)

	resultString := infos.String()
	expect, _ := json.Marshal(infos)
	expectString := string(expect)
	assert.Equal(t, expectString, resultString)

}

func Test_Len(t *testing.T) {
	modeInfo1 := ModeInfo{
		Id:     85,
		Width:  1920,
		Height: 1080,
		Rate:   60.1,
	}
	modeInfo2 := ModeInfo{
		Id:     130,
		Width:  720,
		Height: 240,
		Rate:   59.9,
	}
	modeInfo3 := ModeInfo{
		Id:     135,
		Width:  720,
		Height: 240,
		Rate:   60.0,
	}
	var infos ModeInfos
	infos = append(infos, modeInfo1)
	infos = append(infos, modeInfo2)
	infos = append(infos, modeInfo3)

	assert.Equal(t, 3, infos.Len())

}

func Test_Swap(t *testing.T) {
	modeInfo1 := ModeInfo{
		Id:     85,
		Width:  1920,
		Height: 1080,
		Rate:   60.1,
	}
	modeInfo2 := ModeInfo{
		Id:     130,
		Width:  720,
		Height: 240,
		Rate:   59.9,
	}
	modeInfo3 := ModeInfo{
		Id:     135,
		Width:  720,
		Height: 240,
		Rate:   60.0,
	}
	var infos ModeInfos
	infos = append(infos, modeInfo1)
	infos = append(infos, modeInfo2)
	infos = append(infos, modeInfo3)

	infos.Swap(0, 2)
	assert.Equal(t, modeInfo1, infos[2])
	assert.Equal(t, modeInfo3, infos[0])

}

func Test_toModeInfo(t *testing.T) {
	modeInfo1 := randr.ModeInfo{
		Id:         85,
		Width:      1920,
		Height:     1080,
		DotClock:   14850000,
		HSyncStart: 2448,
		HSyncEnd:   2492,
		HTotal:     2640,
		HSkew:      0,
		VSyncStart: 1084,
		VSyncEnd:   1089,
		VTotal:     1125,
		Name:       "1920x1080",
		ModeFlags:  1 << 4,
	}
	modeInfo2 := randr.ModeInfo{
		Id:         130,
		Width:      720,
		Height:     240,
		DotClock:   13487760,
		HSyncStart: 739,
		HSyncEnd:   801,
		HTotal:     858,
		HSkew:      0,
		VSyncStart: 244,
		VSyncEnd:   247,
		VTotal:     262,
		Name:       "720x240",
		ModeFlags:  1 << 5,
	}
	modeInfo3 := randr.ModeInfo{
		Id:         135,
		Width:      720,
		Height:     240,
		DotClock:   13514000,
		HSyncStart: 739,
		HSyncEnd:   801,
		HTotal:     0,
		HSkew:      0,
		VSyncStart: 244,
		VSyncEnd:   247,
		VTotal:     262,
		Name:       "720x240",
		ModeFlags:  1 << 5,
	}
	tests := []struct {
		Input  randr.ModeInfo
		Expect ModeInfo
	}{
		{
			Input: modeInfo1,
			Expect: ModeInfo{
				Id:     85,
				Width:  1920,
				Height: 1080,
				Rate:   10,
			},
		},
		{
			Input: modeInfo2,
			Expect: ModeInfo{
				Id:     130,
				Width:  720,
				Height: 240,
				Rate:   30,
			},
		},
		{
			Input: modeInfo3,
			Expect: ModeInfo{
				Id:     135,
				Width:  720,
				Height: 240,
				Rate:   0,
			},
		},
	}
	for i, test := range tests {
		t.Run("Test_toModeInfo"+strconv.Itoa(i), func(t *testing.T) {
			modeInfo := toModeInfo(test.Input)
			assert.Equal(t, test.Expect, modeInfo)
		})
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
