// SPDX-FileCopyrightText: 2022 - 2026 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package language_support

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePkgDepends(t *testing.T) {
	pkgDepnds, err := parsePkgDepends("testdata/pkg_depends")
	jsonData, _ := json.Marshal(pkgDepnds)

	assert.Nil(t, err)
	assert.NotNil(t, pkgDepnds)
	assert.NotNil(t, jsonData)
}

func TestLangCodeFromLocale(t *testing.T) {
	locale := langCodeFromLocale("zh_CN")
	assert.Equal(t, locale, "zh-hans")

	locale = langCodeFromLocale("zh_SG")
	assert.Equal(t, locale, "zh-hans")

	locale = langCodeFromLocale("zh_TW")
	assert.Equal(t, locale, "zh-hant")

	locale = langCodeFromLocale("en_US")
	assert.Equal(t, locale, "en")

	locale = langCodeFromLocale("en")
	assert.Equal(t, locale, "en")

	locale = langCodeFromLocale("")
	assert.Equal(t, locale, "")
}

func TestExpendPkgPattern(t *testing.T) {
	pkgs := expendPkgPattern("[p]", "en_US")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]en", "[p]enus", "[p]en-us"})

	pkgs = expendPkgPattern("[p]", "en")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]en"})

	pkgs = expendPkgPattern("[p]", "")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]"})

	pkgs = expendPkgPattern("[p]", "zh_CN")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]zh", "[p]zhcn", "[p]zh-cn", "[p]zh-hans"})

	pkgs = expendPkgPattern("[p]", "zh_SG")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]zh", "[p]zhsg", "[p]zh-sg", "[p]zh-hans"})

	pkgs = expendPkgPattern("[p]", "zh_TW")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]zh", "[p]zhtw", "[p]zh-tw", "[p]zh-hant"})

	pkgs = expendPkgPattern("[p]", "zh_HK")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]zh", "[p]zhhk", "[p]zh-hk", "[p]zh-hant"})

	pkgs = expendPkgPattern("[p]", "wa_BE@euro")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]wa", "[p]wabe", "[p]wa-be", "[p]wa-euro", "[p]wa-be-euro"})
}

type commandCall struct {
	name string
	args []string
}

type fakeCommands struct {
	calls   []commandCall
	dpkgOut string
	dpkgErr error
	aptOut  string
	aptErr  error
}

const (
	dpkgCmdName = "dpkg-query"
	aptCmdName  = "apt-cache"
)

// install stubs both commands with the canned output and records each call.
func (f *fakeCommands) install(t *testing.T) {
	t.Helper()
	prevDpkg, prevApt := dpkgQueryCommand, aptCacheCommand
	t.Cleanup(func() {
		dpkgQueryCommand = prevDpkg
		aptCacheCommand = prevApt
	})
	dpkgQueryCommand = func(args ...string) ([]byte, error) {
		f.calls = append(f.calls, commandCall{name: dpkgCmdName, args: args})
		return []byte(f.dpkgOut), f.dpkgErr
	}
	aptCacheCommand = func(args ...string) ([]byte, error) {
		f.calls = append(f.calls, commandCall{name: aptCmdName, args: args})
		return []byte(f.aptOut), f.aptErr
	}
}

func (f *fakeCommands) callsOf(name string) (out []commandCall) {
	for _, call := range f.calls {
		if call.name == name {
			out = append(out, call)
		}
	}
	return
}

// testLanguageSupport builds a LanguageSupport from pkg_depends content.
func testLanguageSupport(t *testing.T, depends string) *LanguageSupport {
	t.Helper()
	path := filepath.Join(t.TempDir(), "pkg_depends")
	assert.NoError(t, os.WriteFile(path, []byte(depends), 0644))
	pkgDepends, err := parsePkgDepends(path)
	assert.NoError(t, err)
	return &LanguageSupport{pkgDepends: pkgDepends}
}

const batchDepends = `tr::firefox:firefox-locale-
tr::thunderbird:thunderbird-locale-
tr::libreoffice-common:hunspell-
wa:de:libreoffice-common:hunspell-de-at-frami
im:de::fcitx
im:de::fcitx-frontend-gtk3
`

