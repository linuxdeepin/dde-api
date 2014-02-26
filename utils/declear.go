/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

type Manager struct{}

const (
        UTILS_DEST = "com.deepin.api.Utils"
        UTILS_PATH = "/com/deepin/api/Utils"
        UTILS_IFC  = "com.deepin.api.Utils"

        URI_STRING_FILE  = "file://"
        URI_STRING_FTP   = "ftp://"
        URI_STRING_HTTP  = "http://"
        URI_STRING_HTTPS = "https://"
        URI_STRING_SMB   = "smb://"

        URI_TYPE_FILE  = 0
        URI_TYPE_FTP   = 1
        URI_TYPE_HTTP  = 2
        URI_TYPE_HTTPS = 3
        URI_TYPE_SMB   = 4

        ELEMENT_TYPE_INT    = 0
        ELEMENT_TYPE_STRING = 1
        ELEMENT_TYPE_BYTE   = 2
)
