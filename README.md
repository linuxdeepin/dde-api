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

### Runtime Dependencies

* xcur2png
* deepin-desktop-base

## Installation

Install prerequisites

```shell
$ go get github.com/BurntSushi/xgbutil
$ go get gopkg.in/alecthomas/kingpin.v2
$ go get github.com/disintegration/imaging
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
* [WiKi](http://wiki.deepin.org/)

## Getting involved

We encourage you to report issues and contribute changes.

* [Contribution guide for users](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Users)
* [Contribution guide for developers](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Developers)

## License

DDE API is licensed under [GPLv3](LICENSE).
