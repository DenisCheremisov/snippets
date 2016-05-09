#include <list>
#include <string>
#include <memory>

#include <gtest/gtest.h>

#include <logtools/main.hh>
#include <ioutils/writer.hh>
#include <ioutils/bufio.hh>
#include <controller.hh>
#include <scanner.hh>

#include <iostream>


class DummyMainReader: public main::LineReader {
    std::list<std::string> buffer;

public:
    DummyMainReader& operator<< (std::string &&data) {
        buffer.push_back(data);
        return *this;
    }

    int get_lines(main::Visitor *ref, time_t stamp) {
        if (buffer.size() == 0) {
            return -1;
        }
        std::string data = buffer.front();
        ref->append_line(data.c_str(), data.length());
        buffer.pop_front();
        return 1;
    }
};


class DummyMeta: public logtools::LogMeta {
    int counter = -1;

public:
    bool next() {
        counter += 1;
        return counter <= 3;
    }

    std::string name() {
        return "test/data/file" + std::to_string(counter);
    }

    time_t line_stamp(const char *line, int len) {
        return 0;
    }
};


TEST(MainTests, Visitor) {
    io::StrWriter writer;
    main::Visitor visitor(&writer, 2);

    DummyMainReader r1, r2;
    r1 << "1";
    r2 << "2" << "output2";
    visitor.add_reader(&r1);
    visitor.add_reader(&r2);

    ASSERT_EQ(visitor.visit_all(0), 2);
    ASSERT_EQ(writer.buf, "1\n");

    ASSERT_EQ(visitor.visit_all(0), 1);
    ASSERT_EQ(writer.buf, "1\n2\n");
    visitor.flush();

    ASSERT_EQ(visitor.visit_all(0), 0);
    ASSERT_EQ(writer.buf, "1\n2\noutput2\n");
    ASSERT_EQ(visitor.visit_all(0), -1);
    ASSERT_EQ(writer.buf, "1\n2\noutput2\n");
}


TEST(MainTests, ContiniousReader) {
    io::StrWriter writer;
    DummyMeta meta;
    main::ContiniousFileReader reader(&meta);
    bufio::Reader br(&reader, 512*1024);
    bufio::Buf line;

    int i = 0;
    while (br.readline(&line) == bufio::_OK) {
        writer.write(line.buf, line.len);
        ASSERT_EQ(reader.file_count(), i++);
    }
    ASSERT_EQ(writer.buf, "file0file1file2 gzippedfile3");
}


class DummyOddSkipper: public logtools::LogMeta {
    int counter = 0;

public:
    bool next() {
        return true;
    }

    std::string name() {
        return "LOL";
    }

    time_t line_stamp(const char *line, int len) {
        time_t res = ((counter % 2) != 1) ? counter / 2 : -1;
        counter += 1;
        return res;
    }
};


TEST(MainTests, CareLogReader) {
    io::StrReader rawReader(
        "line 1\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7");
    bufio::Reader reader(&rawReader, 512*1024);
    DummyOddSkipper meta;

    main::LogScanner scanner(&meta, &reader, 2);
    bufio::Buf dst;
    const char *samples[] = {
            "line 1\nline 2", "line 3\nline 4", "line 5\nline 6", "line 7"
    };

    int counter = 0;
    while (scanner.readline(&dst, counter) != main::LogScanner::END) {
        ASSERT_EQ(std::string(dst.buf, dst.len), samples[counter]);
        counter += 1;
    }
    ASSERT_EQ(counter, 4);
}


TEST(MainTests, CareLogReaderWait) {
    io::StrReader rawReader(
        "line 1\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7");
    bufio::Reader reader(&rawReader, 512*1024);
    DummyOddSkipper meta;

    meta.line_stamp(NULL, 0);
    meta.line_stamp(NULL, 0);
    main::LogScanner scanner(&meta, &reader, 2);
    bufio::Buf dst;
    const char *samples[] = {
            "line 1\nline 2", "line 3\nline 4", "line 5\nline 6", "line 7"
    };

    int counter = 0;
    int abs_counter = 0;
    main::LogScanner::scanner_read_t res;
    while (res = scanner.readline(&dst, abs_counter),
           res != main::LogScanner::END) {
        if (res != main::LogScanner::EMPTY) {
            ASSERT_EQ(std::string(dst.buf, dst.len), samples[counter]);
            counter += 1;
        }
        abs_counter += 1;
    }
    ASSERT_EQ(abs_counter, 5);
}


class DummyFilter: public main::Filter {
    bool current;

public:
    DummyFilter(bool cur): current(cur) {}

    bool match(const char *line, uint64_t len) {
        current = !current;
        return current;
    }
};


class SuperDummyMeta: public logtools::LogMeta {
    int counter;

public:
    SuperDummyMeta(int c): counter(c) {}

    bool next() {
        return true;
    }
    std::string name() {
        return "";
    }
    time_t line_stamp(const char *line, int len) {
        return counter++;
    }
};


class ReaderFactory {
    std::shared_ptr<io::Reader> rawReader;
    std::shared_ptr<bufio::Reader> reader;
    std::shared_ptr<logtools::LogMeta> meta;
    std::shared_ptr<main::LogScanner> scanner;
    std::shared_ptr<main::Filter> filter;

public:
    ReaderFactory(io::Reader *r, main::Filter *f) {
        rawReader = std::shared_ptr<io::Reader>(r);
        reader = std::shared_ptr<bufio::Reader>(
            new bufio::Reader(r, 512*1024));
        meta = std::shared_ptr<logtools::LogMeta>(
            new DummyOddSkipper());
        scanner = std::shared_ptr<main::LogScanner>(
            new main::LogScanner(meta.get(), reader.get(), 1024));
        filter = std::shared_ptr<main::Filter>(f);
    }

    std::shared_ptr<main::FilteringLineReader> get() {
        return std::shared_ptr<main::FilteringLineReader>(
            new main::FilteringLineReader(scanner.get(), filter.get()));
    }
};



TEST(MainTests, FilteringLineReaderTest) {
    ReaderFactory f1(new io::StrReader("f1\nf2\nf3\nf4\nf5\nf6\nf7"),
                   new DummyFilter(false));
    ReaderFactory f2(new io::StrReader("s1\ns2\ns3\ns4\ns5\ns6\ns7"),
                   new DummyFilter(true));
    auto r1 = f1.get();
    auto r2 = f2.get();

    io::StrWriter writer;
    main::Visitor visitor(&writer, 1024);

    visitor.add_reader(r1.get());
    visitor.add_reader(r2.get());

    ASSERT_EQ(visitor.visit_all(0), 2);
    visitor.flush();
    ASSERT_EQ(writer.buf, "f1\nf2\n");

    ASSERT_EQ(visitor.visit_all(1), 2);
    visitor.flush();
    ASSERT_EQ(writer.buf, "f1\nf2\ns3\ns4\n");

    ASSERT_EQ(visitor.visit_all(2), 2);
    visitor.flush();
    ASSERT_EQ(writer.buf, "f1\nf2\ns3\ns4\nf5\nf6\n");

    ASSERT_EQ(visitor.visit_all(3), 0);
    visitor.flush();
    ASSERT_EQ(writer.buf, "f1\nf2\ns3\ns4\nf5\nf6\ns7\n");
}
