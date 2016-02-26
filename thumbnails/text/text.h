/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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

int text_thumbnail(char** text, char* dest, ThumbInfo* info);

#endif
