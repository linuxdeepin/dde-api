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
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAppInfos(t *testing.T) {
	Convey("Delete info", t, func() {
		var infos = AppInfos{
			&AppInfo{
				Id:   "gvim.desktop",
				Name: "gvim",
				Exec: "gvim",
			},
			&AppInfo{
				Id:   "firefox.desktop",
				Name: "Firefox",
				Exec: "firefox",
			}}
		So(len(infos.Delete("gvim.desktop")), ShouldEqual, 1)
		So(len(infos.Delete("vim.desktop")), ShouldEqual, 2)
	})
}

func TestUnmarshal(t *testing.T) {
	Convey("Test unmarsal", t, func() {
		table, err := unmarshal("testdata/data.json")
		So(err, ShouldBeNil)
		So(len(table.Apps), ShouldEqual, 2)

		So(table.Apps[0].AppId, ShouldEqual, "org.gnome.Nautilus.desktop")
		So(table.Apps[0].AppType, ShouldEqual, "file-manager")
		So(table.Apps[0].Types, ShouldResemble, []string{
			"inode/directory",
			"application/x-gnome-saved-search",
		})

		So(table.Apps[1].AppId, ShouldEqual, "org.gnome.gedit.desktop")
		So(table.Apps[1].AppType, ShouldEqual, "editor")
		So(table.Apps[1].Types, ShouldResemble, []string{
			"text/plain",
		})

	})
}

func TestMarshal(t *testing.T) {
	Convey("Marshal info", t, func() {
		content, err := marshal(&AppInfo{
			Id:   "gvim.desktop",
			Name: "gvim",
			Exec: "gvim",
		})
		So(err, ShouldBeNil)
		So(content, ShouldEqual,
			"{\"Id\":\"gvim.desktop\","+
				"\"Name\":\"gvim\","+
				"\"Exec\":\"gvim\"}")
	})

	Convey("Marshal info list", t, func() {
		content, err := marshal(AppInfos{
			&AppInfo{
				Id:   "gvim.desktop",
				Name: "gvim",
				Exec: "gvim",
			},
			&AppInfo{
				Id:   "firefox.desktop",
				Name: "Firefox",
				Exec: "firefox",
			},
		})
		So(err, ShouldBeNil)
		So(content, ShouldEqual, "["+
			"{\"Id\":\"gvim.desktop\","+
			"\"Name\":\"gvim\","+
			"\"Exec\":\"gvim\"},"+
			"{\"Id\":\"firefox.desktop\","+
			"\"Name\":\"Firefox\","+
			"\"Exec\":\"firefox\"}"+
			"]")
	})

	Convey("Marshal nil", t, func() {
		content, err := marshal(nil)
		So(content, ShouldEqual, "null")
		So(err, ShouldBeNil)
	})
}

func TestIsStrInList(t *testing.T) {
	Convey("Test str whether in list", t, func() {
		var list = []string{"abc", "abs"}
		So(isStrInList("abs", list), ShouldEqual, true)
		So(isStrInList("abd", list), ShouldEqual, false)
	})
}

func TestDelStrFromList(t *testing.T) {
	Convey("Test delete str from list", t, func() {
		var list = []string{"abc", "abs"}
		ret, deleted := delStrFromList("abs", list)
		So(deleted, ShouldEqual, true)
		So(ret, ShouldResemble, []string{"abc"})

		ret, deleted = delStrFromList("abd", list)
		So(deleted, ShouldEqual, false)
		So(ret, ShouldResemble, list)
	})
}
