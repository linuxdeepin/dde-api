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
	"github.com/disintegration/imaging"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path"
	"pkg.deepin.io/lib/graphic"
	dutils "pkg.deepin.io/lib/utils"
	"runtime/debug"
)

const (
	destDir = "/var/cache/image-blur/"

	defaultSigma float64 = 20.0
)

var (
	force = kingpin.Flag("force", "Force to blur image").Short('f').Bool()
	sigma = kingpin.Flag("sigma", "The blur sigma").Short('s').Float64()
	src   = kingpin.Arg("src", "The src file, may be directory").String()
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

	if *sigma == 0 {
		*sigma = defaultSigma
	}
	for _, image := range images {
		err := blurImage(image, *sigma, getDestPath(image))
		if err != nil {
			fmt.Printf("Blur '%s' failed: %v\n", image, err)
		}
	}
}

func blurImage(file string, sigma float64, dest string) error {
	img, err := imaging.Open(file)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Dir(dest), 0755)
	if err != nil {
		return err
	}

	defer debug.FreeOSMemory()

	nrgb := imaging.Blur(img, sigma)
	return imaging.Save(nrgb, dest)
}

func getDestPath(src string) string {
	id, _ := dutils.SumStrMd5(src)
	return destDir + id + path.Ext(src)
}
