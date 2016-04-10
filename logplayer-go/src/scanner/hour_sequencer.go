package scanner

import (
	"log"
	"time"
)

// TimeHourSequencer represents 1-hour time sequence between
type TimeHourSequencer struct {
	from time.Time
	to   time.Time
	cur  time.Time
}

func NewTimeHourSequencer(from time.Time, to time.Time) TimeHourSequencer {
	if from.After(to) {
		log.Printf("from must be after to, got %s > %s instead",
			from.Format(time.RFC3339), to.Format(time.RFC3339))
		panic("Wrong from and to precedence")
	}
	from = time.Date(
		from.Year(), from.Month(), from.Day(), from.Hour(), 0, 0, 0, from.Location())
	return TimeHourSequencer{
		from: from, to: to, cur: from.Add(-time.Second * 3600),
	}
}

func (seq *TimeHourSequencer) Next() bool {
	seq.cur = seq.cur.Add(time.Second * 3600)
	return seq.cur.Before(seq.to)
}
