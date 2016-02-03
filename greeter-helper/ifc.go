/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import "fmt"

func (m *Manager) SetLayout(user, layout string) error {
	layout = formatLayout(layout)
	if len(layout) == 0 {
		return fmt.Errorf("Invalid layout: %s", layout)
	}
	return m.set(user, kfKeyLayout, layout)
}

func (m *Manager) SetLayoutList(user string, list []string) error {
	ret := formatLayoutList(list)
	if len(ret) == 0 {
		return fmt.Errorf("Invalid layout list: %v", list)
	}
	return m.set(user, kfKeyLayoutList, ret)
}

func (m *Manager) SetTheme(user, theme string) error {
	return m.set(user, kfKeyTheme, theme)
}
