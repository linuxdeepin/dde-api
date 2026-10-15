// SPDX-FileCopyrightText: 2022 - 2026 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package language_support

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	libLocale "github.com/linuxdeepin/go-lib/locale"
)

type LanguageSupport struct {
	pkgDepends map[string]map[string]map[string][]string
	//nolint
	langCountryMap int
}

func NewLanguageSupport() (ls *LanguageSupport, err error) {
	ls = &LanguageSupport{}

	ls.pkgDepends, err = parsePkgDepends(defaultDependsFile)
	if err != nil {
		return nil, err
	}

	return ls, nil
}

func (ls *LanguageSupport) Destroy() {
}

// commandTimeout bounds a single dpkg-query or apt-cache run.
const commandTimeout = 30 * time.Second

// dpkgQueryCommand and aptCacheCommand are variables so that tests can supply
// canned output.
var dpkgQueryCommand = func(args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	return exec.CommandContext(ctx, "/usr/bin/dpkg-query", args...).Output()
}

var aptCacheCommand = func(args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/usr/bin/apt-cache", args...)
	// Only LC_ALL=C: a translated apt prints localized field names.
	cmd.Env = []string{"LC_ALL=C"}
	return cmd.Output()
}

// pkgQuery is the dpkg/apt snapshot of a single call: one dpkg-query and one
// apt-cache process, not cached across calls.
type pkgQuery struct {
	installed  map[string]bool
	candidates map[string]bool
}

func newPkgQuery() *pkgQuery {
	return &pkgQuery{candidates: make(map[string]bool)}
}

// isInstalled reads the whole dpkg database once; a multi-arch entry
// ("name:arch") also marks the plain name.
func (q *pkgQuery) isInstalled(name string) bool {
	if q.installed == nil {
		q.installed = make(map[string]bool)
		output, err := dpkgQueryCommand("-W", "-f",
			"${binary:Package} ${db:Status-Status}\n")
		if err == nil {
			for _, line := range strings.Split(string(output), "\n") {
				pkg, status, found := strings.Cut(line, " ")
				if !found || status != "installed" {
					continue
				}
				q.installed[pkg] = true
				if i := strings.IndexByte(pkg, ':'); i > 0 {
					q.installed[pkg[:i]] = true
				}
			}
		}
	}
	return q.installed[name]
}

// loadCandidates asks apt once for all given names; unknown names and
// "(none)" candidates stay unmarked.
func (q *pkgQuery) loadCandidates(names []string) {
	if len(names) == 0 {
		return
	}
	args := append([]string{"policy", "--"}, names...)
	output, err := aptCacheCommand(args...)
	if err != nil {
		return
	}

	var pkg string
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}
		if line[0] != ' ' && strings.HasSuffix(line, ":") { // "package:" header
			pkg = strings.TrimSuffix(line, ":")
			continue
		}
		version, found := strings.CutPrefix(strings.TrimSpace(line), "Candidate:")
		if !found || pkg == "" {
			continue
		}
		version = strings.TrimSpace(version)
		q.candidates[pkg] = version != "" && version != "(none)"
	}
}

// candidateNames lists the packages offered for one trigger and locale.
func (ls *LanguageSupport) candidateNames(package0, locale string) (names []string) {
	depMap := ls.pkgDepends[package0]

	// check explicit entries for that locale
	for _, pkgs := range depMap[langCodeFromLocale(locale)] {
		names = append(names, pkgs...)
	}

	// check patterns for empty locale string (i. e. applies to any locale)
	for _, patternList := range depMap[""] {
		for _, pattern := range patternList {
			names = append(names, expendPkgPattern(pattern, locale)...)
		}
	}
	return
}

