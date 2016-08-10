/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
	"path"
	"pkg.deepin.io/lib/graphic"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	destDir = "/var/cache/image-blur/"

	// defaultSigma float64 = 20.0
	defaultRadius int8   = 9
	defaultRounds uint64 = 10
)

var (
	force = kingpin.Flag("force", "Force to blur image").Short('f').Bool()
	// sigma = kingpin.Flag("sigma", "The blur sigma").Short('s').Float64()
	radius = kingpin.Flag("radius", "The radius, range: [3 - 19], must odd number").Short('r').Int8()
	rounds = kingpin.Flag("rounds", "The number of round").Short('p').Uint64()
	src    = kingpin.Arg("src", "The src file, may be directory").String()
)

func main() {
	if len(os.Args) == 1 {
		kingpin.Usage()
		return
	}

	kingpin.Parse()
	if !dutils.IsFileExist(*src) {
		fmt.Println("Not found this file:", *src)
		return
	}

	var images []string
	if dutils.IsDir(*src) {
		tmp, err := graphic.GetImagesInDir(*src)
		if err != nil {
			fmt.Printf("Get images from dir '%s' failed: %v\n", *src, err)
			return
		}
		images = tmp
	} else {
		images = append(images, *src)
	}

	if *radius == 0 || *radius < 3 || *radius > 13 {
		*radius = defaultRadius
	}
	if *rounds == 0 {
		*rounds = defaultRounds
	}

	for _, image := range images {
		dest := getDestPath(image)
		if !*force && dutils.IsFileExist(dest) {
			continue
		}

		cmd := fmt.Sprintf("exec blur_image -b -r %v -p %v %q -o %s", *radius, *rounds, image, dest)
		out, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
		if err != nil {
			fmt.Printf("Blur '%s' failed: %v\n", image, string(out))
		}
	}
}

func getDestPath(src string) string {
	id, _ := dutils.SumStrMd5(src)
	return destDir + id + path.Ext(src)
}
