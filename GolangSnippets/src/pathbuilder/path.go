package pathbuilder

import "fmt"

// Path represents protobuf's node path which is defined in
// https://github.com/google/protobuf/blob/b189389e2f2ca01dd534a8e9ba3ac38ea45cdba6/src/google/protobuf/descriptor.proto#L718
type Path []int32

func oveflow(v interface{}) {
	panic(fmt.Errorf("value %d (%T) is out of range for int32", v, v))
}

// Append creates a New path consisting of all current Path elements with `data` appended
// All data elements must be either of int or uint types. Every other type causes panicking.
func (p Path) Append(data ...interface{}) (result Path) {
	res := make([]int32, len(p)+len(data))
	copy(res, []int32(p))
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
		result = Path(res)
	}()
	for i, value := range data {
		var val int32
		switch v := value.(type) {
		case int8:
			val = int32(v)
		case int16:
			val = int32(v)
		case int32:
			val = v
		case int64:
			if v != int64(int32(v)) {
				oveflow(value)
			}
			val = int32(v)
		case int:
			if v != int(int32(v)) {
				oveflow(value)
			}
			val = int32(v)
		case uint8:
			val = int32(v)
		case uint16:
			val = int32(v)
		case uint32:
			if v > 0x7fffffff {
				oveflow(value)
			}
			val = int32(v)
		case uint64:
			if v > 0x7fffffff {
				oveflow(value)
			}
			val = int32(v)
		case uint:
			if v > 0x7fffffff {
				oveflow(value)
			}
			val = int32(v)
		default:
			panic(fmt.Errorf("value of type %T cannot be appended, only integer types are supported", value))
		}
		res[i+len(p)] = val
	}
	return
}

// Cut copies all but last `n` elements from current Path into the new Path and returns it
func (p Path) Cut(n int) Path {
	if n > len(p) {
		panic(
			fmt.Errorf(
				"attempt to cut out more elements from the path than it actually has (%d > %d)",
				n,
				len(p),
			),
		)
	}
	res := make([]int32, len(p)-n)
	copy(res, p[:len(p)-n])
	return Path(res)
}