// packagesFor batches the candidates of all given triggers and applies the
// filtering rules to the result.
func (ls *LanguageSupport) packagesFor(
	triggers []string, locale string, includeInstalled bool, q *pkgQuery) []string {

	names := make([]string, 0)
	seen := make(map[string]struct{})
	for _, trigger := range triggers {
		for _, name := range ls.candidateNames(trigger, locale) {
			if _, dup := seen[name]; dup {
				continue
			}
			seen[name] = struct{}{}
			names = append(names, name)
		}
	}

	q.loadCandidates(names)

	packagesTemp := make(map[string]struct{})
	for _, name := range names {
		if !q.candidates[name] {
			continue
		}
		if !includeInstalled && q.isInstalled(name) { // not show installed
			continue
		}
		packagesTemp[name] = struct{}{}
	}

	// exclude Fcitx packages if GNOME desktop
	desktop := os.Getenv("XDG_CURRENT_DESKTOP")
	var noFcitx bool
	for _, desktopItem := range strings.Split(desktop, ":") {
		if desktopItem == "GNOME" {
			noFcitx = true
			break
		}
	}

	if noFcitx {
		for pkg := range packagesTemp {
			if strings.HasPrefix(pkg, "fcitx") {
				delete(packagesTemp, pkg)
			}
		}
	}

	// exclude hunspell-de-XX since they conflict with -frami
	for _, country := range []string{"de", "at", "ch"} {
		delete(packagesTemp, "hunspell-de-"+country)
	}

	// exclude hunspell-gl since it conflicts with hunspell-gl-es
	// https://launchpad.net/bugs/1578821
	delete(packagesTemp, "hunspell-gl")

	packages := make([]string, 0, len(packagesTemp))
	for pkg := range packagesTemp {
		packages = append(packages, pkg)
	}
	return packages
}

// ByPackageAndLocale get language support packages for a package and locale.
func (ls *LanguageSupport) ByPackageAndLocale(
	package0 string, locale string, includeInstalled bool) (packages []string) {

	return ls.packagesFor([]string{package0}, locale, includeInstalled, newPkgQuery())
}

// ByLocale get language support packages for a locale
func (ls *LanguageSupport) ByLocale(locale string, includeInstalled bool) []string {
	q := newPkgQuery()

	var triggers []string
	for trigger := range ls.pkgDepends {
		if trigger == "" || q.isInstalled(trigger) {
			triggers = append(triggers, trigger)
		}
	}

	return ls.packagesFor(triggers, locale, includeInstalled, q)
}

func expendPkgPattern(pattern, locale string) (pkgs []string) {
	comp := libLocale.ExplodeLocale(locale)
	lang := strings.ToLower(comp.Language)
	country := strings.ToLower(comp.Territory)
	variant := strings.ToLower(comp.Modifier)

	pkgs = []string{pattern, pattern + lang}

	if country != "" {
		pkgs = append(pkgs,
			pattern+lang+country,
			pattern+lang+"-"+country)
	}

	if variant != "" {
		pkgs = append(pkgs, pattern+lang+"-"+variant)
	}

	if country != "" && variant != "" {
		pkgs = append(pkgs, pattern+lang+"-"+country+"-"+variant)
	}

	if lang == "zh" {
		if country == "cn" || country == "sg" {
			pkgs = append(pkgs, pattern+"zh-hans")
		} else {
			pkgs = append(pkgs, pattern+"zh-hant")
		}
	}
	return
}

func langCodeFromLocale(locale string) string {
	if strings.HasPrefix(locale, "zh_CN") || strings.HasPrefix(locale, "zh_SG") {
		return "zh-hans"
	}
	if strings.HasPrefix(locale, "zh_") {
		return "zh-hant"
	}

	parts := strings.SplitN(locale, "_", 2)
	return parts[0]
}

const defaultDependsFile = "/usr/share/dde-api/data/pkg_depends"

// parse pkg_depends file
func parsePkgDepends(filename string) (ret map[string]map[string]map[string][]string, err error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fh.Close()
	}()

	reader := csv.NewReader(fh)
	reader.Comma = ':'
	reader.Comment = '#'
	reader.FieldsPerRecord = 4
	reader.TrimLeadingSpace = true

	//          trigger    langCode    category dependency
	ret = make(map[string]map[string]map[string][]string)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		category := record[0]
		langCode := record[1]
		trigger := record[2]
		dependency := record[3]

		langCodeMap := ret[trigger]
		if langCodeMap == nil {
			langCodeMap = make(map[string]map[string][]string)
			ret[trigger] = langCodeMap
		}

		categoryCodeMap := langCodeMap[langCode]
		if categoryCodeMap == nil {
			categoryCodeMap = make(map[string][]string)
			langCodeMap[langCode] = categoryCodeMap
		}

		dependencySlice := categoryCodeMap[category]
		categoryCodeMap[category] = append(dependencySlice, dependency)
	}
	return ret, nil
}
