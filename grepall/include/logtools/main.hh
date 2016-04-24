#pragma once

#include <string>
#include <ctime>

namespace logtools {
    class LogMeta {
    public:
        virtual ~LogMeta() {};

        virtual bool next() = 0;
        virtual std::string name() = 0;
        virtual time_t line_stamp(const char *line, int len) = 0;
    };
}
