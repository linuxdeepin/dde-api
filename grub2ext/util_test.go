package main

import (
	. "launchpad.net/gocheck"
	"testing"
)

type Util struct{}

var _ = Suite(&Util{})

func Test(t *testing.T) { TestingT(t) }

func (u *Util) TestGetImgClipSizeByResolution(c *C) {
	var tests = []struct {
		sw, sh, iw, ih, wantw, wanth int32
	}{
		{1024, 768, 1024, 768, 1024, 768},
		{1024, 768, 1920, 1080, 1024, 768},
	}
	for _, t := range tests {
		w, h := getImgClipSizeByResolution(uint16(t.sw), uint16(t.sh), t.iw, t.ih)
		c.Check(w, Equals, t.wantw)
		c.Check(h, Equals, t.wanth)
	}
	var iw, ih int32 = 1920, 1080
	for sw := 1; sw < 3000; sw += 5 {
		for sh := 1; sh < 3000; sh += 5 {
			w, h := getImgClipSizeByResolution(uint16(sw), uint16(sh), 1920, 1080)
			if w > iw || h > ih {
				c.Fatalf("sw=%d, sh=%d, iw=%d, ih=%d, w=%d, h=%d", sw, sh, iw, ih, w, h)
			}
		}
	}
}
