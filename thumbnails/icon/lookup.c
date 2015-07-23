#include <gtk/gtk.h>

#include "lookup.h"

#define ICON_SIZE 48

char*
lookup_icon(char* theme, char* name)
{
	if (!gtk_init_check(NULL, NULL)) {
		g_warning("Init gtk environment failed");
		return NULL;
	}

	GtkIconTheme* icon_theme = gtk_icon_theme_new();
	gtk_icon_theme_set_custom_theme(icon_theme, theme);
	GtkIconInfo* info = gtk_icon_theme_lookup_icon(icon_theme, 
			name, ICON_SIZE, GTK_ICON_LOOKUP_FORCE_SVG);
	g_object_unref(G_OBJECT(icon_theme));
	if (!info) {
		g_warning("Not fount icon: '%s' in '%s'", name, theme);
		return NULL;
	}

	char* file = g_strdup(gtk_icon_info_get_filename(info));
	g_object_unref(G_OBJECT(info));

	return file;
}
