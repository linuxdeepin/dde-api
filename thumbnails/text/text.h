#ifndef __TEXT_H__
#define __TEXT_H__

typedef struct _THUMB_INFO
{
    int width;
    int height;
    int xborder;
    int yborder;
    int canvasWidth;
    int canvasHeight;
    int fontSize;
} ThumbInfo;

int do_gen_thumbnail(char** text, char* dest, ThumbInfo* info);

#endif
