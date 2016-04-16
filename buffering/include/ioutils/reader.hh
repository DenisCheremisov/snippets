#pragma once

#include <algorithm>
#include <cerrno>
#include <stdint.h>
#include <cstring>

#include <unistd.h>

#include <ioutils/io.hh>

#include <iostream>

namespace io {

    class FileReader: public Reader {
        int fd;
    public:
        FileReader(int _fd): fd(_fd) {};

        int64_t read(char *buf, int64_t count) {
            return (int64_t) ::read(fd, buf, (size_t)count);
        }
    };

    class StrReader: public Reader {
        char *buf;
        uint64_t capacity;
        uint64_t pos;

    public:
        StrReader(const char *data) {
            capacity = strlen(data);
            buf = new char[capacity];
            pos = 0;
            memcpy(buf, data, capacity);
        }
        ~StrReader() {
            delete [] buf;
        }

        int64_t read(char *_buf, int64_t count) {
            uint64_t delta = std::min(capacity - pos, (uint64_t)count);
            memcpy(_buf, this->buf + pos, delta);
            pos += delta;
            return delta;
        }
    };
}
