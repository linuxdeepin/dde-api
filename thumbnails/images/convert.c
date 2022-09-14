// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include <gdk/gdk.h>
#include <librsvg/rsvg.h>

#include "convert.h"

int
svg_to_png(const char* file, const char* dest)
{
	if (!gdk_init_check(NULL, NULL)) {
		g_warning("Init gdk environment failed");
		return -1;
	}

	GError* error = NULL;
	RsvgHandle* handler = rsvg_handle_new_from_file(file, &error);
	if (error) {
		g_warning("New RsvgHandle failed: %s", error->message);
		g_error_free(error);
		return -1;
	}

	GdkPixbuf* pbuf = rsvg_handle_get_pixbuf(handler);
	g_object_unref(G_OBJECT(handler));

	error = NULL;
	gdk_pixbuf_save(pbuf, dest, "png", &error, NULL);
	g_object_unref(G_OBJECT(pbuf));
	if (error) {
		g_warning("Save to png file failed: %s", error->message);
		g_error_free(error);
		return -1;
	}

	return 0;
}
