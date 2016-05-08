#pragma once

#include <algorithm>
#include <cerrno>
#include <cstdint>
#include <cstring>
#include <cstdio>
#include <cassert>

#include <unistd.h>
#include <zlib.h>
#include <fcntl.h>

#include <ioutils/io.hh>

#include <iostream>


namespace io {

    class FileReader: public Reader {
        int fd;
        uint64_t pos;

    public:
        FileReader(int _fd): fd(_fd), pos(0) {};
        ~FileReader() {
            close(fd);
        }

        int64_t read(char *buf, int64_t count) {
            auto bytes_read = (int64_t) ::read(fd, buf, (size_t)count);
            pos += bytes_read;
            readahead(fd, pos, count);
            return bytes_read;
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

    template <uint64_t CHUNK>
    class GzipReader: public Reader {
        FileReader *reader;

        char in_buf[CHUNK];
        char out_buf[CHUNK];
        uint64_t out_pos;
        uint64_t have;

        z_stream strm;
        bool ended;

    public:
        GzipReader(FileReader *_reader) {
            reader = _reader;
            ended = false;
            have = 0;
            out_pos = 0;
            strm.zalloc = Z_NULL;
            strm.zfree = Z_NULL;
            strm.avail_in = 0;
            strm.next_in = Z_NULL;
            assert(inflateInit2(&strm, 16 + MAX_WBITS) == Z_OK);
        }

        int64_t read(char*buf, int64_t n) {
            int64_t read_len = 0;
            int ret_code;
            while (true) {
                if ((have - out_pos) >= n) {
                    memcpy(buf, out_buf + out_pos, n);
                    out_pos += n;
                    return read_len + n;
                }
                if (out_pos < have) {
                    memcpy(buf, out_buf + out_pos, have - out_pos);
                    buf += have - out_pos;
                    n -= have - out_pos;
                    read_len += have - out_pos;
                    out_pos = have = 0;
                } else {
                    out_pos = have = 0;
                }
                if (ended) {
                    return read_len;
                }

                if (strm.avail_in == 0) {
                    auto k = reader->read(in_buf, CHUNK);
                    strm.avail_in = k;
                    if (strm.avail_in <= 0) {
                        inflateEnd(&strm);
                        return strm.avail_in;
                    }
                    strm.next_in = (unsigned char*)in_buf;
                }

                strm.avail_out = CHUNK;
                strm.next_out = (unsigned char*)out_buf;
                ret_code = inflate(&strm, Z_NO_FLUSH);
                switch (ret_code) {
                case Z_STREAM_ERROR:
                case Z_NEED_DICT:
                case Z_DATA_ERROR:
                case Z_MEM_ERROR:
                case Z_BUF_ERROR:
                case Z_VERSION_ERROR:
                    inflateEnd(&strm);
                    return ret_code;
                    break;
                case Z_STREAM_END:
                    ended = true;
                    inflateEnd(&strm);
                }
                have = CHUNK - strm.avail_out;
            }
        }
    };
}
