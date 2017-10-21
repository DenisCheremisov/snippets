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
}
