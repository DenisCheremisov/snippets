package output

import (
	"io"
	"time"
)

type Output struct {
	firstLogStamp    int64
	firstOutputStamp int64
	writer           io.Writer
	speedup          int64

	buf       []byte
	curLength int
}

func NewOutput(writer io.Writer, speedup int64) *Output {
	return &Output{
		writer:           writer,
		firstLogStamp:    0,
		firstOutputStamp: 0,
		buf:              make([]byte, 128*1024),
		curLength:        0,
		speedup:          speedup,
	}
}

func (o *Output) append(line []byte) {
	for len(line)+o.curLength+1 > len(o.buf) {
		newBuf := make([]byte, len(o.buf)*2)
		copy(newBuf, o.buf)
		o.buf = newBuf
	}

	o.buf[o.curLength] = '\n'
	copy(o.buf[o.curLength+1:], line)
	o.curLength = o.curLength + 1 + len(line)
}

func (o *Output) Write(res []LineItem) bool {
	if len(res) == 0 {
		return false
	}
	delta := int64(0)
	lastStamp := res[0].stamp
	if o.firstOutputStamp == 0 {
		o.firstOutputStamp = time.Now().UnixNano()
		o.firstLogStamp = res[0].stamp
	} else {
		now := time.Now().UnixNano()
		delta = (lastStamp-o.firstLogStamp)/o.speedup - (now - o.firstOutputStamp)
	}

	for _, line := range res {
		if line.stamp == lastStamp {
			o.append(line.data)
			continue
		}
		if delta > 0 {
			time.Sleep(time.Duration(delta))
			delta = 0
		}
		delta += (line.stamp - lastStamp) / o.speedup
		lastStamp = line.stamp
		o.writer.Write(o.buf[:o.curLength])
		o.curLength = 0
		o.append(line.data)
	}
	if o.curLength > 0 {
		time.Sleep(time.Duration(delta))
		o.writer.Write(o.buf[:o.curLength])
		o.curLength = 0
	}
	return true
}
