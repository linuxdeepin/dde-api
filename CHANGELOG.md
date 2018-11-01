## [3.3.0] - 2018-11-01
*   chore(adjust-grub-theme): small appearance adjustment
*   fix: update-grub not executed after adjust grub theme
*   feat(adjust-grub-theme): use pkg.deepin.io/lib/log
*   fix(adjust-grub-theme): save boot menu position rel value into theme.txt
*   auto sync po files from transifex

## [3.2.0] - 2018-10-25
*   feat: add binary adjust-grub-theme

## [3.1.30] - 2018-08-14
*   fix(`pkg_depends`): not found cangjie input method
*   chore(`pkg_depneds`): zh-hant do not install fcitx-ui-qimpanel

## [3.1.29] - 2018-08-12
*   chore: auto sync ts files from transifex
*   chore: add unblock-bluetooth-device policy to tx config

## [3.1.28] - 2018-08-07
*   fix(device): rfkillBin use absolute path
*   feat(locale-helper): more safe call locale-gen
*   feat(device): limit com.deepin.api.Device service
*   feat: do not use the root user to run sound theme player

## [3.1.27] - 2018-07-20
*   chore(debian): update require versions
*   feat(mouse): interface to config accel profile via libinput
*   auto sync po files from transifex
*   fix(drandr): rate calculation inaccuracies
*   chore(locale-helper): no use pkg.deepin.io/lib/polkit
*   chore: use go-dbus-factory
*   fix(dde-open): panic when run dde-open xxxxxxx:///xxxxxxxx
*   chore(debian): update debian control
*   chore(drandr): use go-x11-client

## [3.1.26] - 2018-06-07
*   chore: update makefile for arch `sw_64`

## [3.1.25] - 2018-05-14
*   auto sync po files from transifex

## [3.1.24] - 2018-05-14
*   feat(theme-thumb): use rsvg-convert convert svg to png
*   fix(theme-thumb): printf argument wrong

## [3.1.23] - 2018-04-18
*   feat: add dde-open
*   refactor: fix some typos

## [3.1.22] - 2018-03-19
*   chore: use lib dbusutil new api

## [3.1.21] - 2018-03-07
*   auto sync po files from transifex
*   feat(sound-theme-player): use newly lib dbusutil
*   feat(cursor-helper): use newly lib dbusutil
*   feat(hans2pinyin): use newly lib dbusutil
*   feat(lunar-calendar): use newly lib dbusutil
*   feat(locale-helper): use newly lib dbusutil
*   feat(graphic): use newly lib dbusutil
*   feat(device): use newly lib dbusutil
*   chore: update license

## [3.1.20] - 2018-01-24
*   sound-theme-player: use newly `sound_effect` lib
*   add libs for gccgo build
*   fix Adapt lintian
*   remove x-event-monitor

## [3.1.19] - 2017-12-15
*   add lib `language_support`

## [3.1.18] - 2017-12-13
*   doc: update links in README
*   fix a typo
*   rename mousearea to x-event-monitor
*   soundutils: use `sound_effect` lib

## [3.1.17] - 2017-11-6
*   `theme_thumb`: fix some cursor theme thumbnails generate failed 


## [3.1.16] - 2017-11-3
*   dxinput: disable tap middle button
*   Remove target 'build-dev' from 'install' to 'build'


## [3.1.15] - 2017-10-25
#### Added
*   Add deepin-gettext-tools to build dependencies ([d06b8d39](d06b8d39))
*   Add `theme_thumb` ([c513f59f](c513f59f))

#### Fixed
*   Fix `theme_thumb` build failed with old version gccgo ([f6881342](f6881342))
*   Fix policykit message not using user's locale ([c7a9c53a](c7a9c53a))


## [3.1.14] - 2017-10-12
#### Added
*   Add 'image-blur'
*   Add 'outDir' option in image-blur-helper
*   Add palm settings in dxinput

#### Changed
*   Update license
*   Update sound events name
