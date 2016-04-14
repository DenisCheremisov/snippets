package scanner

import (
	"bos"
	"bufio"
	"log"
	"os"
	"testing"
	"time"
	"utils"

	"github.com/pbnjay/strptime"
	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	from, _ := strptime.Parse("2016-01-14 00:02:00", "%Y-%m-%d %H:%M:%S")
	to, _ := strptime.Parse("2016-01-14 00:02:01", "%Y-%m-%d %H:%M:%S")
	seq := utils.NewTimeHourSequencer(from, to)

	meta := bos.NewFssLogMeta("./data_test/bos_srv-k011e.fss.log.%Y%m%d%H", seq)
	timestamps := []int64{
		1452729595, 1452729595, 1452729595, 1452729596, 1452729596}

	rd, _ := os.Open("./data_test/bos_srv-k011e.fss.log.2016011400")
	defer rd.Close()
	sc := bufio.NewScanner(bufio.NewReader(rd))

	scanner := NewScanner(meta)
	i := 0
	for scanner.Scan() {
		sc.Scan()
		line, tstamp := scanner.Data()
		if !assert.Equal(t, string(line), sc.Text()) {
			return
		}
		if !assert.Equal(t, tstamp, timestamps[i]*int64(time.Second)) {
			return
		}
		i = i + 1
	}

	from = from.Add(time.Second * 3600)
	to = to.Add(time.Second * 3600)
	seq = utils.NewTimeHourSequencer(from, to)
	meta = bos.NewFssLogMeta("./data_test/bos_srv-k011e.fss.log.%Y%m%d%H", seq)
	scanner = NewScanner(meta)
	for scanner.Scan() {
		line, tstamp := scanner.Data()
		log.Println(string(line), tstamp)
	}
}
