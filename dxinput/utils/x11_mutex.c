// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "x11_mutex.h"

// Global mutex for all X11 operations
pthread_mutex_t x11_global_mutex = PTHREAD_MUTEX_INITIALIZER;
