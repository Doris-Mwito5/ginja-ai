package null

import (
	"time"
)

type nullType interface {
	bool | float32 | float64 | int | int64 | string | time.Time
}

func NullValue[T nullType](v T) *T {
	return &v
}

func ValueFromNull[T nullType](v *T) T {

	if v != nil {
		return *v
	}

	var zeroValue T

	return zeroValue
}
