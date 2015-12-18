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

class SearchAnyOf {
    const int len;
    const char **anyof;
    int *lengths;

public:
    SearchAnyOf(int l, const char **ao): len(l), anyof(ao) {
        lengths = new int[l];
        for(int i = 0; i < l; i++) {
            lengths[i] = strlen(ao[i]);
        }
    }
    virtual ~SearchAnyOf() {
        delete lengths;
    }

    bool seek(const Line &line, const char *&prepos, const char *&rest) {
        for(int i = 0; i < len; i++) {
            const char *item = anyof[i];
            prepos = (char*) memmem(
                line.data(), line.len(), anyof[i], lengths[i]);
            if (prepos != NULL) {
                rest = prepos + lengths[i];
                return true;
            }
        }
        return false;
    }
};


class SearchClosest {
    const int len;
    const char **variants;
    int *lengths;

public:
    SearchClosest(int l, const char **v): len(l), variants(v) {
        lengths = new int[l];
        for(int i = 0; i < len; i++) {
            lengths[i] = strlen(variants[i]);
        }
    }
    virtual ~SearchClosest() {
        delete lengths;
    }

    bool seek(const Line &line, const char *&prepos, const char *&rest) {
        const char *_prepos, *_rest;
        bool success = false;
        prepos = line.data() + line.len();
        for(int i = 0; i < len; i++) {
            const char *item = variants[i];
            _prepos = (char*) memmem(
                line.data(), line.len(), variants[i], lengths[i]);
            if (_prepos != NULL) {
                success = true;
                if(_prepos < prepos) {
                    prepos = _prepos;
                    rest = prepos + lengths[i];
                }
            }
        }
        return success;
    }
};

#endif // _MATCHER_HPP_INCLUDED_UQ4938204830498304982390_
