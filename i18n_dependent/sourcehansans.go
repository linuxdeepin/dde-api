/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
