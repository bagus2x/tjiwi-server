package repository

import (
	"math"
	"testing"

	"github.com/bagus2x/tjiwi/pkg/history"
)

func TestDescendingFilterBuilder(t *testing.T) {
	p := history.Params{
		StorageID: 1,
		Status:    "deleted",
		StartDate: math.MaxInt16,
		EndDate:   math.MaxInt32,
		Cursor:    math.MaxInt8,
		Direction: "prev",
	}

	str, v := descendingFilter(&p)
	t.Log(str)
	t.Log(v)
}

func BenchmarkDescendingFilterBuilder(b *testing.B) {
	p := history.Params{
		StorageID: 1,
		Status:    "deleted",
		StartDate: math.MaxInt16,
		EndDate:   math.MaxInt32,
		Cursor:    math.MaxInt8,
		Direction: "asc",
	}

	for i := 0; i < b.N; i++ {
		str, v := descendingFilter(&p)
		_ = str
		_ = v
	}
}
