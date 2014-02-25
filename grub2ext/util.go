/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
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
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"time"
)

func copyFile(src, dest string) (written int64, err error) {
	if dest == src {
		return -1, errors.New("source and destination are same file")
	}

	sf, err := os.Open(src)
	if err != nil {
		return
	}
	defer sf.Close()
	df, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return
	}
	defer df.Close()
	return io.Copy(df, sf)
}

// TODO move to go-dlib/os, dde-api/os
func execAndWait(timeout int, name string, arg ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(name, arg...)
	var buf_stdout, buf_stderr bytes.Buffer
	cmd.Stdout = &buf_stdout
	cmd.Stderr = &buf_stderr
	err = cmd.Start()
	if err != nil {
		_LOGGER.Error(err.Error())
		return
	}

	// wait for process finished
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err = cmd.Process.Kill(); err != nil {
			_LOGGER.Error(err.Error())
			return
		}
		<-done
		_LOGGER.Info("time out and process was killed")
	case err = <-done:
		stdout = buf_stdout.String()
		stderr = buf_stderr.String()
		if err != nil {
			_LOGGER.Error("process done with error = %v", err)
			return
		}
	}
	return
}

func getImgClipSizeByResolution(screenWidth, screenHeight uint16, imgWidth, imgHeight int32) (w int32, h int32) {
	if imgWidth >= int32(screenWidth) && imgHeight >= int32(screenHeight) {
		w = int32(screenWidth)
		h = int32(screenHeight)
	} else {
		scale := float32(screenWidth) / float32(screenHeight)
		w = imgWidth
		h = int32(float32(w) / scale)
		if h > imgHeight {
			h = imgHeight
			w = int32(float32(h) * scale)
		}
	}
	return
}
