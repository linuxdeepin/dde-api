package main

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"strings"
	"sync"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var genID = func() func() int32 {
	var lock sync.Mutex
	id := int32(0)
	return func() int32 {
		lock.Lock()
		tmp := id
		id += 1
		lock.Unlock()
		return tmp
	}
}()

func getIDList(x, y int32) ([]int32, []int32) {
	inList := []int32{}
	outList := []int32{}

	for id, array := range idRangeMap {
		inFlag := false
		for _, area := range array.areas {
			if isInArea(x, y, area) {
				inFlag = true
				if !isInIDList(id, inList) {
					inList = append(inList, id)
				}
			}
		}
		if !inFlag {
			if !isInIDList(id, outList) {
				outList = append(outList, id)
			}
		}
	}

	return inList, outList
}

func isInArea(x, y int32, area coordinateRange) bool {
	if (x >= area.X1 && x <= area.X2) &&
		(y >= area.Y1 && y <= area.Y2) {
		return true
	}

	return false
}

func isInIDList(id int32, list []int32) bool {
	for _, v := range list {
		if id == v {
			return true
		}
	}

	return false
}

var keyCode2Str = func() func(int32) string {
	XU, err := xgbutil.NewConn()
	if err != nil {
		Logger.Error("Can't connect to Xserver")
		return func(int32) string { return "" }
	}
	keybind.Initialize(XU)
	return func(code int32) string {
		keyStr := keybind.LookupString(XU, 0, xproto.Keycode(code))
		if keyStr == " " {
			keyStr = "space"
		}
		keyStr = strings.ToLower(keyStr)
		return keyStr
	}
}()

func buttonCode2str(code int32) string {
	switch code {
	case 1:
		return "LeftButton"
	case 2:
		return "Middlebutton"
	case 3:
		return "Rightbutton"
	case 4:
		return "RollForward"
	case 5:
		return "RollBack"
	}
	return "unknow button"
}
