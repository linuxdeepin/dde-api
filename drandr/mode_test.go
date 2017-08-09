package drandr

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCommonModes(t *testing.T) {
	convey.Convey("Test common modes", t, func() {
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
			convey.ShouldEqual(matches[i].Width, result[i].Width)
			convey.ShouldEqual(matches[i].Height, result[i].Height)
			convey.ShouldEqual(matches[i].Rate, result[i].Rate)
		}
	})
}
