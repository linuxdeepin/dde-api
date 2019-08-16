## DDE API

DDE API provides some dbus interfaces that is used for screen zone detecting, thumbnail generating, sound playing, etc.

## Dependencies


### Build Dependencies

* [go-lib](https://github.com/linuxdeepin/go-lib)
* [dbus-factory](https://github.com/linuxdeepin/dbus-factory)
* gtk+-3.0
* librsvg-2.0
* libcanberra
* poppler-glib
* libsqlite3

### Runtime Dependencies

* xcur2png
* deepin-desktop-base
* libsqlite3

## Installation

Install prerequisites

```shell
$ go get gopkg.in/alecthomas/kingpin.v2
$ go get github.com/disintegration/imaging
$ go get github.com/linuxdeepin/go-x11-client
$ go get -u -v github.com/mattn/go-sqlite3
```

Build:
```
$ make GOPATH=/usr/share/gocode
```

Or, build through gccgo
```
$ make GOPATH=/usr/share/gocode USE_GCCGO=1
```

Install:
```
sudo make install
```

## Getting help

Any usage issues can ask for help via

* [Gitter](https://gitter.im/orgs/linuxdeepin/rooms)
* [IRC channel](https://webchat.freenode.net/?channels=deepin)
* [Forum](https://bbs.deepin.org/)
* [WiKi](https://wiki.deepin.org/)

## Getting involved

We encourage you to report issues and contribute changes.

* [Contribution guide for developers](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers-en). (English)
* [开发者代码贡献指南](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers) (中文)

## License

DDE API is licensed under [GPLv3](LICENSE).
