package dateseq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateSeq(t *testing.T) {
	seq, err := NewDateSeq(
		time.Date(2015, 11, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2015, 12, 6, 0, 0, 0, 0, time.UTC))
	assert.Nil(t, err)

	day_seq := []int{
		14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25,
		26, 27, 28, 29, 30, 1, 2, 3, 4, 5, 6}

	i := 0
	for seq.Next() {
		day := seq.Date()
		assert.Equal(t, day.Day(), day_seq[i])
		if day.Day() < 14 {
			assert.True(t, day.Month() == 12)
		} else {
			assert.True(t, day.Month() == 11)
		}
		i += 1
	}
}
