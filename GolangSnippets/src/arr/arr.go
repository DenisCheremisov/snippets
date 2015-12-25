package arr

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// DateSeq is a date sequence generator
type DateSeq struct {
	cur   time.Time
	until time.Time

	to_be_returned time.Time
}

// NewDateSeq create a new sequence generator for dates from start
// to finish. Start should Be Before Finish
func NewDateSeq(start time.Time, finish time.Time) (*DateSeq, error) {
	if start.After(finish) {
		return nil, errors.New("ERROR: start must be before the finish")
	}
	return &DateSeq{start, finish, start}, nil
}

// Next Check If The Sequence Has Not Been Passed Yet
func (seq *DateSeq) Next() bool {
	if seq.cur.After(seq.until) {
		return false
	} else {
		seq.to_be_returned = seq.cur
		seq.cur = seq.cur.AddDate(0, 0, 1)
		return true
	}
}

// Date returns current date of the sequence
func (seq *DateSeq) Date() time.Time {
	return seq.to_be_returned
}

// NewStoreReader tries to open a store named path
// If there is no such file it tries to open path.gz file and
// treat it as gzipped file.
func NewStoreReader(path string) (io.Reader, error) {
	if file, err := os.Open(path); err == nil {
		return file, nil
	}
	if file, err := os.Open(path + ".gz"); err == nil {
		if gzipped, err := gzip.NewReader(file); err == nil {
			return gzipped, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("ERROR: cannot open %s[.gz]", path))
}
