# SPDX-FileCopyrightText: 2026 UnionTech Software Technology Co., Ltd.
#
# SPDX-License-Identifier: LGPL-3.0-or-later

# DdeApiConfig.cmake
# Provides DdeApi::EventLogger imported target (header-only)
#
# Usage:
#   find_package(DdeApi QUIET)
#   if(DdeApi_FOUND)
#     target_link_libraries(myapp PRIVATE DdeApi::EventLogger)
#   endif()

include_guard(GLOBAL)

set(DDE_API_INCLUDE_DIRS "/usr/include")

if(NOT TARGET DDEAPI::EventLogger)
    add_library(DDEAPI::EventLogger INTERFACE IMPORTED)
    set_target_properties(DDEAPI::EventLogger PROPERTIES
        INTERFACE_INCLUDE_DIRECTORIES "${DDE_API_INCLUDE_DIRS}"
    )
endif()

set(DdeApi_FOUND TRUE)
