## DDE API

DDE API项目提供了一些 dbus 接口，用于屏幕区域检测、缩略图生成、声音播放等。

## 依赖

### 编译依赖

* [go-lib](https://github.com/linuxdeepin/go-lib)
* [dbus-factory](https://github.com/linuxdeepin/dbus-factory)
* gtk+-3.0
* librsvg-2.0
* libcanberra
* poppler-glib
* libsqlite3

### 运行依赖

* xcur2png
* deepin-desktop-base
* libsqlite3

## 安装

dde-api需要预安装以下包

```shell
$ go get gopkg.in/alecthomas/kingpin.v2
$ go get github.com/disintegration/imaging
$ go get github.com/linuxdeepin/go-x11-client
$ go get -u -v github.com/mattn/go-sqlite3
```

构建:
```
$ make GOPATH=/usr/share/gocode
```

通过gccgo构建
```
$ make GOPATH=/usr/share/gocode USE_GCCGO=1
```

安装:
```
sudo make install
```

## 获得帮助

如果您遇到任何其他问题，您可能还会发现这些渠道很有用：

* [Gitter](https://gitter.im/orgs/linuxdeepin/rooms)
* [IRC channel](https://webchat.freenode.net/?channels=deepin)
* [Forum](https://bbs.deepin.org/)
* [WiKi](https://wiki.deepin.org/)

## 贡献指南

我们鼓励您报告问题并做出更改。

* [Contribution guide for developers](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers-en). (English)
* [开发者代码贡献指南](https://github.com/linuxdeepin/developer-center/wiki/Contribution-Guidelines-for-Developers) (中文)

## 开源协议

dde-api项目在LGPL-3.0-or-later开源协议下发布。
