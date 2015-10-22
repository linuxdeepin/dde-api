PREFIX = /usr
GOPATH_DIR = gopath
GOPKG_PREFIX = pkg.deepin.io/dde/api

ifndef USE_GCCGO
    GOBUILD = go build
else
    LDFLAGS = $(shell pkg-config --libs gio-2.0 gdk-3.0 gdk-pixbuf-xlib-2.0 x11 xi libcanberra cairo-ft poppler-glib librsvg-2.0 libmetacity-private)
    GOBUILD = go build -compiler gccgo -gccgoflags "${LDFLAGS}"
endif

LIBRARIES = \
    thumbnails \
    themes

BINARIES =  \
    device \
    graphic \
    greeter-helper \
    locale-helper \
    lunar-calendar \
    mousearea \
    thumbnailer \
    mime-helper \
    sound \
    hans2pinyin \
    cursor-helper \
    gtk-thumbnailer

all: build

prepare:
	@if [ ! -d ${GOPATH_DIR}/src/${GOPKG_PREFIX} ]; then \
		mkdir -p ${GOPATH_DIR}/src/$(dir ${GOPKG_PREFIX}); \
		ln -sf ../../../.. ${GOPATH_DIR}/src/${GOPKG_PREFIX}; \
	fi

out/bin/%:
	env GOPATH="${GOPATH}:${CURDIR}/${GOPATH_DIR}" ${GOBUILD} -o $@  ${GOPKG_PREFIX}/${@F}

# Install go packages
build-dep:
	go get github.com/disintegration/imaging
	go get github.com/BurntSushi/xgb
	go get github.com/BurntSushi/xgbutil
	go get github.com/howeyc/fsnotify
	go get launchpad.net/gocheck

build: prepare $(addprefix out/bin/, ${BINARIES})

install-binary: build
	mkdir -pv ${DESTDIR}${PREFIX}/lib/deepin-api
	cp out/bin/* ${DESTDIR}${PREFIX}/lib/deepin-api/

	mkdir -pv ${DESTDIR}/etc/dbus-1/system.d
	cp misc/conf/*.conf ${DESTDIR}/etc/dbus-1/system.d/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1/services
	cp -v misc/services/*.service ${DESTDIR}${PREFIX}/share/dbus-1/services/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1/system-services
	cp -v misc/system-services/*.service ${DESTDIR}${PREFIX}/share/dbus-1/system-services/

	mkdir -pv ${DESTDIR}${PREFIX}/share
	cp -R misc/dde-api ${DESTDIR}${PREFIX}/share

build/lib/%:
	env GOPATH="${GOPATH}:${CURDIR}/${GOPATH_DIR}" ${GOBUILD} ${GOPKG_PREFIX}/${@F}

build-dev: prepare $(addprefix build/lib/, ${LIBRARIES})

install/lib/%:
	mkdir -pv ${DESTDIR}${PREFIX}/share/gocode/src/${GOPKG_PREFIX}
	cp -R ${CURDIR}/${GOPATH_DIR}/src/${GOPKG_PREFIX}/${@F} ${DESTDIR}${PREFIX}/share/gocode/src/${GOPKG_PREFIX}

install-dev: build-dev ${addprefix install/lib/, ${LIBRARIES}}

install: install-binary install-dev

clean:
	rm -rf out/bin

rebuild: clean build
