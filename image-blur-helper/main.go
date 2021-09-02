/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"pkg.deepin.io/dde/api/blurimage"
	"pkg.deepin.io/lib/graphic"
	dutils "pkg.deepin.io/lib/utils"

	"gopkg.in/alecthomas/kingpin.v2"
)

const defaultOutDir = "/var/cache/image-blur/"

var (
	force      = kingpin.Flag("force", "Force to blur image").Short('f').Default("false").Bool()
	radius     = kingpin.Flag("radius", "The radius [3 - 49], must odd number(default 39)").Short('r').Default("39").Uint8()
	rounds     = kingpin.Flag("rounds", "The number of round(default 14).").Short('p').Default("14").Uint64()
	sigma      = kingpin.Flag("sigma", "The blur sigma(default 20.0).").Short('S').Default("20.0").Float64()
	saturation = kingpin.Flag("saturation", "Multiple current saturation(default 1.5)").Short('s').Default("1.5").Float64()
	lightness  = kingpin.Flag("lightness", "Multiple current lightness(HSL)(default 0.9)").Short('l').Default("0.9").Float64()
	src        = kingpin.Arg("src", "The src file, may be directory").String()
	outDir     = kingpin.Arg("outDir", "The out directory").Default(defaultOutDir).String()
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

	syscall.Setpriority(syscall.PRIO_PROCESS, 0, 18)

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

	for _, image := range images {
		dest := getDestPath(image)
		if !*force && dutils.IsFileExist(dest) {
			continue
		}

		args := []string{
			"-l", strconv.FormatFloat(*lightness, 'f', -1, 64),
			"-s", strconv.FormatFloat(*saturation, 'f', -1, 64),
			"-r", strconv.FormatUint(uint64(*radius), 10),
			"-p", strconv.FormatUint(uint64(*rounds), 10),
			image,
			"-o", dest,
		}
		out, err := exec.Command("blur_image", args...).CombinedOutput()
		if err != nil {
			fmt.Printf("Blur '%s' via 'blur_image' failed: %v, %v, try again...\n", image, string(out), err)
		}
		// fallback
		if !dutils.IsFileExist(dest) {
			err = blurimage.BlurImage(image, *sigma, dest)
			if err != nil {
				fmt.Printf("Blur '%s' via 'blurimage' failed: %s\n", image, err)
			}
		}
	}
}

func getDestPath(src string) string {
	id, _ := dutils.SumStrMd5(src)
	return filepath.Join(*outDir, id+filepath.Ext(src))
}
