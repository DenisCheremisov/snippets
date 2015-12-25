package arr

import "bytes"

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

type SearchAnyByte string

func (seeker SearchAnyByte) Seek(source []byte) (bool, int, []byte) {
	pos := bytes.IndexAny(source, string(seeker))
	if pos < 0 {
		return false, 0, source
	} else {
		return true, pos, source[pos+1:]
	}
}

type SearchAnyOf [][]byte

func (seeker SearchAnyOf) Seek(source []byte) (bool, int, []byte) {
	for _, item := range seeker {
		pos := bytes.Index(source, item)
		if pos < 0 {
			continue
		} else {
			return true, pos, source[pos+len(item):]
		}
	}
	return false, 0, source
}

type SearchClosest [][]byte

func (seeker SearchClosest) Seek(source []byte) (bool, int, []byte) {
	pos := len(source)
	var rest []byte
	success := false

	for _, item := range seeker {
		_pos := bytes.Index(source, item)
		if _pos >= 0 && _pos < pos {
			success = true
			pos = _pos
			rest = source[pos+len(item):]
		}
	}

	return success, pos, rest
}
