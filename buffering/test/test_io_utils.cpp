#include <algorithm>
#include <memory>
#include <cstring>
#include <string>

#include <gtest/gtest.h>

#include <ioutils/reader.hh>
#include <ioutils/bufio.hh>


TEST(IOUtilsTests, StrReader) {
    std::unique_ptr<io::Reader> reader(
        new io::StrReader("123456789abcdefghijklmnopqrstuvwxyz"));
    const char *samples[] = {
        "1", "23", "456", "789a", "bcdef", "ghijkl",
        "mnopqrs", "tuvwxyz"
    };
    char buf[500];
    memset(buf, 0, 500);
    int64_t m;
    for(int i = 0; true; i++) {
        m = reader->read(buf, i + 1);
        if (i == 8) {
            break;
        } else {
            ASSERT_EQ(m, std::min(i + 1, 7));
            ASSERT_EQ(m, strlen(samples[i]));
            ASSERT_EQ(memcmp(samples[i], buf, i + 1), 0);
        }
    }
}


TEST(IOUtilsTests, BufReaderTestScanning) {
    std::unique_ptr<io::Reader> reader(
        new io::StrReader("1\n123\n1234\n12345"));
    const char *samples[] = {
        "1", "123", "1234", "12345"
    };
    bufio::read_result_t expected_status[] = {
        bufio::_OK, bufio::_OK, bufio::_OK, bufio::_OK, bufio::_EOF
    };
    bufio::Reader br(reader.get(), 512*1024);
    for (int i = 0; i < 5; i++) {
        bufio::Buf buf;
        auto res = br.readline(&buf);
        ASSERT_EQ(res, expected_status[i]);
        if (i < 4) {
            ASSERT_EQ(std::string(buf.buf, buf.len), std::string(samples[i]));
        }
    }
}


TEST(IOUtilsTests, BufReaderTestChunking) {
    std::unique_ptr<io::Reader> reader(
        new io::StrReader("1\n123\n1234\n12345"));
    const char *samples[] = {
        "1", "123", "1234", "12345"
    };
    bufio::read_result_t expected_status[] = {
        bufio::_OK, bufio::_OK, bufio::_OK, bufio::_OK, bufio::_EOF
    };
    bufio::Reader br(reader.get(), 2);
    for (int i = 0; i < 5; i++) {
        bufio::Buf buf;
        auto res = br.readline(&buf);
        ASSERT_EQ(res, expected_status[i]);
        if (i < 4) {
            ASSERT_EQ(std::string(buf.buf, buf.len), std::string(samples[i]));
        }
    }

    std::unique_ptr<io::Reader> reader2(
        new io::StrReader("1\n123\n1234\n12345"));
    bufio::Reader br2(reader2.get(), 1);
    for (int i = 0; i < 5; i++) {
        bufio::Buf buf;
        auto res = br2.readline(&buf);
        ASSERT_EQ(res, expected_status[i]);
        if (i < 4) {
            ASSERT_EQ(std::string(buf.buf, buf.len), std::string(samples[i]));
        }
    }
}
