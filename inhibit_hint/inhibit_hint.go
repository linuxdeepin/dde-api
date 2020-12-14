package inhibit_hint

import (
	"errors"
	"sync"
	"unicode/utf8"

	"github.com/godbus/dbus"
	"github.com/gosexy/gettext"
	"pkg.deepin.io/lib/dbusutil"
)

//go:generate dbusutil-gen em -type Object

const (
	dbusPath      = "/com/deepin/InhibitHint"
	dbusInterface = "com.deepin.InhibitHint"
)

type Object struct {
	getMu  sync.Mutex
	domain string
	name   interface{}
	icon   interface{}
}

func New(domain string) *Object {
	return &Object{
		domain: domain,
	}
}

func (o *Object) GetInterfaceName() string {
	return dbusInterface
}

func (o *Object) SetName(name string) {
	o.name = name
}

func (o *Object) SetNameFunc(fn GetFunc) {
	o.name = fn
}

type GetFunc func(why string) string

func get(v interface{}, why string) string {
	switch vv := v.(type) {
	case string:
		return vv
	case GetFunc:
		return vv(why)
	default:
		return ""
	}
}

func (o *Object) getName(why string) string {
	return get(o.name, why)
}

func (o *Object) SetIcon(icon string) {
	o.icon = icon
}

func (o *Object) SetIconFunc(fn GetFunc) {
	o.icon = fn
}

func (o *Object) getIcon(why string) string {
	return get(o.icon, why)
}

var errInvalidUTF8 = errors.New("invalid UTF-8 string")

type HintInfo struct {
	Name string
	Icon string
	Why  string
}

func (o *Object) Get(locale string, why string) (hint *HintInfo, busErr *dbus.Error) {
	o.getMu.Lock()
	defer o.getMu.Unlock()

	gettext.SetLocale(gettext.LC_ALL, locale)
	why1 := gettext.DGettext(o.domain, why)
	if !utf8.ValidString(why1) {
		return nil, dbusutil.ToError(errInvalidUTF8)
	}

	name := gettext.DGettext(o.domain, o.getName(why))
	if !utf8.ValidString(name) {
		return nil, dbusutil.ToError(errInvalidUTF8)
	}

	icon := o.getIcon(why)

	return &HintInfo{
		Name: name,
		Icon: icon,
		Why:  why1,
	}, nil
}

func (o *Object) Export(service *dbusutil.Service) error {
	return service.Export(dbusPath, o)
}

func (o *Object) StopExport(service *dbusutil.Service) error {
	return service.StopExport(o)
}
