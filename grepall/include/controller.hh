#pragma once

#include <cstdint>
#include <vector>
#include <ctime>

#include <ioutils/io.hh>

namespace main {
    class Visitor;

    class LineReader {
    public:
        virtual int get_lines(Visitor *v, time_t stamp) = 0;
    };

    class Visitor {
        io::Writer *writer;

        char *buf;
        uint64_t capacity;
        uint64_t cur;

        std::vector<LineReader*> readers;

    public:
        Visitor(io::Writer *w, uint64_t cap=512*1024) {
            writer = w;
            capacity = cap;
            buf = new char[capacity];
            cur = 0;
        }
        ~Visitor() {
            delete [] buf;
        }

        void add_reader(LineReader *reader) {
            readers.push_back(reader);
        }

        void append_line(const char *data, const uint64_t len) {
            uint64_t l = len + 1;
            if (capacity - cur < l) {
                flush();
                if (capacity < l) {
                    capacity = (3*l + 1)/2;
                    char *newbuf = new char[capacity];
                    memcpy(newbuf, buf, cur);
                    delete [] buf;
                    buf = newbuf;
                }
            }
            memcpy(buf + cur, data, len);
            buf[cur + len] = '\n';
            cur += l;
        }

        void flush() {
            writer->write(buf, cur);
            cur = 0;
        }

        int visit_all(time_t stamp) {
            uint64_t lines = 0;
            if(readers.size() == 0) {
                return -1;
            }
            for(int i = 0; i < readers.size();) {
                int res = readers[i]->get_lines(this, stamp);
                if (res < 0) {
                    readers.erase(readers.begin() + i);
                } else {
                    i++;
                    lines += res;
                }
            }
            return lines;
        }
    };
}
