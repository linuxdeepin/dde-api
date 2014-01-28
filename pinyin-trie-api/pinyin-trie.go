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
	"dlib/logger"
	"strconv"
	"strings"
)

type PinyinTrie struct{}

type TrieInfo struct {
	Pinyins []string
	Key     string
	Value   string
}

const (
	PINYIN_TRIE_PATH = "/com/deepin/dde/api/PinyinTrie"
	PINYIN_TRIE_IFC  = "com.deepin.dde.api.PinyinTrie"
)

var nameMD5Map map[string]string

func (t *PinyinTrie) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		PINYIN_DEST,
		PINYIN_TRIE_PATH,
		PINYIN_TRIE_IFC,
	}
}

func (t *PinyinTrie) NewTrieWithString(values map[string]string, name string) string {
	md5Byte := md5.Sum([]byte(getStringFromArray(values)))
	logger.Println("MD5: ", md5Byte)
	if len(md5Byte) == 0 {
		return ""
	}

	md5Str := getMD5String(md5Byte)
	if isMd5Exist(md5Str) {
		return md5Str
	}

	if isNameExist(name) {
		str, _ := nameMD5Map[name]
		t.DestroyTrie(str)
	}
	nameMD5Map[name] = md5Str

	root := newTrie()
	if values == nil {
		return ""
	}
	go func() {
		infos := getPinyinArray(values)
		strsMD5Map[md5Str] = infos
		root.insertTrieInfo(infos)
		trieMD5Map[md5Str] = root
	}()
	return md5Str
}

/*
func (t *PinyinTrie) TraversalTrie(str string) {
	root := trieMD5Map[str]
	root.traversalTrie()
}
*/

func (t *PinyinTrie) SearchKeys(keys string, str string) []string {
	root, ok := trieMD5Map[str]
	if !ok {
		return nil
	}
	keys = strings.ToLower(keys)
	rets := root.searchTrie(keys)
	tmp := searchKeyFromString(keys, str)
	for _, v := range tmp {
		if !isIdExist(v, rets) {
			rets = append(rets, v)
		}
	}

	return rets
}

func (t *PinyinTrie) DestroyTrie(md5Str string) {
	/*
		root, ok := trieMD5Map[md5Str]
		if !ok {
			return
		}
	*/
	delete(trieMD5Map, md5Str)
	delete(strsMD5Map, md5Str)
}

/*
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
*/

func getStringFromArray(strs map[string]string) string {
	str := ""

	for _, v := range strs {
		str += v
	}

	return str
}

func getPinyinArray(strs map[string]string) []*TrieInfo {
	rets := []*TrieInfo{}
	for k, v := range strs {
		array := getPinyinFromKey(v)
		v = strings.ToLower(v)
		tmp := &TrieInfo{Pinyins: array, Key: k, Value: v}
		rets = append(rets, tmp)
	}

	return rets
}

func getMD5String(bytes [16]byte) string {
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

func searchKeyFromString(key, md5Str string) []string {
	rets := []string{}

	infos := strsMD5Map[md5Str]
	for _, v := range infos {
		if strings.Contains(v.Value, key) {
			rets = append(rets, v.Key)
		}
	}

	return rets
}

func isIdExist(id string, list []string) bool {
	for _, v := range list {
		if v == id {
			return true
		}
	}

	return false
}

func isMd5Exist(md5Str string) bool {
	_, ok := strsMD5Map[md5Str]
	if ok {
		return true
	}

	return false
}

func isNameExist(name string) bool {
	_, ok := nameMD5Map[name]
	if ok {
		return true
	}

	return false
}
