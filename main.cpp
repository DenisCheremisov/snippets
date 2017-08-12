#include <iostream>
#include <list>
#include <cstring>

#include <fcntl.h>
#include <unistd.h>

#include <google/dense_hash_map>
#include <boost/program_options.hpp>


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

    void rebuild(const char *data, int lenght) {
        buf = data;
        end = data + lenght;
    }

    void rebuild(const char *data, const char *dataEnd) {
        buf = data;
        end = dataEnd;
    }

    bool operator==(const Line &another) const {
        return len() == another.len() && memcmp(buf, another.buf, len()) == 0;
    }

    size_t hash() const {
        unsigned int hash = 1315423911;
        for (int i = 0; i < len(); i++) {
            hash ^= ((hash << 5) + buf[i] + (hash >> 2));
        }
        return (hash & 0x7FFFFFFF);
    }

    std::string toString() const {
        if (buf == nullptr) {
            return "<empty>";
        }
        return std::string(buf, len());
    }
};

struct lineEquator {
    bool operator()(const Line &line1, const Line &line2) const {
        return line1 == line2;
    }
};

namespace std {
    template<>
    struct hash<Line> {
        std::size_t operator()(const Line &k) const {
            return k.hash();
        }
    };
}

class LinePool {
    char *storage;
    int cap;
    int taken;

    std::list<const char*> chunks;

    char *allocateChunk(int capacity) {
        return new char[capacity];
    }

    int left() {
        return cap - taken;
    }

public:
    LinePool(int initCapacity) {
        storage = allocateChunk(initCapacity);
        cap = initCapacity;
        taken = 0;
    }

    ~LinePool() {
        delete []storage;
        for(auto it: chunks) {
            delete []it;
        }
    }

    void solid_alloc(Line &line) {
        if (line.len() > left()) {
            chunks.push_back(storage);
            cap = std::max(2 * cap, 2 * line.len());
            taken = 0;
            storage = allocateChunk(cap);
        }
        memcpy(storage + taken, line.data(), line.len());
        line.rebuild(storage + taken, line.len());
        taken += line.len();
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


class TSVScanner {
    Line orig;
    const char *dataStart;
    const char *dataEnd;

    int cols;

    int curCol;

public:
    TSVScanner(int colLimit): cols(colLimit) {}

    void set_source(const Line &src) {
        orig = src;
        dataStart = src.data();
        dataEnd = src.data() + src.len();
        curCol = 1;
    }

    // get_column returns bool if diff_index points to the columns that is out of bound set
    // and throws MalformedLine exception if the line doesn't have enough columns
    bool get_column(int index, Line &dest) throw(MalformedLine) {
        if (index >= cols) {
            return false;
        }

        for(int i = curCol; i < index; i++) {
            const char *pos = reinterpret_cast<const char*>(memchr(dataStart, '\t', dataEnd - dataStart));
            if (pos == nullptr) {
                throw MalformedLine(orig.toString(), "not enough columns");
            }
            dataStart = pos + 1;
        }

	curCol = index;
        if (curCol == cols) {
            dest.rebuild(dataStart, dataEnd);
            dataStart = dataEnd;
            return true;
        }

        const char *end = reinterpret_cast<const char*>(memchr(dataStart, '\t', dataEnd - dataStart));
        if (end == nullptr) {
            throw MalformedLine(orig.toString(), "not enough columns");
        }
        dest.rebuild(dataStart, end);
        dataStart = end + 1;
        curCol++;
        return true;
    }
};


int Atoi(const char *nptr, const char *until) {
    int res = 0;
    for (int i = 0; nptr + i < until; i++) {
        res = res * 10 + nptr[i] - '0';
    }
    return res;
}

int Atoi(const Line &src) {
    return Atoi(src.data(), src.data() + src.len());
}



std::pair<std::string, int> process(
        Reader<BUFSIZE> &src, int keyIndex, int valIndex, int colLimit) throw(MalformedLine) {

    Line dst;

    int status;
    int v;
    TSVScanner scanner(colLimit);
    Line tmp1, tmp2, key;
    int indices[2]{std::min(keyIndex, valIndex), std::max(keyIndex, valIndex)};
    LinePool pool(512*1024); // 512Kb
    google::dense_hash_map<Line, int, std::hash<Line>, lineEquator> storage;
    storage.set_empty_key(Line());

    while (true) {
        status = src.getline(dst);
        if (status <= 0) {
            if (status < 0) {
                perror("Error");
                exit(1);
            }
            break;
        }

        scanner.set_source(dst);
        scanner.get_column(indices[0], tmp1);
        scanner.get_column(indices[1], tmp2);
        if (keyIndex < valIndex) {
            key = tmp1;
            v = Atoi(tmp2);
        } else {
            key = tmp2;
            v = Atoi(tmp1);
        }

        auto keyPos = storage.find(key);
        if (keyPos != storage.end()) {
            keyPos.pos->second += v;
        } else {
            pool.solid_alloc(key);
            storage[key] = v;
        }
    }

    if (storage.size() == 0) {
        return std::pair<std::string, int>("", 0);
    }
    v = 0;
    key = Line();
    for (auto it: storage) {
        if (it.second > v) {
            key = it.first;
            v = it.second;
        }
    }
    return std::pair<std::string, int>(key.toString(), v);
}

namespace po = boost::program_options;


void printUsage(const char **argv, const po::options_description &desc) {
    std::cerr << "Usage: " << argv[0] << " [options] <input file name>" << std::endl;
    std::cerr << desc << std::endl;
}

int main(int argc, const char **argv) {
    int keyIndex, valIndex, colLimit;
    po::options_description desc("Options");
    desc.add_options()
            ("help", "produce help message")
            ("key", po::value<int>(&keyIndex)->default_value(2), "key column index (1, 2, 3, ...)")
            ("val", po::value<int>(&valIndex)->default_value(3), "value column index (1, 2, 3, ...)")
            ("limit", po::value<int>(&colLimit)->default_value(4), "how many columns to scan");
    po::options_description hidden("Hidden options");
    hidden.add_options()("input-file", po::value<std::string>(), "input file");
    po::options_description cmdlineOptions;
    cmdlineOptions.add(desc).add(hidden);

    po::positional_options_description p;
    p.add("input-file", -1);

    po::variables_map pm;
    try {
        po::store(po::command_line_parser(argc, argv).options(cmdlineOptions).positional(p).run(), pm);
        po::notify(pm);
    } catch(const std::exception &exc) {
        printUsage(argv, desc);
        std::cerr << "\033[31m" << exc.what() << "\033[0m\n";
        return 1;
    }

    if (pm.count("help")) {
        printUsage(argv, desc);
        return 0;
    }

    int fileIndex = pm.count("input-file");
    if (!fileIndex) {
        printUsage(argv, desc);
        std::cerr << "\033[31mmissing input file name positional parameter\033[0m\n";
        return 1;
    }
    int rawSrc = open(argv[fileIndex], O_RDONLY);
    if (rawSrc < 0) {
        perror("Error");
    }

    Reader<BUFSIZE> src(rawSrc);
    try {
        auto res = process(src, keyIndex, valIndex, colLimit);
        std::cout << "Max key: " << res.first << ", max value: " << res.second << std::endl;
    } catch (const MalformedLine &exc) {
        std::cerr << exc.what() << std::endl;
        return 1;
    }

    return 0;
}
