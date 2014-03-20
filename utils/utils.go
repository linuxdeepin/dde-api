/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
        "os"
        "strings"
)

func deleteStartSpace(str string) string {
        if len(str) <= 0 {
                return ""
        }

        tmp := strings.TrimLeft(str, " ")

        return tmp
}

func (op *Manager) GetBaseName(path string) (string, bool) {
        if len(path) <= 0 {
                return "", false
        }

        as := strings.Split(path, "/")
        if l := len(as); l > 1 {
                return as[l-1], true
        }

        return "", false
}

func (op *Manager) IsContainFromStart(str, substr string) bool {
        l1 := len(substr)
        l2 := len(str)

        l := 0
        if l1 > l2 {
                l = l2
        } else {
                l = l1
        }

        for i := 0; i < l; i++ {
                if str[i] != substr[i] {
                        return false
                }
        }

        return true
}

func (op *Manager) IsFileExist(filename string) bool {
        if len(filename) <= 0 {
                return false
        }

        path, ok := op.URIToPath(filename)
        if !ok {
                return false
        }
        _, err := os.Stat(path)

        return err == nil || os.IsExist(err)
}

/*
 * t --> type
 *      0 : int64
 *      1 : string
 *      2 : byte
 */
func (op *Manager) IsElementExist(e interface{}, l interface{}) bool {
        switch e.(type) {
        case int32, uint32, int64, uint64:
                element := e.(int64)
                list := l.([]int64)
                for _, v := range list {
                        if element == v {
                                return true
                        }
                }
        case string:
                element := e.(string)
                list := l.([]string)
                for _, v := range list {
                        if element == v {
                                return true
                        }
                }
        case byte:
                element := e.(byte)
                list := l.([]byte)
                for _, v := range list {
                        if element == v {
                                return true
                        }
                }
        }

        return false
}

func (op *Manager) IsListEqual(l1, l2 interface{}) bool {
        switch l1.(type) {
        case []int32, []int64, []uint32:
                list1 := l1.([]int64)
                list2 := l2.([]int64)

                len1 := len(list1)
                len2 := len(list2)

                if len1 != len2 {
                        return false
                }

                for i := 0; i < len1; i++ {
                        if list1[i] != list2[i] {
                                return false
                        }
                }
        case []string:
                list1 := l1.([]string)
                list2 := l2.([]string)

                len1 := len(list1)
                len2 := len(list2)

                if len1 != len2 {
                        return false
                }

                for i := 0; i < len1; i++ {
                        if list1[i] != list2[i] {
                                return false
                        }
                }
        case []byte:
                list1 := l1.([]byte)
                list2 := l2.([]byte)

                len1 := len(list1)
                len2 := len(list2)

                if len1 != len2 {
                        return false
                }

                for i := 0; i < len1; i++ {
                        if list1[i] != list2[i] {
                                return false
                        }
                }
        }

        return true
}
