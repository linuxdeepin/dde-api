/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	C "gopkg.in/check.v1"
	"testing"
)

type testWrapper struct{}

func init() {
	C.Suite(&testWrapper{})
}

func Test(t *testing.T) {
	C.TestingT(t)
}

func (*testWrapper) TestLocaleFile(c *C.C) {
	finfo, err := NewLocaleFileInfo("testdata/locale.gen")
	c.Check(err, C.Equals, nil)
	c.Check(len(finfo.Infos), C.Equals, 471)
	c.Check(len(finfo.GetEnabledLocales()), C.Equals, 5)

	// test locale valid
	c.Check(finfo.IsLocaleValid("zh_CN.UTF-8"), C.Equals, true)
	c.Check(finfo.IsLocaleValid("zh_CNN"), C.Equals, false)

	// enable
	finfo.EnableLocale("zh_CN.UTF-8")
	c.Check(len(finfo.GetEnabledLocales()), C.Equals, 5)
	finfo.EnableLocale("zh_TW.UTF-8")
	c.Check(len(finfo.GetEnabledLocales()), C.Equals, 6)
	var tmp = "/tmp/test_locale"
	err = finfo.Save(tmp)
	c.Check(err, C.Equals, nil)

	finfo, err = NewLocaleFileInfo(tmp)
	c.Check(err, C.Equals, nil)
	c.Check(len(finfo.Infos), C.Equals, 471)
	c.Check(len(finfo.GetEnabledLocales()), C.Equals, 6)

	// disable
	finfo.DisableLocale("zh_HK.UTF-8")
	c.Check(len(finfo.GetEnabledLocales()), C.Equals, 6)
	finfo.DisableLocale("zh_CN.UTF-8")
	c.Check(len(finfo.GetEnabledLocales()), C.Equals, 5)
	var tmp2 = "/tmp/test_locale2"
	err = finfo.Save(tmp2)
	c.Check(err, C.Equals, nil)

	finfo, err = NewLocaleFileInfo(tmp2)
	c.Check(err, C.Equals, nil)
	c.Check(len(finfo.Infos), C.Equals, 471)
	c.Check(len(finfo.GetEnabledLocales()), C.Equals, 5)
}
