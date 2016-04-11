package scanner

import (
	"bufio"
	"log"
	"os"
	"testing"
	"time"

	"github.com/pbnjay/strptime"
	"github.com/stretchr/testify/assert"
)

func TestTimeHourSequencer(t *testing.T) {
	from, _ := strptime.Parse("2016-04-03 10:12:13", "%Y-%m-%d %H:%M:%S")
	to, _ := strptime.Parse("2016-04-04 00:00:00", "%Y-%m-%d %H:%M:%S")

	seq := NewTimeHourSequencer(from, to)
	res := make([]time.Time, 0)
	for seq.Next() {
		res = append(res, seq.cur)
	}
	if !assert.Equal(t, len(res), 14) {
		return
	}
	if !assert.Equal(t, res[0].Format(time.RFC3339), "2016-04-03T10:00:00Z") {
		return
	}
	if !assert.Equal(t, res[13].Format(time.RFC3339), "2016-04-03T23:00:00Z") {
		return
	}
	for i := 1; i < 14; i++ {
		if !assert.Equal(t, res[i].Sub(res[i-1]), time.Second*3600) {
			return
		}
	}
}

func TestFssLogMeta(t *testing.T) {
	day, _ := strptime.Parse("2016-04-03", "%Y-%m-%d")
	seq := NewTimeHourSequencer(day, day.Add(time.Second))
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

func TestScanner(t *testing.T) {
	from, _ := strptime.Parse("2016-01-14 00:02:00", "%Y-%m-%d %H:%M:%S")
	to, _ := strptime.Parse("2016-01-14 00:02:01", "%Y-%m-%d %H:%M:%S")
	seq := NewTimeHourSequencer(from, to)

	meta := NewFssLogMeta("./data_test/bos_srv-k011e.fss.log.%Y%m%d%H", seq)
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
	seq = NewTimeHourSequencer(from, to)
	meta = NewFssLogMeta("./data_test/bos_srv-k011e.fss.log.%Y%m%d%H", seq)
	scanner = NewScanner(meta)
	for scanner.Scan() {
		line, tstamp := scanner.Data()
		log.Println(string(line), tstamp)
	}
}
