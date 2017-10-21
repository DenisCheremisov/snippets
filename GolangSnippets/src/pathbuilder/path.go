package pathbuilder

import "fmt"

// Path representing protobuf node path
type Path []int32

func oveflow(v interface{}) {
	panic(fmt.Errorf("value %d (%T) is out of range for int32", v, v))
}

// Append appends elements to the path
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
			panic(fmt.Errorf("unsupported type %T", value))
		}
		res[i+len(p)] = val
	}
	return
}

// Cut cuts last offset elements from the path's tail
func (p Path) Cut(offset int) Path {
	if offset > len(p) {
		panic(fmt.Errorf("try to cut more elements from the path than it has (%d > %d)", offset, len(p)))
	}
	res := make([]int32, len(p)-offset)
	copy(res, p[:len(p)-offset])
	return Path(res)
}
