#pragma once

#include <iostream>

#include <unistd.h>

#include <logtools/main.hh>
#include <timeutils/seq.hh>


const int BUF_LENGTH = 1024;

namespace logtools {
    class ErrLogMeta: public LogMeta {
        nstime::Sequence *seq;
        int prev_day;
        time_t prev_day_tstamp;
        std::string path;

    public:
        ErrLogMeta(std::string _path, nstime::Sequence *_seq) {
            prev_day = -1;
            path = _path;
            seq = _seq;
        }
        ~ErrLogMeta() {};

        bool next() {
            return seq->next();
        }

        std::string name() {
            char buf[BUF_LENGTH];
            tm tm_val = seq->time();
            size_t len = strftime(buf, BUF_LENGTH - 1, path.c_str(), &tm_val);
            return std::string(buf, len);
        }

        time_t line_stamp(const char *line, int len) {
            if (len < 10 || line[9] != ':') {
                return -1;
            }

            int day = (line[0] - '0')*10 + (line[1] - '0');
            if (day != prev_day) {
                tm cur_day = seq->time();
                time_t cur_tstamp = mktime(&cur_day);
                cur_tstamp -= 3600*24;
                tm tmp1;
                localtime_r(&cur_tstamp, &tmp1);
                cur_tstamp += 3600*24*2;
                tm tmp2;
                localtime_r(&cur_tstamp, &tmp2);
                if (tmp1.tm_mday == day) {
                    cur_day = tmp1;
                } else if (tmp2.tm_mday == day) {
                    cur_day = tmp2;
                } else if (cur_day.tm_mday != day) {
                    std::cerr << std::string(line, len) << std::endl;
                    return -2;
                }
                cur_day.tm_hour = cur_day.tm_min = cur_day.tm_sec = 0;
                cur_day.tm_isdst = 0;

                prev_day_tstamp = mktime(&cur_day);
                prev_day = cur_day.tm_mday;
            }

            return prev_day_tstamp +
                ((line[3] - '0')*10 + line[4] - '0')*3600 +
                ((line[5] - '0')*10 + line[6] - '0')*60 +
                ((line[7] - '0')*10 + line[8] - '0');
        }
    };
}
