// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"log"

	"github.com/linuxdeepin/dde-api/blurimage"
)

var sigma = flag.Float64("sigma", 20.0, "control the strength of the blurring effect")

func main() {
	flag.Parse()
	args := flag.Args()
	input := args[0]
	output := args[1]

	err := blurimage.BlurImage(input, *sigma, output)
	if err != nil {
		log.Fatal(err)
	}
}