// dpkg-query -W output of a machine with a native and a foreign arch entry for
// the same package, an installed trigger, a removed trigger and an installed
// package that must stay excluded by the dictionary rule.
const batchDpkgOut = `firefox:amd64 installed
firefox-locale-de:amd64 installed
firefox-locale-de:i386 installed
thunderbird:amd64 config-files
libreoffice-common:amd64 installed
hunspell-de-de:amd64 installed
hunspell-de-de:i386 installed
`

// apt-cache policy output for one batch: "(none)" candidates, a package with no
// block at all, and an unrequested foreign-arch block that must not leak in.
const batchAptOut = `firefox-locale-:
  Installed: (none)
  Candidate: (none)
firefox-locale-de:
  Installed: (none)
  Candidate: 1.0
firefox-locale-dede:
  Installed: (none)
  Candidate: (none)
firefox-locale-de-de:
  Installed: (none)
  Candidate: 1.0
hunspell-de-at-frami:
  Installed: (none)
  Candidate: 1.0
hunspell-:
  Installed: (none)
  Candidate: (none)
hunspell-de:
  Installed: (none)
  Candidate: 1.0
hunspell-dede:
  Installed: (none)
  Candidate: (none)
hunspell-de-de:
  Installed: (none)
  Candidate: 1.0
fcitx:
  Installed: (none)
  Candidate: 1.0
fcitx-frontend-gtk3:
  Installed: (none)
  Candidate: 1.0
libc6:i386:
  Installed: (none)
  Candidate: 9.9
`

var batchCandidates = []string{
	"firefox-locale-", "firefox-locale-de", "firefox-locale-dede", "firefox-locale-de-de",
	"hunspell-", "hunspell-de", "hunspell-dede", "hunspell-de-de", "hunspell-de-at-frami",
	"fcitx", "fcitx-frontend-gtk3",
}

func TestByLocaleBatchesQueries(t *testing.T) {
	t.Setenv("XDG_CURRENT_DESKTOP", "DDE")
	fake := &fakeCommands{dpkgOut: batchDpkgOut, aptOut: batchAptOut}
	fake.install(t)
	ls := testLanguageSupport(t, batchDepends)

	packages := ls.ByLocale("de_DE", false)

	assert.ElementsMatch(t, []string{
		"firefox-locale-de-de", "hunspell-de", "hunspell-de-at-frami",
		"fcitx", "fcitx-frontend-gtk3",
	}, packages)

	// one dpkg snapshot and one apt batch per call, no per-package processes
	dpkgCalls := fake.callsOf(dpkgCmdName)
	aptCalls := fake.callsOf(aptCmdName)
	assert.Len(t, dpkgCalls, 1)
	assert.Len(t, aptCalls, 1)
	assert.Equal(t, []string{"-W", "-f", "${binary:Package} ${db:Status-Status}\n"}, dpkgCalls[0].args)

	// every candidate is asked for exactly once, in a single apt invocation
	assert.Equal(t, "policy", aptCalls[0].args[0])
	assert.Equal(t, "--", aptCalls[0].args[1])
	assert.ElementsMatch(t, batchCandidates, aptCalls[0].args[2:])
	assert.Len(t, aptCalls[0].args[2:], len(batchCandidates))
}

func TestByLocaleIncludeInstalled(t *testing.T) {
	t.Setenv("XDG_CURRENT_DESKTOP", "DDE")
	fake := &fakeCommands{dpkgOut: batchDpkgOut, aptOut: batchAptOut}
	fake.install(t)
	ls := testLanguageSupport(t, batchDepends)

	packages := ls.ByLocale("de_DE", true)

	// installed packages show up again, hunspell-de-de stays out (conflicts
	// with hunspell-de-de-frami)
	assert.ElementsMatch(t, []string{
		"firefox-locale-de", "firefox-locale-de-de", "hunspell-de",
		"hunspell-de-at-frami", "fcitx", "fcitx-frontend-gtk3",
	}, packages)
}

