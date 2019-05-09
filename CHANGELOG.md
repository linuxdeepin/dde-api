[3.18.4] 2019-05-09
*   chore: correct a typo
*   fix(device): checkAuthorization is not secure
*   feat(cursor-helper): set com.deepin.wm cursorTheme

[3.18.3] 2019-04-10
*   chore: auto pull translation files from transifex

[3.18.2] 2019-04-01
*   feat(sound-theme-player): add retry for alsactl restore

[3.18.1] 2019-03-27
*   fix(adjust-grub-theme): terminal font has cursor image residue
*   chore(deb): override `dh_systemd_start`

[3.18.0] 2019-03-15
*   change(api): sound-theme-player add more methods
*   chore: modify GOPATH directroy order

[3.17.1] 2019-03-01
*   chore: use pkg.deepin.io/gir

[3.17.0] 2019-02-22
*   chore: `pkg_depends` remove fcitx-sogoupinyin-uk
*   chore(`language_support`): do not use lastore dbus methods
*   fix: deepin grub theme text typo
*   fix: can not play shutdown music

[3.16.0] 2019-01-25
*   fix: build miss pkg userenv
*   feat: add package userenv
*   feat: grub-theme add key e and c usage
*   auto sync po files from transifex

## [3.15.0] - 2019-01-03
*   chore: depends remove grub-common

## [3.14.0] - 2018-12-29
*   feat(soundutils): and method for get sound file

## [3.13.0] - 2018-12-25
*   chore: change background image format from png to jpeg
*   chore: compile with sw arch no longer needs to use gccgo

## [3.12.0] - 2018-12-14
*   fix(adjust-grub-theme): only 4 lines of options can be displayed when the resolution is 1280x1024

## [3.11.0] - 2018-12-07
*   auto sync po files from transifex
*   chore(adjust-grub-theme): use lib imgutil
*   feat: remove code about setting `GRUB_GFXMODE` from the debian/postinst script
*   chore: `pkg_depends` add rule for thunderbird-l10n-

## [3.10.0] - 2018-11-23
*   fix: some types
*   fix: crash on dde-open when no appinfo available
*   fix(adjust-grub-theme): set background failed

## [3.9.0] - 2018-11-15
*   fix(adjust-grub-theme): no append `GRUB_GFXMODE=1024x768`

## [3.8.0] - 2018-11-14
*   feat(adjust-grub-theme): reset `GRUB_GFXMODE` to 1024x768 only once

## [3.7.0] - 2018-11-13
*   fix: icon theme deepin-dark thumbnail not well

## [3.6.0] - 2018-11-13
*   fix: correctly parse rfkill output

## [3.5.0] - 2018-11-08
*   feat(adjust-grub-theme): add scrollbar thumb
*   fix(adjust-grub-theme): failed to display 5 items
*   chore(adjust-grub-theme): update background.origin.png
*   chore(adjust-grub-theme): update version
*   fix: update-grub not executed after adjust grub theme
*   chore(grub-theme): add os.svg
*   fix(adjust-grub-theme): set background copy file failed
*   feat(dxinput): add libinput pointer rotation supported

## [3.4.0] - 2018-11-01
*   chore: no call grub-mkconfig

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
