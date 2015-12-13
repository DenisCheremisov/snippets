#include <cstring>
#include <cstdlib>
#include <cstdio>

#include "matcher.hpp"
#include "bufio.hpp"


#define BUF_LENGTH 524288 /* 512Kb */



int main(int argc, char *argv[]) {
    Line line;

    SearchChar tab('\t');
    SearchString delim("apropasio");
    End end;
    TakeCharsUntil<SearchChar> column1(&tab);
    TakeCharsUntil<SearchChar> column2(&tab);
    PassCharsUntil<SearchChar> pass(&tab);
    PassCharsUntil<End> last(&end);

    Matcher matcher;
    matcher << column1 << pass << pass << column2 << last;
    const char *rest;

    Writer<BUF_LENGTH> dst(STDOUT_FILENO);
    Writer<BUF_LENGTH> err(STDERR_FILENO);
    Reader<BUF_LENGTH> reader(STDIN_FILENO);

    int res;
    int i = 0;
    while( (res = reader.getline(line)) > 0 ) {
        bool res = matcher.feed(line, rest);
        if (res) {
            dst.append(column1.chars());
            dst.append("\t", 1);
            dst.append(column2.chars());
            dst.append("\n", 1);
        } else {
            err.append("ERROR: unparsed ", 16);
            err.append(line);
            err.append("\n", 1);
            err.flush();
        }
    }
    if (res < 0) {
        perror("");
        return -1;
    }
}
