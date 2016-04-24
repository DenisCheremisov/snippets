#pragma once

#include <algorithm>

namespace io {
    class FileWriter: public Writer {
        int fd;
    public:
        FileWriter(int _fd): fd(_fd) {};

        int64_t write(char *buf, int64_t count) {
            return (int64_t) ::write(fd, buf, (size_t)count);
        }
    };

    class StrWriter: public Writer {
    public:
        char *buf;
        uint64_t capacity;
        uint64_t cur;

        StrWriter(uint64_t _capacity) {
            capacity = _capacity;
            cur = 0;
            buf = new char[capacity];
        }
        ~StrWriter() {
            delete [] buf;
        }

        int64_t write(const char *data, int64_t count) {
            if (count + cur > capacity) {
                uint64_t new_capacity = capacity + std::max(
                    ((uint64_t)count + cur - capacity)*2, capacity/2);
                char *newbuf = new char[new_capacity];
                memcpy(newbuf, buf, cur);
                delete [] buf;
                buf = newbuf;
                capacity = new_capacity;
            }
            memcpy(buf + cur, data, count);
            cur += count;
            return count;
        }
    };

}
