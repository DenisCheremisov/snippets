package scanner

import (
	"log"
	"time"
)

// TimeHourSequencer represents 1-hour time sequence between
type BasicTimeSequencer struct {
	from time.Time
	to   time.Time
	cur  time.Time
}

func NewBasicTimeSequencer(from time.Time, to time.Time) BasicTimeSequencer {
	if from.After(to) {
		log.Printf("from must be after to, got %s > %s instead",
			from.Format(time.RFC3339), to.Format(time.RFC3339))
		panic("Wrong from and to precedence")
	}
	from = time.Date(
		from.Year(), from.Month(), from.Day(), from.Hour(), 0, 0, 0, from.Location())
	return BasicTimeSequencer{
		from: from, to: to, cur: from.Add(-time.Second * 3600),
	}
}

func (seq *BasicTimeSequencer) Next() bool {
	seq.cur = seq.cur.Add(time.Second * 3600)
	return seq.cur.Before(seq.to)
}

func (seq *BasicTimeSequencer) Time() time.Time {
	return seq.cur
}

type TimeHourSequencer struct {
	BasicTimeSequencer
}

func NewTimeHourSequencer(from time.Time, to time.Time) *TimeHourSequencer {
	return &TimeHourSequencer{NewBasicTimeSequencer(from, to)}
}

func (seq *TimeHourSequencer) Next() bool {
	seq.cur = seq.cur.Add(time.Second * 3600)
	return seq.cur.Before(seq.to)
}

type TimeDaySequencer struct {
	BasicTimeSequencer
}

func NewTimeDaySequencer(from time.Time, to time.Time) *TimeDaySequencer {
	return &TimeDaySequencer{NewBasicTimeSequencer(from, to)}
}

func (seq *TimeDaySequencer) Next() bool {
	seq.cur = seq.cur.AddDate(0, 0, 1)
	return seq.cur.Before(seq.to)
}
