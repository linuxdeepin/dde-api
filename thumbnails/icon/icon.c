/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include <gtk/gtk.h>
#include "icon.h"

char *choose_icon(char *theme, const char **names, int size)
{
	if (!gtk_init_check(NULL, NULL)) {
		g_warning("Init gtk environment failed");
		return NULL;
	}

	GtkIconTheme *icon_theme = gtk_icon_theme_new();
	gtk_icon_theme_set_custom_theme(icon_theme, theme);
	GtkIconInfo *info = gtk_icon_theme_choose_icon(icon_theme,
						       names, size, 0);
	g_object_unref(G_OBJECT(icon_theme));
	if (!info) {
		return NULL;
	}

	char *file = g_strdup(gtk_icon_info_get_filename(info));
	g_object_unref(G_OBJECT(info));

	return file;
}
