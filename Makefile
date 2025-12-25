PREFIX = /usr
GOBUILD_DIR = gobuild
GOPKG_PREFIX = github.com/linuxdeepin/dde-api
GOSITE_DIR = ${PREFIX}/share/gocode
libdir = /lib
SYSTEMD_LIB_DIR = ${libdir}
SYSTEMD_SERVICE_DIR = ${SYSTEMD_LIB_DIR}/systemd/system/
GOBUILD = env GOPATH="${CURDIR}/${GOBUILD_DIR}:${GOPATH}" go build
INSTALL_LOCALE_HELPER ?= 0

TESTS = \
	${GOPKG_PREFIX}/adjust-grub-theme \
	${GOPKG_PREFIX}/blurimage \
	${GOPKG_PREFIX}/dde-open \
	${GOPKG_PREFIX}/deepin-shutdown-sound \
	${GOPKG_PREFIX}/device \
	${GOPKG_PREFIX}/drandr \
	${GOPKG_PREFIX}/dxinput \
	${GOPKG_PREFIX}/dxinput/common \
	${GOPKG_PREFIX}/dxinput/kwayland \
	${GOPKG_PREFIX}/dxinput/utils \
	${GOPKG_PREFIX}/graphic \
	${GOPKG_PREFIX}/grub_theme/font \
	${GOPKG_PREFIX}/grub_theme/themetxt \
	${GOPKG_PREFIX}/gtk-thumbnailer \
	${GOPKG_PREFIX}/hans2pinyin \
	${GOPKG_PREFIX}/i18n_dependent \
	${GOPKG_PREFIX}/image-blur \
	${GOPKG_PREFIX}/image-blur-helper \
	${GOPKG_PREFIX}/inhibit_hint \
	${GOPKG_PREFIX}/lang_info \
	${GOPKG_PREFIX}/language_support \
	${GOPKG_PREFIX}/locale-helper \
	${GOPKG_PREFIX}/polkit \
	${GOPKG_PREFIX}/powersupply \
	${GOPKG_PREFIX}/powersupply/battery \
	${GOPKG_PREFIX}/session \
	${GOPKG_PREFIX}/sound-theme-player \
	${GOPKG_PREFIX}/soundutils \
	${GOPKG_PREFIX}/userenv \
	${GOPKG_PREFIX}/validator

LIBRARIES = \
    dxinput \
    drandr \
    soundutils \
    lang_info \
    i18n_dependent \
    session \
    language_support \
    userenv \
    inhibit_hint \
    powersupply \
    polkit \

ININSTALLS = \
    ${LIBRARIES} \
    go.sum \
    go.mod

BINARIES =  \
    device \
    graphic \
    locale-helper \
    hans2pinyin \
    gtk-thumbnailer \
    sound-theme-player \
    deepin-shutdown-sound \
    dde-open \
    adjust-grub-theme \
    image-blur \
    image-blur-helper

all: build-binary build-dev ts-to-policy

prepare:
	@mkdir -p ${GOBUILD_DIR}/src/$(dir ${GOPKG_PREFIX});
	@ln -snf ../../../.. ${GOBUILD_DIR}/src/${GOPKG_PREFIX};

ts:
	deepin-policy-ts-convert policy2ts misc/polkit-action/org.deepin.dde.locale-helper.policy.in misc/ts/org.deepin.dde.locale-helper.policy
	deepin-policy-ts-convert policy2ts misc/polkit-action/org.deepin.dde.device.unblock-bluetooth-devices.policy.in misc/ts/org.deepin.dde.device.unblock-bluetooth-devices.policy

ts-to-policy:
	deepin-policy-ts-convert ts2policy misc/polkit-action/org.deepin.dde.locale-helper.policy.in misc/ts/org.deepin.dde.locale-helper.policy misc/polkit-action/org.deepin.dde.locale-helper.policy
	deepin-policy-ts-convert ts2policy misc/polkit-action/org.deepin.dde.device.unblock-bluetooth-devices.policy.in misc/ts/org.deepin.dde.device.unblock-bluetooth-devices.policy misc/polkit-action/org.deepin.dde.device.unblock-bluetooth-devices.policy

out/bin/%:
	${GOBUILD} -o $@ ${GOBUILD_OPTIONS} ${GOPKG_PREFIX}/${@F}

build-binary: prepare $(addprefix out/bin/, ${BINARIES})

install-binary:
	mkdir -pv ${DESTDIR}${PREFIX}${libdir}/deepin-api
	cp out/bin/* ${DESTDIR}${PREFIX}${libdir}/deepin-api/
	cp misc/scripts/* ${DESTDIR}${PREFIX}${libdir}/deepin-api/

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

	mkdir -pv ${DESTDIR}/var/lib/polkit-1/rules.d
	cp misc/polkit-rules/*.rules ${DESTDIR}/var/lib/polkit-1/rules.d/
	
	mkdir -pv ${DESTDIR}${PREFIX}/share/dde-api
	cp -R misc/data ${DESTDIR}${PREFIX}/share/dde-api

	mkdir -pv ${DESTDIR}${SYSTEMD_SERVICE_DIR}
	cp -R misc/systemd/system/*.service ${DESTDIR}${SYSTEMD_SERVICE_DIR}
	# 默认不安装 deepin-locale-helper.service，只有显式开启时才保留
ifneq ($(INSTALL_LOCALE_HELPER), 1)
	rm -f ${DESTDIR}${SYSTEMD_SERVICE_DIR}/deepin-locale-helper.service; 
	rm -f ${DESTDIR}${PREFIX}/share/dbus-1/system-services/org.deepin.dde.LocaleHelper1.service; 
	rm -f ${DESTDIR}${PREFIX}/share/polkit-1/actions/org.deepin.dde.locale-helper.policy; 
	rm -f ${DESTDIR}${PREFIX}/share/dbus-1/system.d/org.deepin.dde.LocaleHelper1.conf; 
endif

	mkdir -pv ${DESTDIR}${PREFIX}/share/icons/hicolor
	cp -R misc/icons/* ${DESTDIR}${PREFIX}/share/icons/hicolor

build-dev: prepare
	${GOBUILD} $(addprefix ${GOPKG_PREFIX}/, ${LIBRARIES})

install/lib/%:
	mkdir -pv ${DESTDIR}${GOSITE_DIR}/src/${GOPKG_PREFIX}
	cp -R ${CURDIR}/${GOBUILD_DIR}/src/${GOPKG_PREFIX}/${@F} ${DESTDIR}${GOSITE_DIR}/src/${GOPKG_PREFIX}

install-dev: ${addprefix install/lib/, ${ININSTALLS}}

install: install-binary install-dev

clean:
	rm -rf out/bin gobuild out

rebuild: clean build

check_code_quality: prepare
	env GOPATH="${CURDIR}/${GOBUILD_DIR}:${GOPATH}" go vet ./...

test: prepare
	env GOPATH="${CURDIR}/${GOBUILD_DIR}:${GOPATH}" go test -v ${TESTS}

print_gopath: prepare
	GOPATH="${CURDIR}/${GOBUILD_DIR}:${GOPATH}"

test-coverage: prepare
	env GOPATH="${CURDIR}/${GOBUILD_DIR}:${GOPATH}" go test -cover -v ./... | awk '$$1 ~ "^(ok|\\?)" {print $$2","$$5}' | sed "s:${CURDIR}::g" | sed 's/files\]/0\.0%/g' > coverage.csv
