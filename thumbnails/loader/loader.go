// Load thumbnail handlers
package loader

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/graphic"
	dutils "pkg.deepin.io/lib/utils"
	"sync"
	"time"
)

const (
	SizeFlagLarge  int = 256
	SizeFlagNormal     = 128
	SizeFlagSmall      = 64
)

// args: src, bg, width, height, force
// rets: dest, error
type HandleType func(string, string, int, int, bool) (string, error)

type Manager struct {
	handlers map[string]HandleType
	locker   *sync.RWMutex
}

var mInitializer sync.Once

var getManager = func() func() *Manager {
	var m *Manager
	return func() *Manager {
		mInitializer.Do(func() {
			m = &Manager{
				handlers: make(map[string]HandleType),
				locker:   new(sync.RWMutex),
			}
		})
		return m
	}
}()

func Register(ty string, handler HandleType) {
	err := getManager().Add(ty, handler)
	if err != nil {
		fmt.Println("Register failed:", err)
	}
}

func IsStrInList(item string, items []string) bool {
	for _, v := range items {
		if item == v {
			return true
		}
	}
	return false
}

func GetHandler(ty string) (HandleType, error) {
	handler := getManager().Get(ty)
	if handler == nil {
		return nil, fmt.Errorf("Cann't find generator for '%s'", ty)
	}
	return handler, nil
}

func ThumbnailImage(src, dest string, width, height int) error {
	return graphic.ScaleImage(src, dest, width, height,
		graphic.FormatPng)
}

func GetThumbnailDest(uri string, width, height int) (string, error) {
	md5, ok := dutils.SumFileMd5(dutils.DecodeURI(uri))
	if !ok {
		return "", fmt.Errorf("md5sum '%s' failed", uri)
	}

	var mid string
	if width >= SizeFlagLarge || height >= SizeFlagLarge {
		mid = "thumbnails/large"
	} else if width >= SizeFlagNormal || height >= SizeFlagNormal {
		mid = "thumbnails/normal"
	} else {
		mid = "thumbnails/small"
	}
	dir := path.Join(glib.GetUserCacheDir(), mid)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return "", err
	}

	return path.Join(dir, md5+".png"), nil
}

func CompositeIcons(icons []string, bg, dest string,
	iconSize, width, height int) error {
	iconNum := len(icons)
	if iconNum == 0 {
		return fmt.Errorf("No icon files")
	}

	var err error
	if !dutils.IsFileExist(bg) {
		bg, err = GetBackground(width, height)
		if err != nil {
			return err
		}
		defer os.Remove(bg)
	}

	deltaX := (width - iconSize*iconNum) / (iconNum + 1)
	deltaY := (height - iconSize) / 2

	xPoint := deltaX
	yPoint := deltaY
	for _, icon := range icons {
		err = graphic.CompositeImage(bg, icon, dest,
			xPoint, yPoint, graphic.FormatPng)
		if err != nil {
			return err
		}
		bg = dest
		xPoint += deltaX + iconSize
	}
	return nil
}

func GetBackground(width, height int) (string, error) {
	var dest = GetTmpImage()
	err := graphic.NewImageWithColor(dest, int(width), int(height),
		//uint8(192), uint8(192), uint8(192), uint8(250),
		uint8(250), uint8(250), uint8(250), uint8(0),
		graphic.FormatPng)
	if err != nil {
		return "", err
	}

	return dest, nil
}

func GetTmpImage() string {
	var (
		seedStr = "0123456789abcdefghijklmnopqrstuvwxyz"
		ret     string
	)
	length := len(seedStr)
	for i := 0; i < 8; i++ {
		rand.Seed(time.Now().UnixNano())
		ret += string(seedStr[rand.Intn(length)])
	}
	return path.Join("/tmp", ret+".png")
}

func (m *Manager) Add(ty string, handler HandleType) error {
	v := m.Get(ty)
	if v != nil {
		return fmt.Errorf("'%s' has been registed", ty)
	}

	m.locker.Lock()
	defer m.locker.Unlock()
	m.handlers[ty] = handler
	return nil
}

func (m *Manager) Delete(ty string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	delete(m.handlers, ty)
}

func (m *Manager) Get(ty string) HandleType {
	m.locker.RLock()
	defer m.locker.RUnlock()
	handler, ok := m.handlers[ty]
	if !ok {
		return nil
	}

	return handler
}
