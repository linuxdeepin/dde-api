// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
