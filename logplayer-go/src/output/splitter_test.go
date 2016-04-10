package output

import (
	"scanner"
	"testing"
	"time"

	"github.com/pbnjay/strptime"
	"github.com/stretchr/testify/assert"
)

func TestLogSplit(t *testing.T) {
	from, _ := strptime.Parse("2016-01-13 00:02:00", "%Y-%m-%d %H:%M:%S")
	to, _ := strptime.Parse("2016-01-13 00:02:01", "%Y-%m-%d %H:%M:%S")

	meta := scanner.NewFssLogMeta("data_test/bos_srv-k011a.fss.log.{DATE}", from, to)
	scanner := scanner.NewScanner(meta)
	channel := make(chan []LineItem)

	splitter := NewSplitter(5*int64(time.Second), scanner, channel)
	go splitter.Split()

	res := <-channel
	if !assert.Equal(t, len(res), 1) {
		return
	}
	stamp := res[0].stamp

	res = <-channel
	if !assert.Equal(t, len(res), 2) {
		return
	}
	if !assert.Equal(t, stamp+int64(time.Second*123), res[0].stamp) {
		return
	}
	if !assert.Equal(t, res[0].stamp+int64(time.Second), res[1].stamp) {
		return
	}

	res = <-channel
	if !assert.Equal(t, len(res), 1) {
		return
	}
	res = <-channel
	if !assert.Equal(t, len(res), 3) {
		return
	}
	if !assert.Equal(t, res[0].stamp, res[1].stamp) {
		return
	}
	if !assert.Equal(t, res[0].stamp, res[2].stamp) {
		return
	}
}
