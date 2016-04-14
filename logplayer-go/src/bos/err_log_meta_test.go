package bos

import (
	"testing"
	"time"
	"utils"

	"github.com/pbnjay/strptime"
	"github.com/stretchr/testify/assert"
)

func TestErrLogMeta(t *testing.T) {
	day, _ := strptime.Parse("2016-04-03", "%Y-%m-%d")
	seq := utils.NewTimeHourSequencer(day, day.Add(time.Second))
	meta := NewErrLogMeta("/1/2/3/bos_srv-k01a.err.%Y%m%d%H", seq)
	if !assert.True(t, meta.Next()) {
		return
	}
	if !assert.Equal(t, meta.Name(), "/1/2/3/bos_srv-k01a.err.2016040300") {
		return
	}

	dtime, _ := strptime.Parse("2016-04-02T23:59:42", "%Y-%m-%dT%H:%M:%S")
	tstamp, ok := meta.LineTimestamp([]byte("02/235942:"))
	if !assert.True(t, ok) {
		return
	}
	if !assert.Equal(t, dtime.UnixNano(), tstamp) {
		return
	}
	tstamp, ok = meta.LineTimestamp([]byte("03/000001: ahahah"))
	if !assert.True(t, ok) {
		return
	}
	if !assert.Equal(t, dtime.UnixNano()+19*int64(time.Second), tstamp) {
		return
	}
}
