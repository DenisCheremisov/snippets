#include <iostream>
#include <fcntl.h>
#include <unistd.h>
#include <cstring>

class Line {
    const char *buf;
    const char *end;

public:
    Line(): buf(nullptr), end(nullptr) {}
    Line(const char *s, const char *f): buf(s), end(f) {}
    Line(const char *s, int len): buf(s), end(s + len) {}
    Line(const Line &line): buf(line.buf), end(line.end) {}

    inline int len() const {
        return end - buf;
    }

    inline const char *data() const {
        return buf;
    }
};


const int BUFSIZE = 512*1024;

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


class MalformedLine: public std::exception {
private:
    std::string data;
    void _create_basic(const std::string &line) {
        data = "Malformed line `" + line + "`";
    }

public:
    MalformedLine(std::string &&line) {
        _create_basic(line);
    }

    MalformedLine(std::string &&line, std::string &&reason) {
        _create_basic(line);
        data += ": " + reason;
    }

    const char *what() const noexcept {
        return data.c_str();
    }
};


int Atoi(const char *nptr, const char *until) {
    int i = 0;
    int res = 0;
    for (int i = 0; nptr + i < until; i++) {
        res = res*10 + nptr[i] - '0';
    }
    return res;
}


void process(Reader<BUFSIZE> &src) throw(MalformedLine) {
    Line dst;

    const int CAPACITY = 2009;
    int data[CAPACITY];
    int status;
    for(int i = 0; i < sizeof(data)/sizeof(int); i++) {
        data[i] = 0;
    }
    const char *col2start, *col2end, *col3start, *col3end;
    int restLen;
    int k, v;
    std::string reason;

    while (true) {
        status = src.getline(dst);
        if (status <= 0) {
            if (status < 0) {
                perror("Error");
                exit(1);
            }
            break;
        }

        col2start = reinterpret_cast<const char*>(memchr(dst.data(), '\t', dst.len()));
        if (col2start == nullptr) {
            throw MalformedLine(std::string(dst.data(), dst.len()));
        }
        col2start = col2start + 1;
        restLen = dst.len() - (col2start - dst.data());

        col2end = reinterpret_cast<const char*>(memchr(col2start, '\t', restLen));
        if (col2end == nullptr) {
            throw MalformedLine(std::string(dst.data(), dst.len()));
        }
        col3start = col2end + 1;
        restLen = dst.len() - (col3start - dst.data());

        col3end = reinterpret_cast<const char*>(memchr(col3start, '\t', restLen));
        if (col3end == nullptr) {
            throw MalformedLine(std::string(dst.data(), dst.len()));
        }
        k = Atoi(col2start, col2end);
        v = Atoi(col3start, col3end);
        if (k >= CAPACITY) {
            throw MalformedLine(std::string(dst.data(), dst.len()), "2nd column value is out of range");
        }
        data[k] += v;
    }

    k = -1;
    v = 0;
    for(int i = 0; i < sizeof(data)/sizeof(int); i++) {
        if (data[i] > v) {
            k = i;
            v = data[i];
        }
    }
    std::cout << "Max key: " << k << ", max value: " << v << std::endl;
}

int main(int argc, char **argv) {
    int rawSrc = open(argv[1], O_RDONLY);
    if (rawSrc < 0) {
        perror("Error");
    }

    Reader<BUFSIZE> src(rawSrc);
    try {
        process(src);
    } catch (const MalformedLine &exc) {
        std::cerr << exc.what() << std::endl;
        return 1;
    }

    return 0;
}