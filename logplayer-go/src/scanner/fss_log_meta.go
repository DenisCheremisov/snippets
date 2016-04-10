package scanner

import (
	"strings"
	"time"

	"github.com/jehiah/go-strftime"
	"github.com/pbnjay/strptime"
)

type FssLogMeta struct {
	TimeHourSequencer

	path          string
	prevDay       int
	prevDayTstamp int64
}

func NewFssLogMeta(path string, from time.Time, to time.Time) *FssLogMeta {
	return &FssLogMeta{
		TimeHourSequencer: NewTimeHourSequencer(from, to),
		path:              path,
		prevDay:           -1111,
		prevDayTstamp:     0,
	}
}

func (meta *FssLogMeta) Name() string {
	return strings.Replace(
		meta.path, "{DATE}", strftime.Format("%Y%m%d%H", meta.cur), 1)
}

func (meta *FssLogMeta) LineTimestamp(line []byte) (int64, bool) {
	if len(line) <= 19 || line[19] != '|' {
		return 0, false
	}
	day := int(line[8]-'0')*10 + int(line[9]-'0')
	if day != meta.prevDay {
		tstamp, err := strptime.Parse(string(line[:10]), "%Y-%m-%d")
		if err != nil {
			return 0, false
		}
		meta.prevDayTstamp = tstamp.UnixNano()
		meta.prevDay = day
	}
	return meta.prevDayTstamp +
		((int64(line[11]-'0')*10+int64(line[12]-'0'))*3600+
			(int64(line[14]-'0')*10+int64(line[15]-'0'))*60+
			(int64(line[17]-'0')*10+int64(line[18]-'0')))*int64(time.Second), true
}
