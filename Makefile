PREFIX = /usr
GOBUILD_DIR = gobuild
GOPKG_PREFIX = pkg.deepin.io/dde/api
GOSITE_DIR = ${PREFIX}/share/gocode
libdir = /lib
SYSTEMD_LIB_DIR = ${libdir}
SYSTEMD_SERVICE_DIR = ${SYSTEMD_LIB_DIR}/systemd/system/

ifndef USE_GCCGO
    GOBUILD = go build
else
    LDFLAGS = $(shell pkg-config --libs gio-2.0 gtk+-3.0 gdk-pixbuf-xlib-2.0 x11 xi xfixes xcursor libcanberra cairo-ft poppler-glib librsvg-2.0 alsa libpulse-simple)
    GOBUILD = go build -compiler gccgo -gccgoflags "${LDFLAGS}"
endif

LIBRARIES = \
    thumbnails \
    themes \
    theme_thumb\
    dxinput \
    drandr \
    soundutils \
    lang_info \
    i18n_dependent \
    session \
    language_support \
    powersupply

BINARIES =  \
    device \
    graphic \
    locale-helper \
    lunar-calendar \
    x-event-monitor \
    thumbnailer \
    hans2pinyin \
    cursor-helper \
    gtk-thumbnailer \
    sound-theme-player \
    deepin-shutdown-sound \
    image-blur \
    image-blur-helper

all: build-binary build-dev ts-to-policy

prepare:
	@if [ ! -d ${GOBUILD_DIR}/src/${GOPKG_PREFIX} ]; then \
		mkdir -p ${GOBUILD_DIR}/src/$(dir ${GOPKG_PREFIX}); \
		ln -sf ../../../.. ${GOBUILD_DIR}/src/${GOPKG_PREFIX}; \
	fi

ts:
	deepin-policy-ts-convert policy2ts misc/polkit-action/com.deepin.api.locale-helper.policy.in misc/ts/com.deepin.api.locale-helper.policy

ts-to-policy:
	deepin-policy-ts-convert ts2policy misc/polkit-action/com.deepin.api.locale-helper.policy.in misc/ts/com.deepin.api.locale-helper.policy misc/polkit-action/com.deepin.api.locale-helper.policy

out/bin/%:
	env GOPATH="${CURDIR}/${GOBUILD_DIR}:${GOPATH}" ${GOBUILD} -o $@  ${GOPKG_PREFIX}/${@F}

# Install go packages
build-dep:
	go get github.com/disintegration/imaging
	go get github.com/BurntSushi/xgb
	go get github.com/BurntSushi/xgbutil
	go get gopkg.in/check.v1

build-binary: prepare $(addprefix out/bin/, ${BINARIES})

install-binary:
	mkdir -pv ${DESTDIR}${PREFIX}${libdir}/deepin-api
	cp out/bin/* ${DESTDIR}${PREFIX}${libdir}/deepin-api/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1/system.d
	cp misc/conf/*.conf ${DESTDIR}${PREFIX}/share/dbus-1/system.d/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1/services
	cp -v misc/services/*.service ${DESTDIR}${PREFIX}/share/dbus-1/services/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1/system-services
	cp -v misc/system-services/*.service ${DESTDIR}${PREFIX}/share/dbus-1/system-services/

	mkdir -pv ${DESTDIR}${PREFIX}/share/polkit-1/actions
	cp misc/polkit-action/*.policy ${DESTDIR}${PREFIX}/share/polkit-1/actions/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dde-api
	cp -R misc/data ${DESTDIR}${PREFIX}/share/dde-api

	mkdir -pv ${DESTDIR}${SYSTEMD_SERVICE_DIR}
	cp -R misc/systemd/system/*.service ${DESTDIR}${SYSTEMD_SERVICE_DIR}

	mkdir -pv ${DESTDIR}${PREFIX}/share/icons/hicolor
	cp -R misc/icons/* ${DESTDIR}${PREFIX}/share/icons/hicolor

build-dev: prepare
	env GOPATH="${GOPATH}:${CURDIR}/${GOBUILD_DIR}" ${GOBUILD} $(addprefix ${GOPKG_PREFIX}/, ${LIBRARIES})

install/lib/%:
	mkdir -pv ${DESTDIR}${GOSITE_DIR}/src/${GOPKG_PREFIX}
	cp -R ${CURDIR}/${GOBUILD_DIR}/src/${GOPKG_PREFIX}/${@F} ${DESTDIR}${GOSITE_DIR}/src/${GOPKG_PREFIX}

install-dev: ${addprefix install/lib/, ${LIBRARIES}}

install: install-binary install-dev

clean:
	rm -rf out/bin gobuild out

rebuild: clean build
