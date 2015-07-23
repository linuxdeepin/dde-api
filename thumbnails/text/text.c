#include <gdk/gdk.h>

#define PIC_SIZE 256

#define FONT_SIZE 12
#define FONT_NAME "sans-serif"

#define TEXT_START_X 15
#define TEXT_START_Y 15
#define TEXT_DELTA_Y (FONT_SIZE+FONT_SIZE/3)

static void do_show_text(cairo_t* cr, char** text);

int
do_gen_thumbnail(char** text, char* dest)
{
	if (!gdk_init_check(NULL, NULL)) {
		g_warning("Init gdk failed");
		return -1;
	}

	cairo_surface_t* surface = cairo_image_surface_create(
			CAIRO_FORMAT_ARGB32, PIC_SIZE, PIC_SIZE);
	if (!surface) {
		g_warning("Create surface failed");
		return -1;
	}

	cairo_t* cr = cairo_create(surface);
	cairo_surface_destroy(surface);
	if (!cr) {
		g_warning("Create cairo failed");
		return -1;
	}

	//background color: gray
	cairo_set_source_rgb(cr, 0.22, 0.22, 0.22);
	cairo_paint(cr);
	do_show_text(cr, text);

	cairo_status_t status = cairo_surface_write_to_png(
			cairo_get_target(cr),
			dest);
	cairo_destroy(cr);
	if (status != CAIRO_STATUS_SUCCESS) {
		g_warning("Write cairo to file '%s' failed", dest);
		return -1;
	}

	return 0;
}

static void
do_show_text(cairo_t* cr, char** text)
{
	cairo_select_font_face(cr, FONT_NAME, 
	                       CAIRO_FONT_SLANT_NORMAL,
	                       CAIRO_FONT_WEIGHT_BOLD);
	cairo_set_font_size(cr, FONT_SIZE);

	// text color: black
	cairo_set_source_rgb(cr, 1, 1, 1);

	int i = 0;
	int height = TEXT_START_Y;
	while (text[i]) {
		if (height > PIC_SIZE - TEXT_START_Y) {
			break;
		}
		cairo_move_to(cr, TEXT_START_X, height);
		cairo_show_text(cr, text[i]);
		height += TEXT_DELTA_Y;
		i++;
	}
}
