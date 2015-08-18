/**
 * PDF thumbnail generator
 *
 * Reference xfce tumbler.
 **/
#include <glib/poppler-document.h>
#include <glib/poppler-page.h>

#include <cairo.h>

static PopplerDocument* create_poppler_document(gchar* uri);
static cairo_surface_t* get_thumbnail_surface(PopplerDocument* doc, gint index);
static gint save_thumbnail(cairo_surface_t* surface, gchar* dest);
static cairo_surface_t* get_thumbnail_surface_from_page(PopplerPage* page);

int
gen_thumbnail(char* uri, char* dest)
{
        PopplerDocument* doc = create_poppler_document(uri);
        if (!doc) {
                return -1;
        }

        // get the first page surface
        cairo_surface_t* surface = get_thumbnail_surface(doc, 0);
        g_object_unref(doc);
        if (!surface) {
                return -1;
        }

        int ret = save_thumbnail(surface, dest);
        cairo_surface_destroy(surface);

        return ret;
}

static PopplerDocument*
create_poppler_document(gchar* uri)
{
    GError* error = NULL;
    PopplerDocument* doc = poppler_document_new_from_file(uri, NULL, &error);
    // TODO: if doc == NULL, create PopplerDocument from file contents
    if (error) {
        g_print("Open file failed: %s\n", error->message);
        g_error_free(error);
        return NULL;
    }

    return doc;
}

static gint
save_thumbnail(cairo_surface_t* surface, gchar* dest){
        cairo_status_t ret = cairo_surface_write_to_png(surface, dest);
        if (ret != CAIRO_STATUS_SUCCESS) {
                return -1;
        }

        return 0;
}

static cairo_surface_t*
get_thumbnail_surface(PopplerDocument* doc, gint index)
{
    PopplerPage* page = poppler_document_get_page(doc, index);
    if (!page) {
        g_printerr("Get the '%d' page failed\n", index);
        return NULL;
    }

    cairo_surface_t* surface = poppler_page_get_thumbnail(page);
    if (!surface) {
            surface = get_thumbnail_surface_from_page(page);
    }

    g_object_unref(page);
    return surface;
}

static cairo_surface_t*
get_thumbnail_surface_from_page(PopplerPage* page)
{
    gdouble width, height;
    poppler_page_get_size(page, &width, &height);

    cairo_surface_t* surface = cairo_image_surface_create(CAIRO_FORMAT_ARGB32,
                                                          width,
                                                          height);
    cairo_t* cr = cairo_create(surface);
    cairo_save(cr);
    poppler_page_render(page, cr);
    cairo_restore(cr);

    cairo_set_operator(cr, CAIRO_OPERATOR_DEST_OVER);
    cairo_set_source_rgb(cr, 1.0, 1.0, 1.0);
    cairo_paint(cr);
    cairo_destroy(cr);

    return surface;
}
