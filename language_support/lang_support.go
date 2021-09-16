package language_support

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"os/exec"
	"strings"

	libLocale "pkg.deepin.io/lib/locale"
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

func (ls *LanguageSupport) isPkgInstalled(name string) (bool, error) {
	output, err := exec.Command("/usr/bin/dpkg-query", "-W", "-f",
		"${db:Status-Status}", "--", name).Output()
	if err != nil {
		return false, err
	}
	return bytes.Equal(bytes.TrimSpace(output), []byte("installed")), nil
}

func (ls *LanguageSupport) isPkgInstallable(name string) (bool, error) {
	cmd := exec.Command("/usr/bin/apt-cache", "policy",
		"--", name)
	cmd.Env = []string{"LANG=en_US.UTF-8"}
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	const candidate = "Candidate: "
	idx := bytes.Index(output, []byte(candidate))
	if idx == -1 {
		return false, nil
	}

	if bytes.HasPrefix(output[idx+len(candidate):], []byte("(none)")) {
		return false, nil
	}
	return true, nil
}

// ByPackageAndLocale get language support packages for a package and locale.
func (ls *LanguageSupport) ByPackageAndLocale(
	package0 string, locale string, includeInstalled bool) (packages []string) {

	packagesTemp := make(map[string]struct{})
	depMap := ls.pkgDepends[package0]

	// check explicit entries for that locale
	for _, pkgs := range depMap[langCodeFromLocale(locale)] {
		for _, pkg := range pkgs {
			installable, err := ls.isPkgInstallable(pkg)
			if err != nil {
				continue
			}
			if installable {
				packagesTemp[pkg] = struct{}{}
			}
		}
	}

	// check patterns for empty locale string (i. e. applies to any locale)
	for _, patternList := range depMap[""] {
		for _, pattern := range patternList {
			pkgs := expendPkgPattern(pattern, locale)
			for _, pkg := range pkgs {
				installable, err := ls.isPkgInstallable(pkg)
				if err != nil {
					continue
				}
				if installable {
					packagesTemp[pkg] = struct{}{}
				}
			}
		}
	}

	if !includeInstalled { // not show installed
		for pkg := range packagesTemp {
			isInstalled, err := ls.isPkgInstalled(pkg)
			if err != nil {
				continue
			}

			// filter out installed packages
			if isInstalled {
				delete(packagesTemp, pkg)
			}
		}
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

	packages = make([]string, 0, len(packagesTemp))
	for pkg := range packagesTemp {
		packages = append(packages, pkg)
	}
	return
}

// ByLocale get language support packages for a locale
func (ls *LanguageSupport) ByLocale(locale string, includeInstalled bool) []string {
	packagesTemp := make(map[string]struct{})
	for trigger := range ls.pkgDepends {
		var add bool
		if trigger == "" {
			add = true
		} else {
			pkgInstalled, err := ls.isPkgInstalled(trigger)
			if err == nil && pkgInstalled {
				add = true
			}
		}

		if add {
			pkgs := ls.ByPackageAndLocale(trigger, locale, includeInstalled)
			for _, pkg := range pkgs {
				packagesTemp[pkg] = struct{}{}
			}
		}
	}

	packages := make([]string, 0, len(packagesTemp))
	for pkg := range packagesTemp {
		packages = append(packages, pkg)
	}
	return packages
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
