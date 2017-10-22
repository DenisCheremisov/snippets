package pathbuilder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {
	path := Path{}
	path = path.Append(1, 2, 3, 4)
	path = path.Append(12)
	require.Equal(t, Path([]int32{1, 2, 3, 4, 12}), path)

	path = path.Cut(3)
	require.Equal(t, Path([]int32{1, 2}), path)

	failures := []interface{}{
		"123",
		int64(0x7fffffff) + int64(1),
		-int64(0x80000000) - int64(3),
		uint32(0x80000000),
		uint64(0x80000000),
		int(0x7fffffff) + 1,
		uint(0x80000000),
	}
	for _, failure := range failures {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("Conversion of %v (type %T) must cause a panic", failure, failure)
				}
			}()
			path = path.Append(failure)
		}()
	}
}
