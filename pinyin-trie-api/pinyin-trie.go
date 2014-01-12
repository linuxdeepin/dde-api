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
	"crypto/md5"
	"dlib/dbus"
	"fmt"
	"strconv"
)

type PinyinTrie struct{}

type TrieInfo struct {
	Pinyins []string
	Value   string
}

const (
	PINYIN_TRIE_PATH = "/com/deepin/dde/api/PinyinTrie"
	PINYIN_TRIE_IFC  = "com.deepin.dde.api.PinyinTrie"
)

func (t *PinyinTrie) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		PINYIN_DEST,
		PINYIN_TRIE_PATH,
		PINYIN_TRIE_IFC,
	}
}

func (t *PinyinTrie) NewTrieWithString(values []string) string {
	root := NewTrie()
	if values == nil {
                return ""
	}

	md5Byte := md5.Sum([]byte(GetStringFromArray(values)))
	fmt.Println("MD5: ", md5Byte)
	if len(md5Byte) == 0 {
		return ""
	}
	md5Str := GetMD5String(md5Byte)
	infos := GetPinyinArray(values)
	strsMD5Map[md5Str] = infos
	root.InsertTrieInfo(infos)
	trieMD5Map[md5Str] = root
	fmt.Println(md5Str)
	return md5Str
}

func (t *PinyinTrie) TraversalTrie(str string) {
	root := trieMD5Map[str]
	root.TraversalTrie()
}

func (t *PinyinTrie) SearchTrieWithString(keys string,
	str string) []string {
	root, ok := trieMD5Map[str]
	if !ok {
		return nil
	}
	rets := root.SearchTrie(keys)

	retStrs := []string{}
	infos := strsMD5Map[str]
	for _, v := range rets {
		retStrs = append(retStrs, infos[v].Value)
	}

	return retStrs
}

func GetStringFromArray(strs []string) string {
	str := ""

	for i, _ := range strs {
		str += strs[i]
	}

	return str
}

func GetPinyinArray(strs []string) []*TrieInfo {
	rets := []*TrieInfo{}
	for _, k := range strs {
		array := GetPinyinFromKey(k)
		tmp := &TrieInfo{Pinyins: array, Value: k}
		rets = append(rets, tmp)
	}

	return rets
}

func GetMD5String(bytes [16]byte) string {
	str := ""

	for _, v := range bytes {
		s := strconv.FormatInt(int64(v), 16)
		if len(s) == 1 {
			str += "0" + s
		} else {
			str += s
		}
	}

	return str
}
