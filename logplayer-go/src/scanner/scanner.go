package scanner

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

type TimeSequencer interface {
	Next() bool
	Time() time.Time
}

// FileMeta define an abstraction over some log file's attributes
type LogMeta interface {
	Next() bool
	Name() string
	LineTimestamp([]byte) (int64, bool)
}

type Scanner struct {
	r    *bufio.Reader
	rraw *os.File

	line   []byte
	tstamp int64

	res       []byte
	resTstamp int64

	meta LogMeta
}

func NewScanner(meta LogMeta) *Scanner {
	return &Scanner{
		meta: meta,

		r:         nil,
		rraw:      nil,
		line:      nil,
		tstamp:    0,
		res:       nil,
		resTstamp: 0,
	}
}

func (s *Scanner) reopen() (ok bool, err error) {
	if s.rraw != nil {
		s.rraw.Close()
		s.rraw = nil
		s.r = nil
	}
	if !s.meta.Next() {
		return false, nil
	}
	s.rraw, err = os.Open(s.meta.Name())
	if err != nil {
		return false, err
	}
	s.r = bufio.NewReader(s.rraw)
	return true, nil
}

func (s *Scanner) Scan() bool {
	if s.r == nil {
		ok, err := s.reopen()
		if !ok {
			if err != nil {
				log.Fatal(err)
			}
			return false
		}
	}
	for {
		line, is_prefix, err := s.r.ReadLine()
		if err == io.EOF && len(line) == 0 {
			if len(s.line) > 0 {
				s.res = s.line
				s.resTstamp = s.tstamp
				s.line = line
				return true
			} else {
				ok, err := s.reopen()
				if err != nil {
					log.Fatal(err)
				}
				if !ok {
					return false
				}
			}
		}
		for is_prefix {
			var add []byte
			add, is_prefix, err = s.r.ReadLine()
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			new_line := make([]byte, len(line)+len(add))
			copy(new_line[:len(line)], line)
			copy(new_line[len(line):], add)
			line = new_line
			if err == io.EOF {
				break
			}
		}
		tstamp, ok := s.meta.LineTimestamp(line)
		if !ok {
			if s.line == nil {
				continue
			} else {
				new_line := make([]byte, len(s.line)+len(line)+1)
				copy(new_line[:len(s.line)], s.line)
				new_line[len(s.line)] = '\n'
				copy(new_line[len(s.line)+1:], line)
				s.line = new_line
				continue
			}
		}
		if s.line == nil {
			s.line = line
			s.tstamp = tstamp
			continue
		}
		s.res = s.line
		s.resTstamp = s.tstamp
		s.line = line
		s.tstamp = tstamp
		return true
	}
}

func (s *Scanner) Data() (line []byte, timestamp int64) {
	return s.res, s.resTstamp
}
