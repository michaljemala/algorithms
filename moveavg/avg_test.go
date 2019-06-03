package moveavg_test

import (
	"math"
	"testing"

	"github.com/michaljemala/algorithms/moveavg"
)

func TestMoveAvg(t *testing.T) {
	for _, tc := range []struct {
		name       string
		buffSize   int
		data       []float64
		windowSize int
		result     float64
		errFunc    func(*testing.T, error)
	}{
		{
			name:       "not enough captured items",
			buffSize:   50,
			data:       []float64{1.0, 1.3, 1.2, 1.4, 1.9},
			windowSize: 10,
			errFunc: func(t *testing.T, err error) {
				if err != moveavg.ErrNotEnoughValues {
					t.Fatalf("want %v, have %v", moveavg.ErrNotEnoughValues, err)
				}
			},
		},
		{
			name:       "too big window size",
			buffSize:   10,
			windowSize: 50,
			errFunc: func(t *testing.T, err error) {
				if err != moveavg.ErrWindowTooBig {
					t.Fatalf("want %v, have %v", moveavg.ErrWindowTooBig, err)
				}
			},
		},
		{
			name:     "incorrect input items",
			buffSize: 10,
			data: []float64{
				math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN(),
			},
			windowSize: 5,
			errFunc: func(t *testing.T, err error) {
				if err != moveavg.ErrNotEnoughValues {
					t.Fatalf("want %v, have %v", moveavg.ErrNotEnoughValues, err)
				}
			},
		},
		{
			name:     "invalid window",
			buffSize: 50,
			data: []float64{
				1.0, 1.3, 1.2, 1.4, 1.9,
				1.1, 1.4, 1.7, 1.1, 1.8,
			},
			windowSize: 0,
			errFunc: func(t *testing.T, err error) {
				if err != moveavg.ErrInvalidWindow {
					t.Fatalf("want %v, have %v", moveavg.ErrInvalidWindow, err)
				}
			},
		},
		{
			name:     "correct moving average",
			buffSize: 50,
			data: []float64{
				1.0, 1.3, 1.2, 1.4, 1.9,
				1.1, 1.4, 1.7, 1.1, 1.8,
				1.1, 1.1, 1.5, 1.4, 1.0,
			},
			windowSize: 10,
			result:     1.32,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := moveavg.NewAverager(tc.buffSize)
			for _, v := range tc.data {
				r.Add(v)
			}
			v, err := r.MovingAverage(tc.windowSize)
			if err != nil {
				if tc.errFunc == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				tc.errFunc(t, err)
			}
			if want, have := tc.result, v; !float64Equals(want, have) {
				t.Fatalf("want %v, have %v", want, have)
			}
		})
	}
}

func float64Equals(a, b float64) bool {
	const threshold = 1e-9
	return math.Abs(a-b) <= threshold
}
