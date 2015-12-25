package arr

import "bytes"

// Seeker is an interface for seek objects
type Seeker interface {
	// Seek returns true, pos, rest as success
	// where pos is the position in the source where the match is from
	// rest is the part of source started after the end of the match
	//
	// Seek  returns false, <undef>, source on failure
	Seek([]byte) (bool, int, []byte)
}

type Taker interface {
	Take([]byte) (bool, []byte)
}

type Matcher []Taker

func NewMatcher(takers ...Taker) Matcher {
	return Matcher(takers)
}

func (matcher Matcher) Feed(source []byte) (bool, []byte) {
	var success bool
	for _, item := range matcher {
		if success, source = item.Take(source); !success {
			return false, source
		}
	}
	return true, source
}

type PassBytes struct {
	seeker Seeker
}

func PassBytesUntil(seeker Seeker) *PassBytes {
	return &PassBytes{seeker}
}

func (taker *PassBytes) Take(source []byte) (bool, []byte) {
	success, _, rest := taker.seeker.Seek(source)
	if !success {
		return false, source
	} else {
		return true, rest
	}
}

type TakeBytes struct {
	seeker Seeker
	val    []byte
}

func TakeBytesUntil(seeker Seeker) *TakeBytes {
	return &TakeBytes{seeker, nil}
}

func (taker *TakeBytes) Take(source []byte) (bool, []byte) {
	sucess, pos, rest := taker.seeker.Seek(source)
	if !sucess {
		return false, source
	} else {
		taker.val = source[:pos]
		return true, rest
	}
}

func (taker *TakeBytes) Bytes() []byte {
	return taker.val
}

type OnTheEnd bool

func (taker OnTheEnd) Take(source []byte) (bool, []byte) {
	if bool(taker) {
		return len(source) == 0, source
	} else {
		return len(source) != 0, source
	}
}

type PassPattern []byte

func (taker PassPattern) Take(source []byte) (bool, []byte) {
	if len(source) < len(taker) {
		return false, source
	}
	if bytes.Equal(taker, source[:len(taker)]) {
		return true, source[len(taker):]
	} else {
		return false, source
	}
}

type PassFuncCheck func(byte) bool

func (taker PassFuncCheck) Take(source []byte) (bool, []byte) {
	if len(source) > 1 {
		if taker(source[0]) {
			return true, source[1:]
		}
	}
	return false, source
}
