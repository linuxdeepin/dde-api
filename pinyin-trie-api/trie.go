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
	"fmt"
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

func GetNode(ch byte) *Trie {
	node := new(Trie)
	node.Key = ch
	return node
}

func NewTrie() *Trie {
	root := GetNode(' ')
	return root
}

func (root *Trie) InsertTrieInfo(values []*TrieInfo) {
	for i, v := range values {
		root.InsertStringArray(v.Pinyins, int32(i))
	}
}

func (root *Trie) InsertStringArray(strs []string, pos int32) {
	if strs == nil {
		return
	}

	for _, v := range strs {
		root.InsertString(v, pos)
	}
}

func (root *Trie) InsertString(str string, pos int32) {
	if l := len(str); l == 0 {
		return
	}
	low := strings.ToLower(str)
	curNode := root

	fmt.Println("keys: ", low)
	for i, _ := range str {
		index := low[i] - 'a'
		if curNode.NextNode[index] == nil {
			curNode.NextNode[index] = GetNode(low[i])
		}
		curNode.NextNode[index].IndexArray = append(curNode.NextNode[index].IndexArray, pos)
		curNode = curNode.NextNode[index]
	}
}

func (node *Trie) TraversalTrie() {
	if node == nil {
		fmt.Println("trie is nil")
		return
	}

	for i := 0; i < _TRIE_CHILD_LEN; i++ {
		if node.NextNode[i] != nil {
			node.NextNode[i].TraversalTrie()
			/*fmt.Println(node.NextNode[i].Key)*/
			fmt.Println(node.NextNode[i].IndexArray)
		}
	}
}

func (root *Trie) SearchTrie(keys string) []int32 {
	if root == nil {
		return nil
	}
	if l := len(keys); l == 0 {
		return nil
	}

	curNode := root
	low := strings.ToLower(keys)
	for i, _ := range low {
		index := low[i] - 'a'
		if curNode.NextNode[index] == nil {
			return nil
		}
		curNode = curNode.NextNode[index]
	}

	retArray := curNode.IndexArray
	fmt.Println("ret array:", retArray)
	return retArray
}
