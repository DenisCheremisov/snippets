package output

import (
	"bos"
	"bufio"
	"log"
	"os"
	"scanner"
	"testing"
	"time"
	"utils"

	"github.com/pbnjay/strptime"
	"github.com/stretchr/testify/assert"
)

func TestLogSplit(t *testing.T) {
	from, _ := strptime.Parse("2016-01-13 00:02:00", "%Y-%m-%d %H:%M:%S")
	to, _ := strptime.Parse("2016-01-13 00:02:01", "%Y-%m-%d %H:%M:%S")
	seq := utils.NewTimeHourSequencer(from, to)

	meta := bos.NewFssLogMeta("data_test/bos_srv-k011a.fss.log.%Y%m%d%H", seq)
	scanner := scanner.NewScanner(meta)
	channel := make(chan []LineItem)

	splitter := NewSplitter(5*int64(time.Second), scanner, channel)
	go splitter.Split()

	file, err := os.Open("data_test/bos_srv-k011a.fss.log.2016011300")
	if err != nil {
		t.Error(err)
		return
	}
	scanner_sample := bufio.NewScanner(bufio.NewReader(file))
	samples := []string{}
	for scanner_sample.Scan() {
		text := scanner_sample.Text()
		if text != "" {
			samples = append(samples, scanner_sample.Text())
		}
	}
	if !assert.True(t, len(samples) > 0) {
		return
	}

	res := <-channel
	if !assert.Equal(t, len(res), 1) {
		return
	}
	if !assert.Equal(t, string(res[0].data), samples[0]) {
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
	if !assert.Equal(t, string(res[0].data), samples[1]) {
		return
	}
	if !assert.Equal(t, res[0].stamp+int64(time.Second), res[1].stamp) {
		return
	}
	if !assert.Equal(t, string(res[1].data), samples[2]) {
		return
	}

	res = <-channel
	if !assert.Equal(t, len(res), 1) {
		return
	}
	if !assert.Equal(t, string(res[0].data), samples[3]) {
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
	if !assert.Equal(t, string(res[0].data), samples[4]) {
		return
	}
	if !assert.Equal(t, string(res[1].data), samples[5]) {
		return
	}
	if !assert.Equal(t, string(res[2].data), samples[6]) {
		return
	}

	res = <-channel
	if res != nil {
		log.Println("res must be nil after file was read out")
	}

}
