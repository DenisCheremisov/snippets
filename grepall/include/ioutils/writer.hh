#pragma once

#include <algorithm>
#include <string>

#include <ioutils/io.hh>

namespace io {
    class FileWriter: public Writer {
        int fd;
    public:
        FileWriter(int _fd): fd(_fd) {};
        ~FileWriter() {
            close(fd);
        }

        int64_t write(char *buf, int64_t count) {
            return (int64_t) ::write(fd, buf, (size_t)count);
        }
    };

    class StrWriter: public Writer {
    public:
        std::string buf;

        int64_t write(const char *data, int64_t count) {
            buf += std::string(data, count);
            return count;
        }
    };

}
