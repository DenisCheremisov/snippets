#include <cstring>
#include <cstdlib>
#include <cstdio>

#include "bufio.hpp"
#include "line.hpp"

#define BUF_LENGTH 524288


int main(int argc, char *argv[]) {
    Writer<BUF_LENGTH> stdout(STDOUT_FILENO);

    char *sep = (char*) "\t\t\t\t\n";
    char buf[4096];

    for (int i = 0; i < 50000000; i++) {
        for(int k = 0; k < 5; k++) {
            int len = (random() % 80) + 1;
            for (int l = 0; l < len; l++) {
                buf[l] = (random() % 10) + '0';
            }
            buf[len] = sep[k];
            stdout.write(buf, len + 1);
        }
    }
}
