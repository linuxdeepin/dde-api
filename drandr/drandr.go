package drandr

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"sync"
)

var (
	inited     bool
	infoLocker sync.Mutex

	lastConfigTimestamp xproto.Timestamp
)

func GetScreenInfo(conn *xgb.Conn) (OutputInfos, ModeInfos, error) {
	infoLocker.Lock()
	defer infoLocker.Unlock()

	var (
		outputInfos OutputInfos
		modeInfos   ModeInfos
	)
	if !inited {
		err := randr.Init(conn)
		if err != nil {
			return outputInfos, modeInfos, err
		}
		inited = true
	}

	sinfo := xproto.Setup(conn).DefaultScreen(conn)
	resource, err := randr.GetScreenResources(conn, sinfo.Root).Reply()
	if err != nil {
		return outputInfos, modeInfos, err
	}

	lastConfigTimestamp = resource.ConfigTimestamp
	for _, output := range resource.Outputs {
		info := toOuputInfo(conn, output)
		if len(info.Name) == 0 {
			continue
		}

		outputInfos = append(outputInfos, info)
	}

	for _, mode := range resource.Modes {
		modeInfos = append(modeInfos, toModeInfo(conn, mode))
	}
	return outputInfos, modeInfos, nil
}
