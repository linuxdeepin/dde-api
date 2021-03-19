# Run tests in check section
# disable for bootstrapping
%bcond_with check
#%global _unpackaged_files_terminate_build 0
%global with_debug 1

%if 0%{?with_debug}
%global debug_package   %{nil}
%endif

%global sname deepin-api

# out of memory on armv7hl
%ifarch %{arm}
%global _smp_mflags -j1
%endif

%global goipath  pkg.deepin.io/dde/api
%global forgeurl https://github.com/linuxdeepin/dde-api
%global tag      %{version}

Name:           dde-api
Version:        5.4.4
Release:        1
Summary:        Go-lang bingding for dde-daemon
License:        GPLv3+
URL:            https://shuttle.corp.deepin.com/cache/tasks/19177/unstable-amd64/
Source0:        %{name}-%{version}.orig.tar.xz

BuildRequires:  libcanberra-devel
BuildRequires:  deepin-gettext-tools
BuildRequires:  librsvg2-devel
BuildRequires:  sqlite-devel
BuildRequires:  compiler(go-compiler)
BuildRequires:  golang-github-linuxdeepin-go-x11-client-devel
BuildRequires:  gdk-pixbuf2-xlib-devel
BuildRequires:  kf5-kwayland-devel
BuildRequires:  poppler-glib
BuildRequires:  poppler-glib-devel
BuildRequires:  alsa-lib-devel
BuildRequires:  alsa-lib
BuildRequires:  pulseaudio-libs-devel
BuildRequires:  gocode
BuildRequires:  deepin-gir-generator
BuildRequires:  golang-github-linuxdeepin-go-dbus-factory-devel
BuildRequires:  go-lib-devel
BuildRequires:  libgudev-devel
%{?systemd_requires}
Requires:       deepin-desktop-base
Requires:       rfkill
Requires(pre):  shadow-utils

%description
%{summary}.

%package -n %{name}-devel
Summary:        %{summary}
BuildArch:      noarch

%description -n %{name}-devel
%{summary}.

This package contains library source intended for
building other packages which use import path with
%{goipath} prefix.

%prep
%forgeautosetup -p1 -n %{name}-%{version}

sed -i 's|/usr/lib|%{_libexecdir}|' misc/*services/*.service \
    misc/systemd/system/deepin-shutdown-sound.service \
    lunar-calendar/main.go \
    theme_thumb/gtk/gtk.go \
    thumbnails/gtk/gtk.go

sed -i 's|PREFIX}${libdir|LIBDIR|; s|libdir|LIBDIR|' \
    Makefile adjust-grub-theme/main.go

%build
export GOPATH=/usr/share/gocode/:$GOPATH
%make_build

%install
for file in $(find . -iname "*.go" -o -iname "*.c" -o -iname "*.h" -o -iname "*.s"); do
    install -d -p %{buildroot}/%{gopath}/src/%{goipath}/$(dirname $file)
    cp -pav $file %{buildroot}/%{gopath}/src/%{goipath}/$file
    echo "%{gopath}/src/%{goipath}/$file" >> devel.file-list
done
%make_install SYSTEMD_SERVICE_DIR="%{_unitdir}" LIBDIR="%{_libexecdir}"
# HOME directory for user deepin-sound-player
mkdir -p %{buildroot}%{_sharedstatedir}/deepin-sound-player

%pre
getent group deepin-sound-player >/dev/null || groupadd -r deepin-sound-player
getent passwd deepin-sound-player >/dev/null || \
    useradd -r -g deepin-sound-player -d %{_sharedstatedir}/deepin-sound-player\
    -s /sbin/nologin \
    -c "User of com.deepin.api.SoundThemePlayer.service" deepin-sound-player
exit 0

%post
%systemd_post deepin-shutdown-sound.service

%preun
%systemd_preun deepin-shutdown-sound.service

%postun
%systemd_postun_with_restart deepin-shutdown-sound.service

%files
%doc README.md
%license LICENSE
%{_bindir}/dde-open
%{_libexecdir}/%{sname}/
%{_unitdir}/*.service
%{_datadir}/dbus-1/services/*.service
%{_datadir}/dbus-1/system-services/*.service
%{_datadir}/dbus-1/system.d/*.conf
%{_datadir}/icons/hicolor/*/actions/*
%{_datadir}/dde-api/data/huangli.db
%{_datadir}/dde-api/data/huangli.version
%{_datadir}/dde-api/data/pkg_depends
%{_datadir}/dde-api/data/grub-themes/
%{_datadir}/polkit-1/actions/com.deepin.api.locale-helper.policy
%{_datadir}/polkit-1/actions/com.deepin.api.device.unblock-bluetooth-devices.policy
%{_var}/lib/polkit-1/localauthority/10-vendor.d/com.deepin.api.device.pkla
%attr(-, deepin-sound-player, deepin-sound-player) %{_sharedstatedir}/deepin-sound-player

%files -n %{name}-devel -f devel.file-list

%changelog
* Thu Mar 18 2021 uoser <uoser@uniontech.com> - 5.4.4-1
- Update to 5.4.4