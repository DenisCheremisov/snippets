package bos

import (
	"testing"
	"time"
	"utils"

	"github.com/pbnjay/strptime"
	"github.com/stretchr/testify/assert"
)

func TestFssLogMeta(t *testing.T) {
	day, _ := strptime.Parse("2016-04-03", "%Y-%m-%d")
	seq := utils.NewTimeHourSequencer(day, day.Add(time.Second))
	meta := NewFssLogMeta("/1/2/3/bos_srv-k01a.fss.log.%Y%m%d%H", seq)
	if !assert.True(t, meta.Next()) {
		return
	}
	if !assert.Equal(t, meta.Name(), "/1/2/3/bos_srv-k01a.fss.log.2016040300") {
		return
	}

	dtime, _ := strptime.Parse("2016-04-03T01:02:03", "%Y-%m-%dT%H:%M:%S")
	tstamp, ok := meta.LineTimestamp([]byte("2016-04-03 01:02:03|ahahah"))
	if !assert.True(t, ok) {
		return
	}
	if !assert.Equal(t, dtime.UnixNano(), tstamp) {
		return
	}
	tstamp, ok = meta.LineTimestamp([]byte("2016-04-03 01:02:15|ahahah"))
	if !assert.True(t, ok) {
		return
	}
	if !assert.Equal(t, dtime.UnixNano()+12*int64(time.Second), tstamp) {
		return
	}
}
