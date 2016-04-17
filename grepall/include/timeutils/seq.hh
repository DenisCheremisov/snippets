#pragma once

#include <ctime>

namespace nstime {

    class Sequence {
    public:
        virtual bool next() = 0;
        virtual tm time() = 0;
    };
}
