package output

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogSleeper(t *testing.T) {
	s := NewLogSleeper(1)
	tm := time.Now().UnixNano()

	if !assert.True(t, s.Delta(tm) <= 0) {
		return
	}
	tm += int64(time.Second)
	if !assert.Equal(t, s.Delta(tm), int64(time.Second)) {
		return
	}
	tm += int64(time.Second)
	if !assert.Equal(t, s.Delta(tm), int64(time.Second)) {
		return
	}
	tm += int64(time.Second)
	s.Sync()
	d := s.Delta(tm)
	if !assert.True(t, d <= 3*int64(time.Second)) {
		return
	}

	if !assert.True(t, d > 29999*int64(time.Second)/10000) {
		return
	}
}
