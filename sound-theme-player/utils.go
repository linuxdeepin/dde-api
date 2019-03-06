package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/keyfile"
)

func getLightDMAutoLoginUser() (string, error) {
	kf := keyfile.NewKeyFile()
	err := kf.LoadFromFile("/etc/lightdm/lightdm.conf")
	if err != nil {
		return "", err
	}

	v, err := kf.GetString("Seat:*", "autologin-user")
	return v, err
}

func getAsoundStateFile(uid int) string {
	return filepath.Join(homeDir, fmt.Sprintf("asound-state-%d.gz", uid))
}

func runAlsaCtlStore(uid int) error {
	stateFilename := getAsoundStateFile(uid)
	logger.Debug("store ALSA state to file:", stateFilename)
	fh, err := os.Create(stateFilename)
	if err != nil {
		return err
	}
	defer fh.Close()

	bufWriter := bufio.NewWriter(fh)
	gzipWriter := gzip.NewWriter(bufWriter)

	cmd := exec.Command(alsaCtlBin, "-f", "-", "store")
	cmd.Stdout = gzipWriter
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err = cmd.Run()
	if err != nil {
		logger.Warningf("alsactl std err: %s", errBuf.String())
		return err
	}

	err = gzipWriter.Close()
	if err != nil {
		return err
	}

	err = bufWriter.Flush()
	if err != nil {
		return err
	}

	return nil
}

func runALSARestore(uid int) error {
	stateFilename := getAsoundStateFile(uid)
	logger.Debug("restore ALSA state from file:", stateFilename)
	fh, err := os.Open(stateFilename)
	if err != nil {
		return err
	}
	bufReader := bufio.NewReader(fh)
	gzipReader, err := gzip.NewReader(bufReader)
	if err != nil {
		return err
	}

	cmd := exec.Command(alsaCtlBin, "-f", "-", "restore")
	cmd.Stdin = gzipReader
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err = cmd.Run()
	if err != nil {
		logger.Warningf("alsactl std err: %s", errBuf.Bytes())
		return err
	}

	err = gzipReader.Close()
	if err != nil {
		return err
	}

	return nil
}

func getLastUser() (int, error) {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return -1, err
	}
	lockServiceObj := sysBus.Object("com.deepin.dde.LockService",
		"/com/deepin/dde/LockService")
	var userJson string
	err = lockServiceObj.Call("com.deepin.dde.LockService.CurrentUser", 0).Store(&userJson)
	if err != nil {
		return -1, err
	}

	var v struct {
		Uid int
	}
	err = json.Unmarshal([]byte(userJson), &v)
	if err != nil {
		return -1, err
	}
	return v.Uid, nil
}