func TestByLocaleGnomeExcludesFcitx(t *testing.T) {
	t.Setenv("XDG_CURRENT_DESKTOP", "GNOME")
	fake := &fakeCommands{dpkgOut: batchDpkgOut, aptOut: batchAptOut}
	fake.install(t)
	ls := testLanguageSupport(t, batchDepends)

	packages := ls.ByLocale("de_DE", true)

	assert.ElementsMatch(t, []string{
		"firefox-locale-de", "firefox-locale-de-de", "hunspell-de", "hunspell-de-at-frami",
	}, packages)
}

func TestByPackageAndLocaleBatchesCandidates(t *testing.T) {
	t.Setenv("XDG_CURRENT_DESKTOP", "DDE")
	// duplicated entries must not turn into duplicated apt arguments
	ls := testLanguageSupport(t, `tr::libreoffice-common:hunspell-
wa:de:libreoffice-common:hunspell-de-at-frami
wa:de:libreoffice-common:hunspell-de-at-frami
`)

	fake := &fakeCommands{
		dpkgOut: "hunspell-de-at-frami:amd64 installed\n",
		aptOut: `hunspell-:
  Installed: (none)
  Candidate: (none)
hunspell-de:
  Installed: (none)
  Candidate: 1.0
hunspell-dede:
  Installed: (none)
  Candidate: (none)
hunspell-de-at-frami:
  Installed: (none)
  Candidate: 1.0
`,
	}
	fake.install(t)

	packages := ls.ByPackageAndLocale("libreoffice-common", "de_DE", false)
	assert.ElementsMatch(t, []string{"hunspell-de"}, packages)

	aptCalls := fake.callsOf(aptCmdName)
	assert.Len(t, aptCalls, 1)
	assert.ElementsMatch(t, []string{
		"hunspell-", "hunspell-de", "hunspell-dede", "hunspell-de-de", "hunspell-de-at-frami",
	}, aptCalls[0].args[2:])
	assert.Len(t, aptCalls[0].args[2:], 5) // no duplicate argument

	// installed package kept when asked for
	fake.calls = nil
	packages = ls.ByPackageAndLocale("libreoffice-common", "de_DE", true)
	assert.ElementsMatch(t, []string{"hunspell-de", "hunspell-de-at-frami"}, packages)
	assert.Len(t, fake.callsOf(aptCmdName), 1)
}

func TestQueryFailuresDoNotPanicAndKeepOldFiltering(t *testing.T) {
	t.Setenv("XDG_CURRENT_DESKTOP", "DDE")
	ls := testLanguageSupport(t, `tr::firefox:firefox-locale-
wa:de::hunspell-de-at-frami
`)

	// dpkg fails: no trigger counts as installed, so only the empty trigger
	// contributes, and the failing command is not retried per package
	fake := &fakeCommands{
		dpkgErr: assert.AnError,
		aptOut: `firefox-locale-de:
  Installed: (none)
  Candidate: 1.0
hunspell-de-at-frami:
  Installed: (none)
  Candidate: 1.0
`,
	}
	fake.install(t)
	assert.ElementsMatch(t, []string{"hunspell-de-at-frami"}, ls.ByLocale("de_DE", false))
	assert.Len(t, fake.callsOf(dpkgCmdName), 1)

	// apt fails: nothing is installable
	fake.calls = nil
	fake.dpkgErr = nil
	fake.aptErr = assert.AnError
	assert.Empty(t, ls.ByLocale("de_DE", false))
	assert.Len(t, fake.callsOf(aptCmdName), 1)
}

func TestLanguageSupportRefreshesPackageState(t *testing.T) {
	ls := testLanguageSupport(t, "fn:ja::fonts-example\n")
	fake := &fakeCommands{aptOut: "fonts-example:\n  Candidate: 1.0\n"}
	fake.install(t)
	assert.ElementsMatch(t, []string{"fonts-example"}, ls.ByLocale("ja_JP", false))

	// Reusing the object after installation must not return stale missing packages.
	fake.dpkgOut = "fonts-example installed\n"
	assert.Empty(t, ls.ByLocale("ja_JP", false))

	// Repository changes must also be visible to the next query.
	fake.aptOut = "fonts-example:\n  Candidate: (none)\n"
	assert.Empty(t, ls.ByLocale("ja_JP", true))
}
