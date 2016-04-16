#pragma once

#include <ioutils/io.hh>
#include <cstring>



namespace bufio {
    enum read_result_t {
        _OK = 0,
        _EOF,
        _ERR
    };

    class Buf {
    public:
        char *buf;
        uint64_t len;

        Buf() {}
    };

    class Reader {
        char *buf;
        int capacity;
        int curlength;
        int current;
        io::Reader *raw_reader;

    public:
        Reader(io::Reader *rd, uint64_t bufsize) {
            buf = new char[bufsize];
            capacity = bufsize;
            current = 0;
            curlength = 0,
            raw_reader = rd;
        }
        ~Reader() {
            delete [] buf;
        }

        read_result_t readline(Buf *dst) {
            char *ptr = (char*) memchr(buf + current, '\n', curlength - current);
            while (ptr == NULL) {
                if(current > 0) {
                    memcpy(buf, buf + current, curlength - current);
                    curlength -= current;
                    current = 0;
                }
                int size = raw_reader->read(buf + curlength, capacity - curlength);
                if (size < 0) {
                    return _ERR;
                }
                if (size == 0) {
                    break;
                }
                curlength += size;
                ptr = (char*) memchr(buf, '\n', curlength);
                if (ptr != NULL) {
                    break;
                }
                char *newbuf = new char[capacity*2];
                capacity *= 2;
                memcpy(newbuf, buf, curlength);
                delete [] buf;
                buf = newbuf;
            }
            dst->buf = buf + current;
            if (ptr != NULL) {
                dst->len = ptr - buf - current;
                current = ptr - buf + 1;
                return _OK;
            } else {
                if (current == curlength) {
                    return _EOF;
                }
                dst->len = curlength - current;
                current = curlength;
                return _OK;
            }
        }
    };
}
