#include <string.h>
#include <unistd.h>
#include <stdlib.h>
#include <memory.h>
#include <stdio.h>

#define BUF_LENGTH 524288 /* 512Kb */


typedef struct {
    char buf[BUF_LENGTH];
    int current;
    int out_fd;
} writer;


writer *writer_create(int fd) {
    writer *result = (writer*)malloc(sizeof(writer));
    if(result != NULL) {
        result->current = 0;
        result->out_fd = fd;
        return result;
    } else {
        perror("");
        exit(-1);
    }
}


void writer_flush(writer *dst) {
    if(dst->current == 0) {
        return;
    }
    int result = write(dst->out_fd, dst->buf, dst->current);
    if(result < 0) {
        perror("");
        exit(-1);
    }
    dst->current = 0;
}


void writer_append(writer *dst, const char *start, int size) {
    int space_left = BUF_LENGTH - dst->current;
    if(space_left < size) {
        writer_flush(dst);
    }
    memcpy(dst->buf + dst->current, start, size);
    dst->current += size;
}


void writer_add(writer *dst, writer *add) {
    if((dst->current + add->current) > BUF_LENGTH) {
        writer_flush(dst);
    }
    memcpy(dst->buf + dst->current, add->buf, add->current);
    dst->current += add->current;
}


void extract_data(const char *start, const char *end, writer *dst, writer *tmp) {
    tmp->current = 0;

    char *ptr = strstr(start, "2015");
    char *endptr;
    if(ptr == NULL) {
        return;
    }
    if((end - ptr) < 40) {
        return;
    }
    if(ptr[8] != '.') {
        return;
    }

    writer_append(tmp, ptr, 8);
    ptr = strchr(ptr + 8, ':');
    if(ptr == NULL) {
        return;
    }
    ptr += 1;
    if((end - ptr) < 40) {
        return;
    }
    writer_append(tmp, "\t", 1);
    writer_append(tmp, ptr, 9);

    ptr = strstr(ptr + 40, "user=");
    if(ptr == NULL) {
        return;
    }
    ptr += 5;
    endptr = strchr(ptr, ' ');
    if(endptr == NULL) {
        return;
    }
    writer_append(tmp, "\t", 1);
    writer_append(tmp, ptr, endptr - ptr);

    ptr = strstr(endptr, "attached=");
    if(ptr == NULL) {
        return;
    }
    ptr += 9;
    endptr = strchr(ptr, ' ');
    if(endptr == NULL) {
        return;
    }
    writer_append(tmp, "\t", 1);
    writer_append(tmp, ptr, endptr - ptr);

    ptr = strstr(ptr, "msisdn=");
    if(ptr == NULL) {
        return;
    }
    ptr += 7;
    endptr = strchr(ptr, ' ');
    if(endptr == NULL) {
        return;
    }
    writer_append(tmp, "\t", 1);
    writer_append(tmp, ptr, endptr - ptr);
    writer_append(tmp, "\n", 1);

    writer_add(dst, tmp);
}


typedef struct {
    char buf[BUF_LENGTH + 1];
    int current;
    int capacity;
    int fd;
    int final;
} reader;


reader* reader_create(int fd) {
    reader *result = malloc(sizeof(reader));
    if(result != NULL) {
        result->current = 0;
        result->capacity = 0;
        result->fd = fd;
        result->final = 0;
        result->buf[0] = 0;
        return result;
    } else {
        perror("");
        exit(-1);
    }
}


typedef struct {
    char *start;
    char *end;
} line_type;


line_type* reader_getline(reader *src, line_type *res) {
    char *ptr = strchr(src->buf + src->current, '\n');
    if(ptr == NULL) {
        if(src->final) {
            return NULL;
        }
        if(src->current > 0) {
            memcpy(src->buf,
                   src->buf + src->current, src->capacity - src->current);
            src->capacity -= src->current;
            src->current = 0;
        }
        int size = read(src->fd,
                        src->buf + src->capacity, BUF_LENGTH - src->capacity);
        if(size < 0) {
            perror("");
            exit(-1);
        }
        src->final = size == 0;
        src->capacity += size;
        src->buf[src->capacity] = 0;
        ptr = strchr(src->buf, '\n');
        if(ptr != NULL) {
            ptr[0] = 0;
        }
    }
    res->start = src->buf + src->current;
    if(ptr != NULL) {
        res->end = ptr;
        src->current = ptr - src->buf + 1;
    } else {
        res->end = src->buf + src->capacity;
        src->current = src->capacity;
    }
    return res;
}



int main(int argc, char *argv[]) {
    writer *dst = writer_create(STDOUT_FILENO);
    reader *src = reader_create(STDIN_FILENO);

    char tab = '\t';
    char nl = '\n';

    line_type line;

    int res;


    while ( reader_getline(src, &line) ) {

        char *rest = line.start;
        int restlen = line.end - line.start;

        // 1st column
        char *tab_pos = memchr(rest, tab, restlen);
        if(tab_pos == NULL) {
            continue;
        }
        writer_append(dst, rest, tab_pos - rest);
        writer_append(dst, "\t", 1);
        restlen -= tab_pos - rest + 1;
        rest = tab_pos + 1;

        // 2nd column
        tab_pos = memchr(rest, tab, restlen);
        if(tab_pos == NULL) {
            continue;
        }
        restlen -= tab_pos - rest + 1;
        rest = tab_pos + 1;

        // 3nd column
        tab_pos = memchr(rest, tab, restlen);
        if(tab_pos == NULL) {
            continue;
        }
        restlen -= tab_pos - rest + 1;
        rest = tab_pos + 1;

        // 4th column
        tab_pos = memchr(rest, tab, restlen);
        if(tab_pos == NULL) {
            continue;
        }

        writer_append(dst, rest, tab_pos - rest);
        writer_append(dst, "\n", 1);
        restlen -= tab_pos - rest + 1;
        rest += tab_pos - rest + 1;
    }

    writer_flush(dst);
}
