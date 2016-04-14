package output

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	var b bytes.Buffer

	data := []LineItem{
		LineItem{[]byte("0"), 0},
		LineItem{[]byte("1"), 1},
		LineItem{[]byte("2"), 2},
		LineItem{[]byte("3"), 3},
		LineItem{[]byte("4"), 4},
		LineItem{[]byte("5"), 5},
		LineItem{[]byte("6"), 6},
		LineItem{[]byte("7"), 7},
		LineItem{[]byte("8"), 8},
		LineItem{[]byte("9"), 9},
		LineItem{[]byte("ab"), 11},
	}

	start := time.Now().UnixNano()
	sleeper := NewLogSleeper(100)
	o := NewOutput(&b, sleeper)
	o.Write(data)
	end := time.Now().UnixNano()
	if !assert.True(t, (end-start) < int64(time.Second)/5) {
		return
	}

	output := string(b.Bytes())
	if !assert.Equal(t, output, "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\nab\n") {
		return
	}
}
