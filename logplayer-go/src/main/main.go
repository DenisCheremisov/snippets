package main

import (
	"bos"
	"os"
	"output"
	"scanner"
	"time"
	"utils"

	"github.com/pbnjay/strptime"
)

func main() {
	from, _ := strptime.Parse("2016-01-14 01:02:00", "%Y-%m-%d %H:%M:%S")
	to, _ := strptime.Parse("2016-01-14 01:02:01", "%Y-%m-%d %H:%M:%S")
	meta := bos.NewErrLogMeta(""+
		"/home/emacs/Tmp/logs/icq/bos-k011/"+
		"bos_srv-k011e.err.%Y%m%d", utils.NewTimeHourSequencer(from, to))
	scn := scanner.NewScanner(meta)
	channel := make(chan []output.LineItem)

	splitter := output.NewSplitter(3*int64(time.Second), scn, channel)
	speedFactor := int64(100)
	sleeper := output.NewLogSleeper(speedFactor)
	outp := output.NewOutput(os.Stdout, sleeper)

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
