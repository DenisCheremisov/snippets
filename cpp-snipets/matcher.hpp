#ifndef _MATCHER_HPP_INCLUDED_UQ4938204830498304982390_
#define _MATCHER_HPP_INCLUDED_UQ4938204830498304982390_

#include <vector>

#include "line.hpp"

class Seeker {
public:
    virtual bool seek(const Line &, const char *&prepos, const char *&rest) = 0;
};

class Taker {
public:
    virtual bool take(const Line &, const char *&rest) = 0;
};


class Matcher {
    std::vector<Taker*> takers;

public:
    bool feed(const Line &line, const char *&rest) {
        rest = line.data();
        int i = 0;
        int len = line.len();
        for(auto it: this->takers) {
            const char *next;
            bool res = it->take(Line(rest, len), next);
            len -= next - rest;
            rest = next;
            if(!res) {
                return false;
            }
            i++;
        }
        return true;
    }

    Matcher& operator<< (Taker &item) {
        this->takers.push_back(&item);
        return *this;
    }

    Matcher& operator<< (Taker &&item) {
        this->takers.push_back(&item);
        return *this;
    }
};


template <class S>
class PassCharsUntil: public Taker {
    S *seeker;

public:
    PassCharsUntil(S *s): seeker(s) {}

    bool take(const Line &line, const char *&rest) {
        const char *prepos;
        const char *next;
        if (!this->seeker->seek(line, prepos, next)) {
            rest = line.data();
            return false;
        }
        rest = next;
        return true;
    }
};


template <class S>
class TakeCharsUntil: public Taker {
private:
    Line data;
    S *seeker;

public:
    TakeCharsUntil(S *s): seeker(s) {}

    bool take(const Line &line, const char *&rest) {
        const char *next;
        const char *prepos;
        if(!this->seeker->seek(line, prepos, next)) {
            rest = line.data();
            return false;
        }
        rest = next;
        this->data = Line(line.data(), prepos);
        return true;
    }

    const Line &chars() const {
        return data;
    }
};


class SearchString  {
    const char *string;
    int length;

public:
    SearchString(const char *data): string(data), length(strlen(string)) {}

    bool seek(const Line &line, const char *&prepose, const char *&rest) {
        prepose = (char*) memmem(line.data(), line.len(), this->string, length);
        if (prepose == NULL) {
            rest = line.data();
            return false;
        } else {
            rest = prepose + length;
            return true;
        }
    }
};


class SearchChar {
    char chr;

public:
    SearchChar(char c): chr(c) {}

    bool seek(const Line &line, const char *&prepos, const char *&rest) {
        prepos = (char*) memchr(line.data(), chr, line.len());
        if (prepos == NULL) {
            rest = line.data();
            return false;
        } else {
            rest = prepos + 1;
            return true;
        }
    }
};

class End {
public:

    bool seek(const Line &line, const char *&prepos, const char *&rest) {
        prepos = rest = line.data() + line.len();
        return true;
    }
};

#endif // _MATCHER_HPP_INCLUDED_UQ4938204830498304982390_
