package utils

import (
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
