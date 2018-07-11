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
	"fmt"

	"github.com/linuxdeepin/go-x11-client/ext/randr"
)

type ModeInfo struct {
	Id     uint32
	Width  uint16
	Height uint16
	Rate   float64
}
type ModeInfos []ModeInfo

func FindCommonModes(infosGroup ...ModeInfos) ModeInfos {
	length := len(infosGroup)
	if length == 0 {
		return ModeInfos{}
	} else if length == 1 {
		return infosGroup[0]
	}

	var commons = infosGroup[0]
	for i := 1; i < length; i++ {
		commons = doFoundCommonModes(commons, infosGroup[i])
	}
	return commons
}

func (infos ModeInfos) Query(id uint32) ModeInfo {
	for _, info := range infos {
		if info.Id == id {
			return info
		}
	}
	return ModeInfo{}
}

func (infos ModeInfos) QueryBySize(width, height uint16) ModeInfos {
	var matches ModeInfos
	for _, info := range infos {
		if info.Width == width && info.Height == height {
			matches = append(matches, info)
		}
	}
	return matches
}

func (infos ModeInfos) Max() ModeInfo {
	length := len(infos)
	if length == 0 {
		return ModeInfo{}
	} else if length == 1 {
		return infos[0]
	}

	var idx = 0
	for i := 1; i < length; i++ {
		if !infos.Less(idx, i) {
			idx = i
		}
	}
	return infos[idx]
}

func (infos ModeInfos) Equal(list ModeInfos) bool {
	len1, len2 := len(infos), len(list)
	if len1 != len2 {
		return false
	}

	for i := 0; i < len1; i++ {
		if !infos[i].Equal(list[i]) {
			return false
		}
	}
	return true
}

func (infos ModeInfos) String() string {
	data, _ := json.Marshal(infos)
	return string(data)
}

func (infos ModeInfos) Len() int {
	return len(infos)
}

func (infos ModeInfos) Less(i, j int) bool {
	if infos[i].Width == infos[j].Width &&
		infos[i].Height == infos[j].Height {
		return infos[i].Rate > infos[j].Rate
	}

	sum1 := infos[i].Width + infos[j].Height
	sum2 := infos[j].Width + infos[j].Height
	if sum1 != sum2 {
		return sum1 > sum2
	}

	if infos[i].Width == infos[j].Width {
		return infos[i].Height > infos[j].Height
	}
	return infos[i].Width > infos[j].Height
}

func (infos ModeInfos) Swap(i, j int) {
	infos[i], infos[j] = infos[j], infos[i]
}

func (infos ModeInfos) FilterBySize() ModeInfos {
	var set = make(map[string]ModeInfo)
	for _, info := range infos {
		wh := fmt.Sprintf("%d%d", info.Width, info.Height)
		if _, ok := set[wh]; ok {
			continue
		}
		set[wh] = info
	}

	var ret ModeInfos
	for _, info := range set {
		ret = append(ret, info)
	}
	return ret
}

func (infos ModeInfos) HasRefreshRate(rate float64) bool {
	for _, info := range infos {
		if info.Rate == rate {
			return true
		}
	}
	return false
}

func (info ModeInfo) Equal(v ModeInfo) bool {
	return info.Id == v.Id
}

func toModeInfo(info randr.ModeInfo) ModeInfo {
	return ModeInfo{
		Id:     uint32(info.Id),
		Width:  info.Width,
		Height: info.Height,
		Rate:   calcModeRate(info),
	}
}

func calcModeRate(info randr.ModeInfo) float64 {
	vTotal := float64(info.VTotal)
	if (info.ModeFlags & randr.ModeFlagDoubleScan) != 0 {
		/* doublescan doubles the number of lines */
		vTotal *= 2
	}
	if (info.ModeFlags & randr.ModeFlagInterlace) != 0 {
		/* interlace splits the frame into two fields */
		/* the field rate is what is typically reported by monitors */
		vTotal /= 2
	}

	if info.HTotal == 0 || vTotal == 0 {
		return 0
	} else {
		return float64(info.DotClock) / (float64(info.HTotal) * vTotal)
	}
}

// doFoundCommonModes return common modes sorted by x11 preferred
func doFoundCommonModes(modes1, modes2 ModeInfos) ModeInfos {
	var (
		common   ModeInfos
		max, min = modes1, modes2
	)
	if max[0].Width+max[0].Height < min[0].Width+min[0].Height {
		max, min = modes2, modes1
	}
	for _, mode := range min {
		matches := max.QueryBySize(mode.Width, mode.Height)
		if len(matches) == 0 {
			continue
		}

		// filter same mode
		if v := common.QueryBySize(matches[0].Width, matches[0].Height); len(v) != 0 {
			continue
		}
		common = append(common, matches[0])
	}
	return common
}
