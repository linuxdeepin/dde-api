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
