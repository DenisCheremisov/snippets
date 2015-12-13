package dateseq

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcher(t *testing.T) {
	take1 := TakeBytesUntil(SearchString(" = indexOf "))
	matcher1 := NewMatcher(take1)

	success, rest := matcher1.Feed([]byte("123 = indexOf "))
	assert.True(t, success)
	assert.Equal(t, len(rest), 0)
	assert.Equal(t, "123", string(take1.Bytes()))

	wrong_sample := []byte("1232 = index_of")
	success, rest = matcher1.Feed(wrong_sample)
	assert.False(t, success)
	assert.Equal(t, len(rest), len(wrong_sample))

	take2 := TakeBytesUntil(SearchString(" end"))
	matcher2 := NewMatcher(take1, take2)
	success, rest = matcher2.Feed([]byte("ab.def = indexOf array end"))
	assert.True(t, success)
	assert.Equal(t, len(rest), 0)
	assert.Equal(t, "ab.def", string(take1.Bytes()))
	assert.Equal(t, "array", string(take2.Bytes()))

	success, rest = matcher2.Feed([]byte("1 = indexOf bugagaga"))
	assert.False(t, success)
	assert.Equal(t, string(rest), "bugagaga")

	take3 := TakeBytesUntil(EndSeeker(true))
	matcher3 := NewMatcher(PassBytesUntil(SearchString("\t")), take3)
	success, rest = matcher3.Feed([]byte("123213\tabcdef"))
	assert.True(t, success)
	assert.Equal(t, string(take3.Bytes()), "abcdef")
}

func BenchmarkMatcher(b *testing.B) {
	line := []byte("1\t2\t3\t4\t5")

	i := 0
	take1 := TakeBytesUntil(SearchByte('\t'))
	take2 := TakeBytesUntil(SearchByte('\t'))
	take3 := TakeBytesUntil(SearchByte('\t'))
	take4 := TakeBytesUntil(SearchByte('\t'))
	take5 := TakeBytesUntil(EndSeeker(true))
	matcher := NewMatcher(take1, take2, take3, take4, take5)
	for i < b.N {
		success, rest := matcher.Feed(line)
		if !(success && len(rest) == 0) {
			b.Error("Wrong data")
		}
		if take1.Bytes()[0] != '1' || len(take1.Bytes()) > 1 {
			b.Error("Wrong data")
		}
		if take2.Bytes()[0] != '2' || len(take2.Bytes()) > 1 {
			b.Error("Wrong data")
		}
		if take3.Bytes()[0] != '3' || len(take3.Bytes()) > 1 {
			b.Error("Wrong data")
		}
		if take4.Bytes()[0] != '4' || len(take4.Bytes()) > 1 {
			b.Error("Wrong data")
		}
		if take5.Bytes()[0] != '5' || len(take5.Bytes()) > 1 {
			b.Error("Wrong data")
		}
		i += 1
	}
}

func BenchmarkSplit(b *testing.B) {
	line := []byte("1\t2\t3\t4\t5")
	sep := []byte("\t")

	i := 0
	for i < b.N {
		res := bytes.Split(line, sep)
		if len(res) != 5 {
			b.Error("Wrong data")
		}
		if res[0][0] != '1' || len(res[0]) > 1 {
			b.Error("Wrong data")
		}
		if res[1][0] != '2' || len(res[1]) > 1 {
			b.Error("Wrong data")
		}
		if res[2][0] != '3' || len(res[2]) > 1 {
			b.Error("Wrong data")
		}
		if res[3][0] != '4' || len(res[3]) > 1 {
			b.Error("Wrong data")
		}
		if res[4][0] != '5' || len(res[4]) > 1 {
			b.Error("Wrong data")
		}
		i += 1
	}
}
