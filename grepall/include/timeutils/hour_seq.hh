#include <timeutils/seq.hh>

namespace nstime {
    class HourSequence: public Sequence {
    private:
        time_t from, to, cur;

    public:
        HourSequence(tm _from, tm _to) {
            _from.tm_min = 0;
            _from.tm_sec = 0;
            _from.tm_isdst = 0;
            from = mktime(&_from);
            _to.tm_isdst =  0;
            to = mktime(&_to);
            cur = from - 3600;
        }

        bool next() {
            cur = cur + 3600;
            return cur < to;
        }

        tm time() {
            tm res;
            localtime_r(&cur, &res);
            return res;
        }
    };
}
