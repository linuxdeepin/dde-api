package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"pkg.deepin.io/gir/gio-2.0"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.sessionmanager"
	"pkg.deepin.io/lib/dbus1"
)

var versionFlag bool

func init() {
	log.SetFlags(log.Lshortfile)
	flag.BoolVar(&versionFlag, "version", false, "show version")
}

func main() {
	flag.Parse()
	if versionFlag {
		fmt.Println("0.0.1")
		return
	}

	rawUrl := flag.Arg(0)
	if rawUrl == "" {
		log.Fatal("rawUrl empty")
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		log.Fatal(err)
	}

	switch u.Scheme {
	case "", "file":
		openFile(u.Path)

	default:
		openScheme(u.Scheme, rawUrl)
	}
}

func launchApp(desktopFile, filename string) {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	startManager := sessionmanager.NewStartManager(sessionBus)
	err = startManager.LaunchApp(dbus.FlagNoAutoStart, desktopFile, 0,
		[]string{filename})
	if err != nil {
		log.Fatal(err)
	}
}

func openFile(filename string) {
	log.Printf("openFile: %q\n", filename)
	filename, err := filepath.Abs(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	file := gio.FileNewForPath(filename)
	defer file.Unref()

	fileInfo, err := file.QueryInfo(gio.FileAttributeStandardContentType, gio.FileQueryInfoFlagsNone, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer fileInfo.Unref()
	contentType := fileInfo.GetAttributeString(gio.FileAttributeStandardContentType)
	if contentType == "" {
		log.Fatal("failed to get file content type")
	}

	appInfo := gio.AppInfoGetDefaultForType(contentType, false)
	if appInfo == nil {
		log.Fatal("failed to get appInfo")
	}
	defer appInfo.Unref()

	dAppInfo := gio.ToDesktopAppInfo(appInfo)
	desktopFile := dAppInfo.GetFilename()
	log.Println("desktop file:", desktopFile)
	launchApp(desktopFile, filename)
}

func openScheme(scheme, url string) {
	log.Printf("openScheme: %q, %q\n", scheme, url)
	appInfo := gio.AppInfoGetDefaultForUriScheme(scheme)

	if appInfo == nil {
		log.Fatal("failed to get appInfo")
	}

	dAppInfo := gio.ToDesktopAppInfo(appInfo)
	desktopFile := dAppInfo.GetFilename()
	log.Println("desktop file:", desktopFile)
	launchApp(desktopFile, url)
	appInfo.Unref()
}
