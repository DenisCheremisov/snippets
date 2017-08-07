package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
)

// TSVScanner splits input line into chunks separated with TAB character
type TSVScanner struct {
	orig []byte
	rest []byte

	cols   int
	curCol int
}

// NewTSVScanner constructor
func NewTSVScanner(colLimit int) *TSVScanner {
	return &TSVScanner{
		cols: colLimit,
	}
}

// SetSource ...
func (s *TSVScanner) SetSource(line []byte) {
	s.orig = line
	s.rest = line
	s.curCol = 1
}

// GetChunk ...
func (s *TSVScanner) GetChunk(index int) (chunk []byte, err error) {
	if index > s.cols {
		err = fmt.Errorf("Column index (%d) is out of bounds in `%s`", index, string(s.orig))
		return
	}

	for i := s.curCol; i < index; i++ {
		pos := bytes.IndexByte(s.rest, '\t')
		if pos < 0 {
			err = fmt.Errorf("Malformed line `%s`: not enough columns", string(s.orig))
			return
		}
		s.rest = s.rest[pos+1:]
	}

	s.curCol = index
	if s.curCol == s.cols {
		chunk = s.rest
		s.rest = s.rest[len(s.rest):]
		return
	}

	pos := bytes.IndexByte(s.rest, '\t')
	if pos < 0 {
		err = fmt.Errorf("Malformed line `%s`: not enough columns", string(s.orig))
		return
	}
	chunk = s.rest[:pos]
	s.rest = s.rest[pos+1:]
	s.curCol++
	return
}

// Atoi "fast" parsing
func Atoi(s []byte) (int, bool) {
	i := 0
	var x int
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		x = x*10 + int(c) - '0'
	}
	return x, true
}

// Process processed that google's TSV file
func Process(r io.Reader, keyIndex, valIndex, colLimit int) (key string, val int, err error) {
	counters := make(map[string]*int, 16384)
	countValues := make([]int, 0, 16384)
	indices := [2]int{keyIndex, valIndex}
	if keyIndex > valIndex {
		indices[0], indices[1] = indices[1], indices[0]
	}
	s := NewTSVScanner(colLimit)

	reader := bufio.NewReaderSize(r, 512*1024)
	scanner := bufio.NewScanner(reader)
	var t1, t2 []byte
	for scanner.Scan() {
		s.SetSource(scanner.Bytes())

		t1, err = s.GetChunk(indices[0])
		if err != nil {
			return
		}
		t2, err = s.GetChunk(indices[1])
		if err != nil {
			return
		}

		if keyIndex > valIndex {
			t1, t2 = t2, t1
		}
		v, ok := Atoi(t2)
		if !ok {
			err = fmt.Errorf("Cannot parse value `%s` into integer in `%s`", string(t2), scanner.Text())
			return
		}

		if val, ok := counters[string(t1)]; ok {
			*val += v
		} else {
			countValues = append(countValues, v)
			counters[string(t1)] = &(countValues[len(countValues)-1])
		}
	}

	key = ""
	val = -1
	for k, v := range counters {
		if *v > val {
			key = k
			val = *v
		}
	}
	return
}

func main() {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

	var keyIndex, valIndex, colLimit int
	flag.IntVar(&keyIndex, "key", 2, "key column index (1, 2, 3, ...)")
	flag.IntVar(&valIndex, "val", 3, "value column index (1, 2, 3, ...)")
	flag.IntVar(&colLimit, "cols", 4, "how many columns to scan")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	inputFileName := flag.Arg(0)
	file, err := os.Open(inputFileName)
	if err != nil {
		panic(err)
	}

	key, val, err := Process(file, keyIndex, valIndex, colLimit)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Max key: %s, max val: %d\n", key, val)
}
