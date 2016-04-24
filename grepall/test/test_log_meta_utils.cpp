#include <string>

#include <gtest/gtest.h>

#include <logtools/bos_err.hh>
#include <timeutils/hour_seq.hh>

TEST(LogToolsTests, BossErrLogMetaToolTest) {
    tm from, to;
    strptime("2016-04-03T00:00:00", "%Y-%m-%dT%H:%M:%S", &from);
    from.tm_isdst = 0;
    time_t ahah = mktime(&from) + 1;
    localtime_r(&ahah, &to);
    to.tm_isdst = 0;

    nstime::HourSequence seq(from, to);
    logtools::ErrLogMeta meta("/1/2/3/bos_srv-k01a.err.%Y%m%d%H", &seq);
    ASSERT_EQ(meta.next(), true);
    ASSERT_EQ(meta.name(), "/1/2/3/bos_srv-k01a.err.2016040300");

    tm tmp;
    strptime("2016-04-02T23:59:42", "%Y-%m-%dT%H:%M:%S", &tmp);
    tmp.tm_isdst = 0;
    time_t tstamp = meta.line_stamp("02/235942:", 10);
    ASSERT_EQ(tstamp, mktime(&tmp));

    tstamp = meta.line_stamp("03/000001: ahahah", 17);
    ASSERT_EQ(tstamp, mktime(&tmp) + 19);

    tstamp = meta.line_stamp("04/000101: ahahah", 17);
    ASSERT_EQ(tstamp, mktime(&tmp) + 79 + 3600*24);
}
