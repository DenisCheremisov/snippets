package main

import (
	"os"
	"output"
	"scanner"
	"time"

	"github.com/pbnjay/strptime"
)

func main() {
	from, _ := strptime.Parse("2016-01-14 01:02:00", "%Y-%m-%d %H:%M:%S")
	to, _ := strptime.Parse("2016-01-14 01:02:01", "%Y-%m-%d %H:%M:%S")
	meta := scanner.NewFssLogMeta(""+
		"/home/emacs/Desktop/"+
		"bos_srv-k011e.fss.log.%Y%m%d", scanner.NewTimeHourSequencer(from, to))
	scn := scanner.NewScanner(meta)
	channel := make(chan []output.LineItem)

	splitter := output.NewSplitter(3*int64(time.Second), scn, channel)
	speedFactor := int64(4)
	outp := output.NewOutput(os.Stdout, speedFactor)

	go func() {
		for {
			res := <-channel
			if !outp.Write(res) {
				break
			}
		}
	}()

	splitter.Split()
}
