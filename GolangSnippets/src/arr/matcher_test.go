package arr

import (
	"fmt"
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

func TestVariants(t *testing.T) {
	take := TakeBytesUntil(SearchAnyOf([][]byte{
		[]byte("1234"),
		[]byte("4321"),
		[]byte("9999"),
	}))
	matcher := NewMatcher(take)

	success, _ := matcher.Feed([]byte("abcdef1234"))
	assert.True(t, success)
	assert.Equal(t, "abcdef", string(take.Bytes()))

	success, _ = matcher.Feed([]byte("abcdef22344321"))
	assert.True(t, success)
	assert.Equal(t, "abcdef2234", string(take.Bytes()))

	success, rest := matcher.Feed([]byte("abcdef12999934431"))
	assert.True(t, success)
	assert.Equal(t, "abcdef12", string(take.Bytes()))
	assert.Equal(t, string(rest), "34431")

	take = TakeBytesUntil(SearchClosest([][]byte{
		[]byte("12"),
		[]byte("13"),
		[]byte("99"),
	}))
	matcher = NewMatcher(take)

	success, rest = matcher.Feed([]byte("abcdef 99 12"))
	assert.True(t, success)
	assert.Equal(t, "abcdef ", string(take.Bytes()))
	assert.Equal(t, " 12", string(rest))

	success, rest = matcher.Feed([]byte("abcdef 12 99"))
	assert.True(t, success)
	assert.Equal(t, "abcdef ", string(take.Bytes()))
	assert.Equal(t, " 99", string(rest))

	success, rest = matcher.Feed([]byte("abcdef 14 99"))
	assert.True(t, success)
	assert.Equal(t, "abcdef 14 ", string(take.Bytes()))
	assert.Equal(t, len(rest), 0)

	success, _ = matcher.Feed([]byte("abcedef  jsdlfkjk 914"))
	assert.False(t, success)
}

func TestNegative(t *testing.T) {
	matcher := NewMatcher(
		PassBytesUntil(SearchUntilByte(byte('0'))))
	ok, rest := matcher.Feed([]byte("00000000000000000111"))
	assert.True(t, ok)
	assert.Equal(t, string(rest), "111")
}

func TestOnTheEnd(t *testing.T) {
	matcher := NewMatcher(OnTheEnd(true))
	ok, _ := matcher.Feed([]byte(""))
	assert.True(t, ok)

	ok, _ = matcher.Feed([]byte("123"))
	assert.False(t, ok)

	matcher = NewMatcher(OnTheEnd(false))
	ok, _ = matcher.Feed([]byte(""))
	assert.False(t, ok)

	ok, _ = matcher.Feed([]byte("123"))
	assert.True(t, ok)
}

func TestReadyExpression(t *testing.T) {
	pass15 := PassFuncCheck(func(x byte) bool {
		return x == byte('1') || x == byte('5')
	})
	matcher := NewMatcher(
		PassBytesUntil(SearchByte('#')),
		pass15, pass15,
		PassPattern([]byte("1110")),
		PassBytesUntil(SearchUntilByte(byte('0'))),
		OnTheEnd(true))

	line := []byte(
		"2016-06-10 21:39:53.474667642 +0300 MSK -- log entry #15111000")
	res, rest := matcher.Feed(line)
	assert.True(t, res)
	fmt.Println(string(rest))
}
