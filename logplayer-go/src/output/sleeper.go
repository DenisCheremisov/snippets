package output

import "time"

type LogSleeper struct {
	firstImageStamp int64
	lastImageStamp  int64
	firstRealStamp  int64
	factor          int64

	sync bool
}

func NewLogSleeper(factor int64) *LogSleeper {
	return &LogSleeper{
		firstImageStamp: 0,
		lastImageStamp:  0,
		firstRealStamp:  0,
		factor:          factor,
		sync:            true,
	}
}

func (s *LogSleeper) Delta(stamp int64) int64 {
	if s.firstRealStamp == 0 {
		s.firstRealStamp = time.Now().UnixNano()
		s.firstImageStamp = stamp
		s.lastImageStamp = stamp
	}
	var delta int64
	if s.sync {
		now := time.Now().UnixNano()
		delta = (stamp-s.firstImageStamp)/s.factor - (now - s.firstRealStamp)
		s.sync = false
	} else {
		delta = (stamp - s.lastImageStamp) / s.factor
	}
	s.lastImageStamp = stamp
	return delta
}

func (s *LogSleeper) Sync() {
	s.sync = true
}
