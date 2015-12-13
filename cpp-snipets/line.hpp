#ifndef _STRING_HPP_INCLUDED_120498230SDIFCJKJHS088123adsfk_123_
#define _STRING_HPP_INCLUDED_120498230SDIFCJKJHS088123adsfk_123_

class Line {
    const char *buf;
    const char *end;

public:
    Line(): buf(NULL), end(NULL) {}
    Line(const char *s, const char *f): buf(s), end(f) {}
    Line(const char *s, int len): buf(s), end(s + len) {}
    Line(const Line &line): buf(line.buf), end(line.end) {}

    inline int len() const {
        return end - buf;
    }

    inline const char *data() const {
        return buf;
    }
};


#endif // _STRING_HPP_INCLUDED_120498230SDIFCJKJHS088123adsfk_123_
