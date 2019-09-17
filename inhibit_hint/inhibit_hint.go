package inhibit_hint

import (
	"errors"
	"sync"
	"unicode/utf8"

	"github.com/gosexy/gettext"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusPath      = "/com/deepin/inhibit/Hint"
	dbusInterface = "com.deepin.inhibit.Hint"
)

type Object struct {
	getMu   sync.Mutex
	domain  string
	name    string
	methods *struct {
		GetName func() `in:"lang" out:"name"`
		GetText func() `in:"lang,msgId" out:"msgStr"`
	}
}

func New(domain, name string) *Object {
	return &Object{
		domain: domain,
		name:   name,
	}
}

func (o *Object) GetInterfaceName() string {
	return dbusInterface
}

func (o *Object) GetName(locale string) (string, *dbus.Error) {
	return o.GetText(locale, o.name)
}

var errInvalidUTF8 = errors.New("invalid UTF-8 string")

func (o *Object) GetText(locale string, msgId string) (string, *dbus.Error) {
	o.getMu.Lock()
	defer o.getMu.Unlock()
	gettext.SetLocale(gettext.LC_ALL, locale)
	msgStr := gettext.DGettext(o.domain, msgId)
	if !utf8.ValidString(msgStr) {
		return "", dbusutil.ToError(errInvalidUTF8)
	}
	return msgStr, nil
}

func (o *Object) Export(service *dbusutil.Service) error {
	return service.Export(dbusPath, o)
}

func (o *Object) StopExport(service *dbusutil.Service) error {
	return service.StopExport(o)
}
