package output

import "scanner"

type LineItem struct {
	data  []byte
	stamp int64
}

func (item LineItem) Line() string {
	return string(item.data)
}

type Splitter struct {
	prebuf int64

	scanner *scanner.Scanner
	channel chan []LineItem

	lastLine  []byte
	lastStamp int64

	buf []byte
}

func NewSplitter(prebuf int64, scn *scanner.Scanner, ch chan []LineItem) *Splitter {
	return &Splitter{
		prebuf:    prebuf,
		scanner:   scn,
		channel:   ch,
		lastLine:  nil,
		lastStamp: 0,
		buf:       make([]byte, 0, 64*1024),
	}
}

func (idx *Splitter) GetChannel() chan []LineItem {
	return idx.channel
}

func (idx *Splitter) reallocate(line []byte) []byte {
	prev_len := len(idx.buf)
	idx.buf = append(idx.buf, line...)
	return idx.buf[prev_len:]
}

func (idx *Splitter) refreshBuf() {
	capacity := 64 * 1024
	if cap(idx.buf) > capacity {
		capacity = cap(idx.buf)
	}
	idx.buf = make([]byte, 0, capacity)
}

func (idx *Splitter) Split() {
	if idx.lastLine == nil {
		if !idx.scanner.Scan() {
			return
		} else {
			idx.lastLine, idx.lastStamp = idx.scanner.Data()
			idx.lastLine = idx.reallocate(idx.lastLine)
		}
	}

	estimatedCapacity := 1000
	bound := idx.lastStamp
	for {
		bound = bound + idx.prebuf*((idx.lastStamp-bound)/idx.prebuf+1)
		buf := make([]LineItem, 1, estimatedCapacity)

		buf[0] = LineItem{
			data:  idx.lastLine,
			stamp: idx.lastStamp,
		}
		idx.lastLine = nil
		count := 0
		for idx.scanner.Scan() {
			count += 1
			line, stamp := idx.scanner.Data()

			if stamp >= bound {
				idx.refreshBuf()
			}
			line = idx.reallocate(line)

			if stamp >= bound {
				idx.lastLine = line
				idx.lastStamp = stamp
				if cap(buf) > estimatedCapacity {
					estimatedCapacity = cap(buf)
				}
				break
			}

			buf = append(buf, LineItem{
				data:  line,
				stamp: stamp,
			})
		}

		if count == 0 {
			idx.channel <- nil
			break
		} else {
			idx.channel <- buf
		}
	}
}
