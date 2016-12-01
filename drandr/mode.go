package drandr

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"math"
	"sort"
)

type ModeInfo struct {
	Id     uint32
	Width  uint16
	Height uint16
	Rate   float64
}
type ModeInfos []ModeInfo

func FindCommonModes(infosGroup ...ModeInfos) ModeInfos {
	countSet := make(map[string]int)
	tmpSet := make(map[string]ModeInfo)

	for _, infos := range infosGroup {
		sort.Sort(infos)
		for _, info := range infos.FilterBySize() {
			wh := fmt.Sprintf("%d%d", info.Width, info.Height)
			countSet[wh] += 1
			tmpSet[wh] = info
		}
	}

	for wh, count := range countSet {
		// remove not common mode
		if count < len(countSet) {
			delete(tmpSet, wh)
		}
	}

	var commons ModeInfos
	for _, info := range tmpSet {
		commons = append(commons, info)
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

func (infos ModeInfos) QueryBySize(width, height uint16) ModeInfo {
	for _, info := range infos {
		if info.Width == width && info.Height == height {
			return info
		}
	}
	return ModeInfo{}
}

func (infos ModeInfos) Best() ModeInfo {
	length := len(infos)
	if length == 0 {
		return ModeInfo{}
	}

	if length >= 2 {
		sort.Sort(infos)
	}
	return infos[0]
}

func (infos ModeInfos) Equal(list ModeInfos) bool {
	len1, len2 := len(infos), len(list)
	if len1 != len2 {
		return false
	}

	sort.Sort(infos)
	sort.Sort(list)
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

func (info ModeInfo) Equal(v ModeInfo) bool {
	return info.Id == v.Id
}

func toModeInfo(conn *xgb.Conn, info randr.ModeInfo) ModeInfo {
	return ModeInfo{
		Id:     uint32(info.Id),
		Width:  info.Width,
		Height: info.Height,
		Rate:   sumModeRate(info),
	}
}

func sumModeRate(info randr.ModeInfo) float64 {
	var vTotal = info.Vtotal
	if (info.ModeFlags & randr.ModeFlagDoubleScan) != 0 {
		vTotal *= 2
	}
	if (info.ModeFlags & randr.ModeFlagInterlace) != 0 {
		vTotal /= 2
	}

	var rate = float64(info.DotClock) / float64(uint32(info.Htotal)*uint32(vTotal))
	return (math.Floor(rate*10+0.5) / 10)
}
