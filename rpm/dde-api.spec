# Run tests in check section
# disable for bootstrapping
%bcond_with check
%global goipath  pkg.deepin.io/dde/api

%global provider        pkg
%global provider_tld    deepin
%global project         io
%global repo            dde/api

%global provider_prefix %{provider}.%{provider_tld}.%{project}/%{repo}
%global import_path     %{provider_prefix}

%global with_debug 1

%if 0%{?with_debug}
%global debug_package   %{nil}
%endif
%gometa

Name:           dde-api
Version:        5.3.0
Release:        1
Summary:        Go bindings for Deepin Desktop Environment development
License:        GPLv3
URL:            %{gourl}
Source0:        %{name}_%{version}.orig.tar.xz
BuildRequires:  compiler(go-compiler)

%description
DLib is a set of Go bindings/libraries for DDE development.
Containing dbus (forking from guelfey), glib, gdkpixbuf, pulse and more.

%package devel
Summary:        %{summary}
BuildArch:      noarch

BuildRequires:  deepin-gir-generator
BuildRequires:  dbus-x11
BuildRequires:  iso-codes
BuildRequires:  mobile-broadband-provider-info
BuildRequires:  pkgconfig(gio-2.0)
BuildRequires:  pkgconfig(gdk-3.0)
BuildRequires:  pkgconfig(gdk-x11-3.0)
BuildRequires:  pkgconfig(gdk-pixbuf-xlib-2.0)
BuildRequires:  pkgconfig(libpulse)

%description devel
%{summary}.

Provides: golang(pkg.deepin.io/lib)

This package contains library source intended for
building other packages which use import path with
%{goipath} prefix.

%prep
%setup -q -n %{name}-%{version}
%autosetup

%install
install -d -p %{buildroot}/%{gopath}/src/%{import_path}/
for file in $(find . -iname "*.go") ; do
    install -d -p %{buildroot}/%{gopath}/src/%{import_path}/$(dirname $file)
    cp -pav $file %{buildroot}/%{gopath}/src/%{import_path}/$file
    echo "%%{gopath}/src/%%{import_path}/$file" >> devel.file-list
done

%if %{with check}
%check
%gochecks
%endif

%files devel -f devel.file-list
%doc README.md
%license LICENSE

%changelog
* Thu Jun 11 2020 uoser <uoser@uniontech.com> - 5.4.5
- Update to 5.4.5
