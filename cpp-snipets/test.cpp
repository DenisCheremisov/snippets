#include <cstring>
#include <cstdlib>
#include <cstdio>

#include "matcher.hpp"


int main(int argc, char *argv[]) {
    const char *content = "abcdef3033";
    Line line(content, strlen(content));
    const char *data[] = {"2033", "33", "19"};

    SearchVariant variant(3, data);

    TakeCharsUntil<SearchVariant> take(&variant);

    Matcher matcher;
    matcher << take;
    const char *rest;

    bool res = matcher.feed(line, rest);
    if (!res) {
        printf("Failed\n");
    } else {
        for(int i = 0; i < take.chars().len(); i++) {
            putchar(take.chars().data()[i]);
        }
        putchar('\n');
    }
}
