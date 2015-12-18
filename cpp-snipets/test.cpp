#include <string>

#include <gtest/gtest.h>

#include "matcher.hpp"

TEST(VariantTests, HandleAnyOf) {
    const char *content = "abcdef3033";
    Line line(content, strlen(content));
    const char *data[] = {"2033", "33", "19"};

    SearchAnyOf variant(3, data);

    TakeCharsUntil<SearchAnyOf> take(&variant);

    Matcher matcher;
    matcher << take;
    const char *rest;

    bool res = matcher.feed(line, rest);
    ASSERT_TRUE(res);
    ASSERT_EQ(strlen(rest), 0);
}


TEST(VariantTests, HandleClosest) {
    const char *content = "abcdefgh 9999";
    Line line(content, strlen(content));
    const char *data[] = {"1", "2", "99"};

    SearchClosest variant(3, data);
    TakeCharsUntil<SearchClosest> take(&variant);

    Matcher matcher;
    matcher << take;
    const char *rest;

    bool res = matcher.feed(line, rest);
    ASSERT_TRUE(res);
    ASSERT_EQ(std::string(take.chars().data(), take.chars().len()),
              std::string("abcdefgh "));
    ASSERT_EQ(std::string(rest), std::string("99"));

    const char *content2 = "abc945a193";
    line = Line(content2, strlen(content2));
    SearchClosest v2(3, data);
    TakeCharsUntil<SearchClosest> take2(&v2);
    Matcher m2;
    m2 << take2;
    res = m2.feed(line, rest);
    ASSERT_TRUE(res);
    ASSERT_EQ(std::string(take2.chars().data(), take2.chars().len()),
              std::string("abc945a"));
    ASSERT_EQ(std::string(rest), std::string("93"));
}


int main(int argc, char **argv) {
  testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
