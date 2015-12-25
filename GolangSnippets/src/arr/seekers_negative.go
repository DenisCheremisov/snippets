package arr

type SearchUntilByte byte

func (seeker SearchUntilByte) Seek(source []byte) (bool, int, []byte) {
	for i, char := range source {
		if char != byte(seeker) {
			return true, i, source[i:]
		}
	}
	if len(source) > 0 {
		if source[len(source)-1] == byte(seeker) {
			return true, len(source), source[len(source):]
		}
	}
	return true, 0, source
}
