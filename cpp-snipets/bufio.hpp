#ifndef _BUFIO_HPP_INCLUDED_SDA9R890Q8095R80UJIHSILUFHUILDHFUILSDFHDUIS_
#define _BUFIO_HPP_INCLUDED_SDA9R890Q8095R80UJIHSILUFHUILDHFUILSDFHDUIS_
#include <cstring>
#include <cstdio>

#include <unistd.h>

#include "line.hpp"

template <size_t N>
class Writer {
    char buf[N];
    int current;
    int out_fd;

public:
    Writer(int fd) {
        current = 0;
        out_fd = fd;
    }

    virtual ~Writer() {
        if (current > 0) {
            flush();
        }
    }

    int flush() {
        if (current == 0) {
            return 0;
        }
        int result = write(out_fd, buf, current);
        if(result < 0) {
            return -1;
        }
        current = 0;
        return 0;
    }

    int append(const char *start, int size) {
        int space_left = N - current;
        if (space_left < size) {
            if (flush() < 0) {
                return -1;
            }
            if (N < size) {
                if(write(out_fd, start, size) < 0) {
                    return -1;
                }
                current = 0;
                return space_left + size;
            }
        }
        memcpy(buf + current, start, size);
        current += size;
        return 0;
    }

    inline int append(const Line &line) {
        return append(line.data(), line.len());
    }

    int add(Writer<N> &add) {
        if((current + add.current) > N) {
            flush();
        }
        memcpy(buf + current, add.buf, add.current);
        current += add.current;
        add.current = 0;
    }
};


template <size_t N>
class Reader {
    char buf[N+1];
    int current;
    int capacity;
    int fd;
    int final;

public:
    Reader(int _fd) {
        current =  0;
        capacity = 0;
        final = 0;
        fd = _fd;
    }

    int getline(Line &res) {
        char *ptr = (char*) memchr(buf + current, '\n', capacity - current);
        if (ptr == NULL) {
            if (final) {
                return 0;
            }
            if (current > 0) {
                memcpy(buf, buf + current, capacity - current);
                capacity -= current;
                current = 0;
            }
            int size = read(fd, buf + capacity, N - capacity);
            if (size < 0) {
                return -1;
            } else if (size + capacity == 0) {
                return 0;
            }
            final = size == 0;
            capacity += size;
            ptr = (char*) memchr(buf, '\n', capacity);
        }
        if(ptr != NULL) {
            res = Line(buf + current, ptr);
            current = ptr - buf + 1;
        } else {
            res = Line(buf + current, buf + capacity);
            current = capacity;
        }
        return 1;
    }
};

#endif // _BUFIO_HPP_INCLUDED_SDA9R890Q8095R80UJIHSILUFHUILDHFUILSDFHDUIS_
