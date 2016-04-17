#include <string>

#include <gtest/gtest.h>

#include <timeutils/hour_seq.hh>

TEST(TimeUtilsTests, HourSequencerTest) {
    const char *format = "%Y-%m-%dT%H:%M:%S";
    tm from, to;
    strptime("2016-04-03T10:12:13", format, &from);
    strptime("2016-04-04T00:00:00", format, &to);

    tm data[14];
    nstime::HourSequence test(from, to);
    int i = 0;
    while (test.next()) {
        data[i] = test.time();
        i++;
    }
    ASSERT_EQ(i, 14);
    ASSERT_EQ(test.next(), false);

    char buf[500];
    strftime(buf, 499, format, &data[0]);
    ASSERT_EQ(std::string(buf), std::string("2016-04-03T10:00:00"));
    strftime(buf, 499, format, &data[13]);
    ASSERT_EQ(std::string(buf), std::string("2016-04-03T23:00:00"));
    for (int i = 1; i < 14; i++) {
        ASSERT_EQ(mktime(&data[i - 1]) + 3600, mktime(&data[i]));
    }
}
