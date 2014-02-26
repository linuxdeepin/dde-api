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

func (op *Manager) ConvertURIToPath(uri string, t int32) (string, bool) {
        tmp := deleteStartSpace(uri)
        switch t {
        case URI_TYPE_FILE:
                if !op.IsContainFromStart(tmp, URI_STRING_FILE) {
                        return "", false
                }
                return tmp[7:], true
        case URI_TYPE_FTP:
                if !op.IsContainFromStart(tmp, URI_STRING_FTP) {
                        return "", false
                }
                return tmp[6:], true
        case URI_TYPE_HTTP:
                if !op.IsContainFromStart(tmp, URI_STRING_HTTP) {
                        return "", false
                }
                return tmp[7:], true
        case URI_TYPE_HTTPS:
                if !op.IsContainFromStart(tmp, URI_STRING_HTTPS) {
                        return "", false
                }
                return tmp[8:], true
        case URI_TYPE_SMB:
                if !op.IsContainFromStart(tmp, URI_STRING_SMB) {
                        return "", false
                }
                return tmp[6:], true
        }

        return "", false
}

func (op *Manager) ConvertPathToURI(path string, t int32) (string, bool) {
        switch t {
        case URI_TYPE_FILE:
                return URI_STRING_FILE + path, true
        case URI_TYPE_FTP:
                return URI_STRING_FTP + path, true
        case URI_TYPE_HTTP:
                return URI_STRING_HTTP + path, true
        case URI_TYPE_HTTPS:
                return URI_STRING_HTTPS + path, true
        case URI_TYPE_SMB:
                return URI_STRING_SMB + path, true
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

        _, err := os.Stat(filename)
        if os.IsNotExist(err) {
                return false
        }

        return true
}

/*
 * t --> type
 *      0 : int64
 *      1 : string
 *      2 : byte
 */
func (op *Manager) IsElementExist(e interface{}, l interface{}, t int32) bool {
        switch t {
        case ELEMENT_TYPE_INT:
                element := e.(int64)
                list := l.([]int64)
                for _, v := range list {
                        if element == v {
                                return true
                        }
                }
        case ELEMENT_TYPE_STRING:
                element := e.(string)
                list := l.([]string)
                for _, v := range list {
                        if element == v {
                                return true
                        }
                }
        case ELEMENT_TYPE_BYTE:
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

func (op *Manager) IsListEqual(l1, l2 interface{}, t int32) bool {
        switch t {
        case ELEMENT_TYPE_INT:
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
        case ELEMENT_TYPE_STRING:
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
        case ELEMENT_TYPE_BYTE:
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
