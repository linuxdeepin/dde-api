package i18n_dependent

var (
	conflictPkgMap = map[string][]string{
		"fonts-adobe-source-han-sans-cn": []string{
			"fonts-adobe-source-han-sans-tw",
			"fonts-adobe-source-han-sans-jp",
			"fonts-adobe-source-han-sans-kr",
		},
		"fonts-adobe-source-han-sans-tw": []string{
			"fonts-adobe-source-han-sans-cn",
			"fonts-adobe-source-han-sans-jp",
			"fonts-adobe-source-han-sans-kr",
		},
		"fonts-adobe-source-han-sans-jp": []string{
			"fonts-adobe-source-han-sans-cn",
			"fonts-adobe-source-han-sans-tw",
			"fonts-adobe-source-han-sans-kr",
		},
		"fonts-adobe-source-han-sans-kr": []string{
			"fonts-adobe-source-han-sans-cn",
			"fonts-adobe-source-han-sans-tw",
			"fonts-adobe-source-han-sans-jp",
		},
	}
)
