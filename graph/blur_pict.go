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

// #cgo pkg-config: glib-2.0 gdk-pixbuf-2.0
// #cgo LDFLAGS: -lm
// #include <stdlib.h>
// #include "blur_pict.h"
import "C"
import "unsafe"

import (
        "crypto/md5"
        "fmt"
        "os"
        "os/user"
        "strconv"
)

const (
        _BG_BLUR_PICT_CACHE_DIR = ".gaussian-background"
)

var (
        jobInHand map[string]bool
)

func (graph *Graph) BackgroundBlurPictPath(srcPath, destPath string,
        sigma, numsteps float64) (int32, string) {
        if len(srcPath) <= 0 {
                return -1, ""
        }
        homeDir, err := getHomeDir()
        if err != nil {
                fmt.Println("get home dir failed")
                return -1, ""
        }

        srcFlag := true
        if len(destPath) <= 0 {
                destPath = GenerateDestPath(srcPath, homeDir)
        }
        if IsFileValid(srcPath, destPath) {
                return 0, destPath
        }

        if MkGaussianCacheDir() {
                go func() {
                        success := GenerateBlurPict(srcPath, destPath, sigma, numsteps)
                        if success && !srcFlag {
                                userInfo, err := user.Current()
                                if err != nil {
                                        fmt.Println("New User Info Failed:", err)
                                        panic(err)
                                }
                                uidInt, _ := strconv.ParseInt(userInfo.Uid, 10, 64)
                                gidInt, _ := strconv.ParseInt(userInfo.Gid, 10, 64)
                                f, _ := os.Open(destPath)
                                defer f.Close()
                                f.Chown(int(uidInt), int(gidInt))
                                if graph.BlurPictChanged != nil {
                                        graph.BlurPictChanged(srcPath, destPath)
                                }
                        }
                }()
        }

        return 1, srcPath
}

func MkGaussianCacheDir() bool {
        userInfo, err := user.Current()
        if err != nil {
                fmt.Println("New User Info Failed:", err)
                return false
        }

        homeDir := userInfo.HomeDir
        pictPath := homeDir + "/" + _BG_BLUR_PICT_CACHE_DIR
        err = os.MkdirAll(pictPath, os.FileMode(0755))
        if err != nil {
                fmt.Println(err)
                return false
        }
        f, err1 := os.Open(pictPath)
        defer f.Close()
        if err1 != nil {
                fmt.Println("'%s' Open Failed:", err)
                return false
        }
        uidInt, _ := strconv.ParseInt(userInfo.Uid, 10, 64)
        gidInt, _ := strconv.ParseInt(userInfo.Gid, 10, 64)
        f.Chown(int(uidInt), int(gidInt))

        return true

}

func GenerateBlurPict(srcPath, destPath string, sigma, numsteps float64) bool {
        if len(srcPath) <= 0 && len(destPath) <= 0 {
                fmt.Println("args failed")
                return false
        }

        if _, ok := jobInHand[destPath]; ok {
                fmt.Printf("'%s' has been in hand\n", destPath)
                return false
        }
        src := C.CString(srcPath)
        defer C.free(unsafe.Pointer(src))
        dest := C.CString(destPath)
        defer C.free(unsafe.Pointer(dest))

        jobInHand[destPath] = true
        is_ok := C.generate_blur_pict(src, dest, C.double(sigma), C.double(numsteps))
        delete(jobInHand, destPath)
        if is_ok == 0 {
                fmt.Println("generate gaussian picture failed")
                return false
        }

        return true
}

func GenerateDestPath(srcPath, homeDir string) string {
        if len(homeDir) <= 0 && len(srcPath) <= 0 {
                fmt.Println("args failed")
                return ""
        }

        md5Sum := md5.Sum([]byte(srcPath))
        md5Str := ""
        for _, b := range md5Sum {
                s := strconv.FormatInt(int64(b), 16)
                if len(s) == 1 {
                        md5Str += "0" + s
                } else {
                        md5Str += s
                }
        }

        destPath := homeDir + "/" + _BG_BLUR_PICT_CACHE_DIR + "/" + md5Str + ".png"
        return destPath
}

func IsFileValid(srcPath, destPath string) bool {
        if len(srcPath) <= 0 && len(destPath) <= 0 {
                fmt.Println("args failed")
                return false
        }

        if !FileIsExist(destPath) {
                return false
        }

        src := C.CString(srcPath)
        defer C.free(unsafe.Pointer(src))
        dest := C.CString(destPath)
        defer C.free(unsafe.Pointer(dest))
        if C.blur_pict_is_valid(src, dest) == 0 {
                fmt.Println("file invalid")
                return false
        }

        return true
}

func getHomeDir() (string, error) {
        userInfo, err := user.Current()
        if err != nil {
                fmt.Println(err) // TODO
                return "", err
        }

        return userInfo.HomeDir, nil
}

func FileIsExist(filename string) bool {
        if _, err := os.Stat(filename); os.IsNotExist(err) {
                fmt.Printf("File '%s' not exist:%s\n", filename, err)
                return false
        }

        return true
}
