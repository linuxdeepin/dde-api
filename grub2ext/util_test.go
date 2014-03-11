package main

import (
	. "launchpad.net/gocheck"
	"testing"
)

type Util struct{}

var _ = Suite(&Util{})

func Test(t *testing.T) { TestingT(t) }

func (u *Util) TestGetImgClipRectByResolution(c *C) {
	var tests = []struct {
		sw, sh, iw, ih, x0, y0, x1, y1 int32
	}{
		{1024, 768, 1024, 768, 0, 0, 1024, 768},
		{1440, 900, 1920, 1080, 240, 90, 1680, 990},
		{1024, 768, 800, 600, 0, 0, 800, 600},
		{1024, 768, 500, 500, 0, 0, 500, 375},
	}
	for _, t := range tests {
		x0, y0, x1, y1 := getImgClipRectByResolution(uint16(t.sw), uint16(t.sh), t.iw, t.ih)
		c.Check(x0, Equals, t.x0)
		c.Check(y0, Equals, t.y0)
		c.Check(x1, Equals, t.x1)
		c.Check(y1, Equals, t.y1)
	}
	var iw, ih int32 = 1920, 1080
	for sw := 1; sw < 3000; sw += 5 {
		for sh := 1; sh < 3000; sh += 5 {
			x0, y0, x1, y1 := getImgClipRectByResolution(uint16(sw), uint16(sh), 1920, 1080)
			if (x1-x0) > iw || (y1-y0) > ih {
				c.Fatalf("sw=%d, sh=%d, iw=%d, ih=%d, x0=%d, y0=%d, x1=%d, y1=%d", sw, sh, iw, ih, x0, y0, x1, y1)
			}
		}
	}
}
