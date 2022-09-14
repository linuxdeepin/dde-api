// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
