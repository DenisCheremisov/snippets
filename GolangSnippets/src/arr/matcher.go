package dateseq

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

type Search []byte

func SearchString(val string) Search {
	return Search([]byte(val))
}

func (seeker Search) Seek(source []byte) (bool, int, []byte) {
	pos := bytes.Index(source, seeker)
	if pos < 0 {
		return false, 0, source
	} else {
		return true, pos, source[pos+len(seeker):]
	}
}

type EndSeeker bool

func (seeker EndSeeker) Seek(source []byte) (bool, int, []byte) {
	return true, len(source), source[len(source):]
}

type SearchByte byte

func (seeker SearchByte) Seek(source []byte) (bool, int, []byte) {
	pos := bytes.IndexByte(source, byte(seeker))
	if pos < 0 {
		return false, 0, source
	} else {
		return true, pos, source[pos+1:]
	}
}
