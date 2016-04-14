package bos

import (
	"log"
	"time"
	"utils"

	"github.com/jehiah/go-strftime"
)

type ErrLogMeta struct {
	seq utils.TimeSequencer

	path          string
	prevDay       int
	prevDayTstamp int64
}

func NewErrLogMeta(path string, seq utils.TimeSequencer) *ErrLogMeta {
	return &ErrLogMeta{
		seq:           seq,
		path:          path,
		prevDay:       -1111,
		prevDayTstamp: 0,
	}
}

func (meta *ErrLogMeta) Next() bool {
	return meta.seq.Next()
}

func (meta *ErrLogMeta) Name() string {
	return strftime.Format(meta.path, meta.seq.Time())
}

func (meta *ErrLogMeta) LineTimestamp(line []byte) (int64, bool) {
	if len(line) < 10 || line[9] != ':' {
		return 0, false
	}

	day := int(line[0]-'0')*10 + int(line[1]-'0')
	if day != meta.prevDay {
		curDay := meta.seq.Time()
		if curDay.AddDate(0, 0, -1).Day() == day {
			curDay = curDay.AddDate(0, 0, -1)
		} else if curDay.AddDate(0, 0, 1).Day() == day {
			curDay = curDay.AddDate(0, 0, 1)
		} else if curDay.Day() != day {
			log.Fatalf("Cannot extract date from %s", string(line))
		}

		curDay = time.Date(
			curDay.Year(), curDay.Month(), curDay.Day(), curDay.Hour(),
			0, 0, 0, curDay.Location())
		meta.prevDayTstamp = curDay.UnixNano()
		meta.prevDay = day
	}

	return meta.prevDayTstamp +
		((int64(line[3]-'0')*10+int64(line[4]-'0'))*3600+
			(int64(line[5]-'0')*10+int64(line[6]-'0'))*60+
			(int64(line[7]-'0')*10+int64(line[8]-'0')))*int64(time.Second), true

}
