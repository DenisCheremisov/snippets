package main

import (
	"os"
	"output"
	"scanner"
	"time"

	"github.com/pbnjay/strptime"
)

func main() {
	from, _ := strptime.Parse("2016-01-13 01:02:00", "%Y-%m-%d %H:%M:%S")
	to, _ := strptime.Parse("2016-01-13 01:02:01", "%Y-%m-%d %H:%M:%S")
	meta := scanner.NewFssLogMeta(""+
		"/home/emacs/Sources/snippets/logplayer-go/src/output/"+
		"data_test/bos_srv-k011a.fss.log.{DATE}", from, to)
	scn := scanner.NewScanner(meta)
	channel := make(chan []output.LineItem)

	splitter := output.NewSplitter(3*int64(time.Second), scn, channel)
	outp := output.NewOutput(os.Stdout, 1)

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
