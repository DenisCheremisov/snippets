//
// Created by emacs on 11.05.16.
//

#pragma once

#include <cstring>

#include <scanner.hh>

namespace main {
    class FixedFilter: public Filter {
        const char *filter;
        const uint64_t length;

    public:
        FixedFilter(const char *f): filter(f), length(strlen(f)) {}

        bool match(const char *line, uint64_t len) {
            return memmem(line, len, filter, length) != NULL;
        }
    };
}

