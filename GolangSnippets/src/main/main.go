package main

import (
	"arr"
	"bufio"
	"io"
	"log"
	"os"
)

const N = 100000000

func main() {
	tab := '\t'
	nl := '\n'

	take1 := dateseq.TakeBytesUntil(dateseq.SearchByte('\t'))
	take2 := dateseq.TakeBytesUntil(dateseq.SearchByte('\t'))
	take3 := dateseq.TakeBytesUntil(dateseq.SearchByte('\t'))
	take4 := dateseq.TakeBytesUntil(dateseq.SearchByte('\t'))
	take5 := dateseq.TakeBytesUntil(dateseq.EndSeeker(true))
	matcher := dateseq.NewMatcher(take1, take2, take3, take4, take5)

	stdout := bufio.NewWriterSize(os.Stdout, 512*1024)
	defer stdout.Flush()
	stdin := bufio.NewReaderSize(os.Stdin, 512*1024)
	for {
		line, err := stdin.ReadBytes('\n')
		if len(line) > 0 {
			if ok, rest := matcher.Feed(line); ok {
				stdout.Write(take1.Bytes())
				stdout.WriteRune(tab)
				stdout.Write(take2.Bytes())
				stdout.WriteRune(nl)
			} else {
				log.Println("ERROR: " + string(rest))
			}
		}

		if err != nil {
			if err != io.EOF {
				log.Fatalf("ERROR: %s", err)
			} else {
				break
			}
		}
	}
}
