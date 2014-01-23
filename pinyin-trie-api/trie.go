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
	"dlib/logger"
	"strings"
)

type Trie struct {
	Key        byte
	IndexArray []int32
	NextNode   [_TRIE_CHILD_LEN]*Trie
}

const (
	_TRIE_CHILD_LEN = 26
)

var (
	trieMD5Map map[string]*Trie
	strsMD5Map map[string][]*TrieInfo
)

func getNode(ch byte) *Trie {
	node := new(Trie)
	node.Key = ch
	return node
}

func newTrie() *Trie {
	root := getNode(' ')
	return root
}

func (root *Trie) insertTrieInfo(values []*TrieInfo) {
	for i, v := range values {
		root.insertStringArray(v.Pinyins, int32(i))
	}
}

func (root *Trie) insertStringArray(strs []string, pos int32) {
	if strs == nil {
		return
	}

	for _, v := range strs {
		root.insertString(v, pos)
	}
}

func (root *Trie) insertString(str string, pos int32) {
	if l := len(str); l == 0 {
		return
	}
	low := strings.ToLower(str)
	curNode := root

	logger.Println("keys: ", low)
	for i, _ := range str {
		index := low[i] - 'a'
		if curNode.NextNode[index] == nil {
			curNode.NextNode[index] = getNode(low[i])
		}
		curNode.NextNode[index].IndexArray = append(curNode.NextNode[index].IndexArray, pos)
		curNode = curNode.NextNode[index]
	}
}

func (node *Trie) traversalTrie() {
	if node == nil {
		logger.Println("trie is nil")
		return
	}

	for i := 0; i < _TRIE_CHILD_LEN; i++ {
		if node.NextNode[i] != nil {
			node.NextNode[i].traversalTrie()
			/*logger.Println(node.NextNode[i].Key)*/
			logger.Println(node.NextNode[i].IndexArray)
		}
	}
}

func (root *Trie) searchTrie(keys string) []int32 {
	if root == nil {
		return nil
	}
	if l := len(keys); l == 0 {
		return nil
	}

	curNode := root
	low := strings.ToLower(keys)
	for i, _ := range low {
		if i >= 'a' && i <= 'z' {
			index := low[i] - 'a'
			if curNode.NextNode[index] == nil {
				return nil
			}
			curNode = curNode.NextNode[index]
		} else {
			return nil
		}
	}

	retArray := curNode.IndexArray
	logger.Println("ret array:", retArray)
	return retArray
}
