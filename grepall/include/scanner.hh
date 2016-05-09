#pragma once

#include <cstdint>

#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>

#include <ioutils/bufio.hh>
#include <ioutils/reader.hh>
#include <logtools/main.hh>
#include <controller.hh>


namespace main {

    class ContiniousFileReader : public io::Reader {
        logtools::LogMeta *meta;
        io::FileReader *rawReader;
        io::Reader *reader;
        int reader_no;

    public:
        ContiniousFileReader(logtools::LogMeta *m) : meta(m),
                                                     reader(NULL),
                                                     rawReader(NULL),
                                                     reader_no(-1) { }

        ~ContiniousFileReader() {
            if (reader != NULL) {
                delete reader;
            }
            if (rawReader != NULL) {
                delete rawReader;
            }
        }

        int reopen() {
            if (reader != NULL) {
                delete reader;
                reader = NULL;
                if (rawReader != NULL) {
                    delete rawReader;
                    rawReader = NULL;
                }
            }
            while (meta->next()) {
                auto filename = meta->name();
                int fd = open(filename.c_str(), O_RDONLY);
                if (fd > 0) {
                    reader = new io::FileReader(fd);
                    reader_no++;
                    return fd;
                } else {
                    fd = open((filename + ".gz").c_str(), O_RDONLY);
                    if (fd > 0) {
                        rawReader = new io::FileReader(fd);
                        reader = new io::GzipReader<512 * 1024>(rawReader);
                        reader_no++;
                        return fd;
                    }
                }
            }
            return 0;
        }

        int64_t read(char *buf, int64_t count) {
            if (reader == NULL) {
                if (reopen() == 0) {
                    return 0;
                }
            }
            int64_t new_count = reader->read(buf, count);
            if (new_count != 0) {
                return new_count;
            } else {
                int res = reopen();
                if (res == 0) {
                    return 0;
                }
                new_count = reader->read(buf, count);
                return new_count;
            }
        }

        int file_count() const {
            return reader_no;
        }
    };


    // Take care of multiline log records
    class LogScanner {
        logtools::LogMeta *meta;
        bufio::Reader *reader;
        char *buf;
        uint64_t cur;
        uint64_t lineStart;
        uint64_t capacity;
        time_t prevStamp;

    public:
        enum scanner_read_t {
            EMPTY = 0,
            OK,
            END,
            ERROR
        };

        LogScanner(logtools::LogMeta *m, bufio::Reader *r, uint64_t cap) {
            meta = m;
            reader = r;
            capacity = cap;
            buf = new char[capacity];
            cur = 0;
            lineStart = 0;
            prevStamp = -1;
        }

        ~LogScanner() {
            delete [] buf;
        }

        scanner_read_t readline(bufio::Buf *dst, time_t stamp) {
            while (true) {
                if (prevStamp > stamp) {
                    return EMPTY;
                }
                bufio::Buf tmp;
                auto res = reader->readline(&tmp);
                time_t curStamp;
                switch (res) {
                case bufio::_OK:
                    if (tmp.len + 1> capacity - cur) {
                        memmove(buf, buf + lineStart, cur - lineStart);
                        cur = cur - lineStart;
                        lineStart = 0;
                        if (tmp.len + 1 > capacity - cur) {
                            uint64_t newcap =
                                std::max(3*capacity/2, capacity + 2*tmp.len + 2);
                            char *newbuf = new char[newcap];
                            memcpy(newbuf, buf, cur);
                            capacity = newcap;
                            delete [] buf;
                            buf = newbuf;
                        }
                    }

                    curStamp = meta->line_stamp(tmp.buf, tmp.len);
                    if (curStamp >= 0 && prevStamp >= 0) {
                        dst->buf = buf + lineStart;
                        dst->len = cur - lineStart;
                        memcpy(buf + cur, tmp.buf, tmp.len);
                        lineStart = cur;
                        cur += tmp.len;
                        prevStamp = curStamp;
                        return OK;
                    } else {
                        int raise = 0;
                        if (prevStamp >= 0) {
                            buf[cur] = '\n';
                            raise = 1;
                        } else {
                            prevStamp = curStamp;
                        }
                        memcpy(buf + cur + raise, tmp.buf, tmp.len);
                        cur += tmp.len + raise;
                    }
                    break;

                case bufio::_EOF:
                    if (lineStart == cur) {
                        return END;
                    }
                    dst->buf = buf + lineStart;
                    dst->len = cur - lineStart;
                    lineStart = cur;
                    return OK;

                case bufio::_ERR:
                    return ERROR;
                }
            }
        }
    };


    class Filter {
    public:
        virtual bool match(const char *line, uint64_t len) = 0;
    };


    class FilteringLineReader: public LineReader {
        LogScanner *scanner;
        Filter *filter;

    public:
        FilteringLineReader(LogScanner *s, Filter *f) {
            scanner = s;
            filter = f;
        }

        int get_lines(Visitor *v, time_t stamp) {
            bufio::Buf line;
            while (true) {
                switch (scanner->readline(&line, stamp)) {
                case LogScanner::OK:
                    if (filter->match(line.buf, line.len)) {
                        v->append_line(line.buf, line.len);
                    }
                    continue;
                case LogScanner::EMPTY:
                    return 1;
                case LogScanner::END:
                    return 0;
                case LogScanner::ERROR:
                    return -1;
                }
            }
        }
    };
}
