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
}

func NewSplitter(prebuf int64, scn *scanner.Scanner, ch chan []LineItem) *Splitter {
	return &Splitter{
		prebuf:    prebuf,
		scanner:   scn,
		channel:   ch,
		lastLine:  nil,
		lastStamp: 0,
	}
}

func (idx *Splitter) GetChannel() chan []LineItem {
	return idx.channel
}

func (idx *Splitter) Split() {
	if idx.lastLine == nil {
		if !idx.scanner.Scan() {
			return
		} else {
			idx.lastLine, idx.lastStamp = idx.scanner.Data()
		}
	}

	estimatedCapacity := 1000
	bound := idx.lastStamp
	for {
		bound = bound + idx.prebuf*((idx.lastStamp-bound)/idx.prebuf+1)
		buf := make([]LineItem, estimatedCapacity)

		buf[0] = LineItem{
			data:  idx.lastLine,
			stamp: idx.lastStamp,
		}
		curLength := 1
		idx.lastLine = nil
		count := 0
		for idx.scanner.Scan() {
			count += 1
			line, stamp := idx.scanner.Data()

			if stamp >= bound {
				idx.lastLine = line
				idx.lastStamp = stamp
				break
			}
			if curLength == len(buf) {
				newBuf := make([]LineItem, curLength*2)
				copy(newBuf, buf)
				buf = newBuf
				estimatedCapacity = len(buf)
			}
			buf[curLength] = LineItem{
				data:  line,
				stamp: stamp,
			}
			curLength += 1
		}

		idx.channel <- buf[:curLength]
		if count == 0 {
			break
		}
	}
}
