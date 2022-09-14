// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"

	"github.com/linuxdeepin/dde-api/thumbnails"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	SizeTypeNormal int = 128
	SizeTypeLarge  int = 256
)

func main() {
	var (
		src  = kingpin.Flag("src", "Source file").String()
		size = kingpin.Flag("size", "Thumbnail size").Int()
	)
	kingpin.Parse()

	if len(*src) == 0 {
		fmt.Println("Please input source file")
		return
	}

	if *size < 0 {
		fmt.Println("Invalid size:", *size)
		return
	}

	if *size == 0 {
		dest, err := thumbnails.GenThumbnail(*src, SizeTypeNormal)
		if err != nil {
			fmt.Printf("Gen '%s' thumbnail in size '%v' failed: %v\n", *src, SizeTypeNormal, err)
			return
		}
		fmt.Printf("Thumbnail[%v]: %v\n", SizeTypeNormal, dest)

		dest, err = thumbnails.GenThumbnail(*src, SizeTypeLarge)
		if err != nil {
			fmt.Printf("Gen '%s' thumbnail in size '%v' failed: %v\n", *src, SizeTypeLarge, err)
			return
		}
		fmt.Printf("Thumbnail[%v]: %v\n", SizeTypeLarge, dest)
		return
	}

	dest, err := thumbnails.GenThumbnail(*src, *size)
	if err != nil {
		fmt.Printf("Gen '%s' thumbnail in size '%v' failed: %v\n", *src, *size, err)
		return
	}
	fmt.Printf("Thumbnail[%v]: %v\n", *size, dest)
}
