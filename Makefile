PREFIX = /usr

BINARIES =  \
    graphic \
    greeter-utils \
    lunar-calendar \
    mousearea \
    set-date-time \
    sound

all: build

out/bin/%:
	(cd ${@F}; go build -o ../$@)

build: $(addprefix out/bin/, ${BINARIES})

install: build
	mkdir -pv ${DESTDIR}${PREFIX}/lib/deepin-api
	cp out/bin/* ${DESTDIR}${PREFIX}/lib/deepin-api/

	mkdir -pv ${DESTDIR}/etc/dbus-1/system.d
	cp misc/conf/*.conf ${DESTDIR}/etc/dbus-1/system.d/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1/services
	cp -v misc/services/*.service ${DESTDIR}${PREFIX}/share/dbus-1/services/

	mkdir -pv ${DESTDIR}${PREFIX}/share/dbus-1/system-services
	cp -v misc/system-services/*.service ${DESTDIR}${PREFIX}/share/dbus-1/system-services/

clean:
	rm -rf out/bin

rebuild: clean build
