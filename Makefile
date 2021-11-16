PREFIX = /usr
GOBUILD_DIR = gobuild
GOPKG_PREFIX = pkg.deepin.io/dde/api
GOSITE_DIR = ${PREFIX}/share/gocode
libdir = /lib
SYSTEMD_LIB_DIR = ${libdir}
SYSTEMD_SERVICE_DIR = ${SYSTEMD_LIB_DIR}/systemd/system/
GOBUILD = env GOPATH="${CURDIR}/${GOBUILD_DIR}:${GOPATH}" go build
export GO111MODULE=off

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
	userenv \
	inhibit_hint \
    powersupply	\
    polkit

BINARIES =  \
    device \
    graphic \
    locale-helper \
    thumbnailer \
    hans2pinyin \
    cursor-helper \
    gtk-thumbnailer \
    sound-theme-player \
    deepin-shutdown-sound \
    dde-open \
	adjust-grub-theme \
    image-blur \
    image-blur-helper
    #lunar-calendar \

all: build-binary build-dev ts-to-policy

prepare:
	@mkdir -p ${GOBUILD_DIR}/src/$(dir ${GOPKG_PREFIX});
	@ln -snf ../../../.. ${GOBUILD_DIR}/src/${GOPKG_PREFIX};

ts:
	deepin-policy-ts-convert policy2ts misc/polkit-action/com.deepin.api.locale-helper.policy.in misc/ts/com.deepin.api.locale-helper.policy
	deepin-policy-ts-convert policy2ts misc/polkit-action/com.deepin.api.device.unblock-bluetooth-devices.policy.in misc/ts/com.deepin.api.device.unblock-bluetooth-devices.policy

ts-to-policy:
	deepin-policy-ts-convert ts2policy misc/polkit-action/com.deepin.api.locale-helper.policy.in misc/ts/com.deepin.api.locale-helper.policy misc/polkit-action/com.deepin.api.locale-helper.policy
	deepin-policy-ts-convert ts2policy misc/polkit-action/com.deepin.api.device.unblock-bluetooth-devices.policy.in misc/ts/com.deepin.api.device.unblock-bluetooth-devices.policy misc/polkit-action/com.deepin.api.device.unblock-bluetooth-devices.policy

out/bin/%:
	${GOBUILD} -o $@ ${GOBUILD_OPTIONS} ${GOPKG_PREFIX}/${@F}

# Install go packages
build-dep:
	go get github.com/disintegration/imaging
	go get gopkg.in/check.v1
	go get github.com/linuxdeepin/go-x11-client

build-binary: prepare $(addprefix out/bin/, ${BINARIES})

install-binary:
	mkdir -pv ${DESTDIR}${PREFIX}${libdir}/deepin-api
	cp out/bin/* ${DESTDIR}${PREFIX}${libdir}/deepin-api/

	mkdir -pv ${DESTDIR}${PREFIX}/bin
	cp out/bin/dde-open ${DESTDIR}${PREFIX}/bin
	rm ${DESTDIR}${PREFIX}${libdir}/deepin-api/dde-open

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1/system.d
	cp misc/conf/*.conf ${DESTDIR}${PREFIX}/share/dbus-1/system.d/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1/services
	cp -v misc/services/*.service ${DESTDIR}${PREFIX}/share/dbus-1/services/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1/system-services
	cp -v misc/system-services/*.service ${DESTDIR}${PREFIX}/share/dbus-1/system-services/

	mkdir -pv ${DESTDIR}${PREFIX}/share/polkit-1/actions
	cp misc/polkit-action/*.policy ${DESTDIR}${PREFIX}/share/polkit-1/actions/

	mkdir -pv ${DESTDIR}/var/lib/polkit-1/localauthority/10-vendor.d
	cp misc/polkit-localauthority/*.pkla ${DESTDIR}/var/lib/polkit-1/localauthority/10-vendor.d/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dde-api
	cp -R misc/data ${DESTDIR}${PREFIX}/share/dde-api

	mkdir -pv ${DESTDIR}${SYSTEMD_SERVICE_DIR}
	cp -R misc/systemd/system/*.service ${DESTDIR}${SYSTEMD_SERVICE_DIR}

	mkdir -pv ${DESTDIR}${PREFIX}/share/icons/hicolor
	cp -R misc/icons/* ${DESTDIR}${PREFIX}/share/icons/hicolor

build-dev: prepare
	${GOBUILD} $(addprefix ${GOPKG_PREFIX}/, ${LIBRARIES})

install/lib/%:
	mkdir -pv ${DESTDIR}${GOSITE_DIR}/src/${GOPKG_PREFIX}
	cp -R ${CURDIR}/${GOBUILD_DIR}/src/${GOPKG_PREFIX}/${@F} ${DESTDIR}${GOSITE_DIR}/src/${GOPKG_PREFIX}

install-dev: ${addprefix install/lib/, ${LIBRARIES}}

install: install-binary install-dev

clean:
	rm -rf out/bin gobuild out

rebuild: clean build

check_code_quality: prepare
	env GOPATH="${CURDIR}/${GOBUILD_DIR}:${GOPATH}" go vet ./...

test: prepare
	env GOPATH="${CURDIR}/${GOBUILD_DIR}:${GOPATH}" go test -v ./...

print_gopath: prepare
	GOPATH="${CURDIR}/${GOBUILD_DIR}:${GOPATH}"

test-coverage: prepare
	env GOPATH="${CURDIR}/${GOBUILD_DIR}:${GOPATH}" go test -cover -v ./... | awk '$$1 ~ "^(ok|\\?)" {print $$2","$$5}' | sed "s:${CURDIR}::g" | sed 's/files\]/0\.0%/g' > coverage.csv
