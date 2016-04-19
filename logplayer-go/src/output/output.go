package output

import (
	"io"
	"time"
)

type Output struct {
	firstLogStamp    int64
	firstOutputStamp int64
	writer           io.Writer
	sleeper          Sleeper

	buf []byte
}

type Sleeper interface {
	Delta(int64) int64
	Sync()
}

func NewOutput(writer io.Writer, sleeper Sleeper) *Output {
	return &Output{
		writer:           writer,
		firstLogStamp:    0,
		firstOutputStamp: 0,
		buf:              make([]byte, 0, 128*1024),
		sleeper:          sleeper,
	}
}

func (o *Output) append(line []byte) {
	o.buf = append(o.buf, line...)
	o.buf = append(o.buf, '\n')
}

func (o *Output) Write(res []LineItem) bool {
	if len(res) == 0 {
		return false
	}
	lastStamp := res[0].stamp
	o.sleeper.Sync()
	for _, line := range res {
		if line.stamp == lastStamp {
			o.append(line.data)
			continue
		}
		delta := o.sleeper.Delta(lastStamp)
		lastStamp = line.stamp
		time.Sleep(time.Duration(delta))
		o.writer.Write(o.buf)
		o.buf = o.buf[:0]
		o.append(line.data)
	}
	if len(o.buf) > 0 {
		time.Sleep(time.Duration(o.sleeper.Delta(lastStamp)))
		o.writer.Write(o.buf)
		o.buf = o.buf[:0]
	}
	return true
}
