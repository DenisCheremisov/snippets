#include <list>
#include <string>

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
