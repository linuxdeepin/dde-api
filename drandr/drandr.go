package drandr

import (
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"sort"
	"sync"
)

var (
	inited     bool
	infoLocker sync.Mutex

	lastConfigTimestamp xproto.Timestamp
)

type ScreenInfo struct {
	Outputs OutputInfos
	Modes   ModeInfos

	conn   *xgb.Conn
	window xproto.Window
}

func GetScreenInfo(conn *xgb.Conn) (*ScreenInfo, error) {
	infoLocker.Lock()
	defer infoLocker.Unlock()

	if !inited {
		err := randr.Init(conn)
		if err != nil {
			return nil, err
		}
		inited = true
	}

	sinfo := xproto.Setup(conn).DefaultScreen(conn)
	resource, err := randr.GetScreenResources(conn, sinfo.Root).Reply()
	if err != nil {
		return nil, err
	}

	var outputInfos OutputInfos
	lastConfigTimestamp = resource.ConfigTimestamp
	for _, output := range resource.Outputs {
		info := toOuputInfo(conn, output)
		if len(info.Name) == 0 {
			continue
		}

		outputInfos = append(outputInfos, info)
	}

	var modeInfos ModeInfos
	for _, mode := range resource.Modes {
		modeInfos = append(modeInfos, toModeInfo(conn, mode))
	}
	sort.Sort(modeInfos)
	return &ScreenInfo{
		Outputs: outputInfos,
		Modes:   modeInfos,
		conn:    conn,
		window:  sinfo.Root,
	}, nil
}

func (info *ScreenInfo) GetPrimary() (*OutputInfo, error) {
	reply, err := randr.GetOutputPrimary(info.conn, info.window).Reply()
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
	sinfo := xproto.Setup(info.conn).DefaultScreen(info.conn)
	return sinfo.WidthInPixels, sinfo.HeightInPixels
}
