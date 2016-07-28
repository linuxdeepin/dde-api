package powersupply

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestACOnline(t *testing.T) {
	Convey("ACOnline", t, func() {
		acExist, acOnline, err := ACOnline()
		t.Logf("acExist %v, acOnline %v, err %v", acExist, acOnline, err)
	})
}

func TestGetSystemBatteryInfos(t *testing.T) {
	Convey("GetSystemBatteryInfos", t, func() {
		batInfos, err := GetSystemBatteryInfos()
		if err != nil {
			t.Log("err", err)
			return
		}
		for _, batInfo := range batInfos {
			t.Logf("%+v", batInfo)
		}
	})
}
