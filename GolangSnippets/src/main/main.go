package main

import (
	"arr"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

const N = 100000000

func main() {
	pass15 := arr.PassFuncCheck(func(x byte) bool {
		return x == byte('1') || x == byte('5')
	})
	matcher := arr.NewMatcher(
		arr.PassBytesUntil(arr.SearchByte('#')),
		pass15, pass15,
		arr.PassPattern([]byte("1110")),
		arr.PassBytesUntil(arr.SearchUntilByte(byte('0'))),
		arr.OnTheEnd(true))

	stdout := bufio.NewWriterSize(os.Stdout, 512*1024)
	defer stdout.Flush()
	source := bufio.NewReaderSize(os.Stdin, 512*1024)
	for {
		line, err := source.ReadBytes('\n')
		if len(line) > 0 {
			if ok, _ := matcher.Feed(line[:len(line)-1]); ok {
				fmt.Printf("Match: %s", string(line))
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
