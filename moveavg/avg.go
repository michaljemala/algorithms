package moveavg

import (
	"container/ring"
	"fmt"
	"math"
	"math/big"
)

var (
	ErrInvalidWindow   = fmt.Errorf("invalid window")
	ErrNotEnoughValues = fmt.Errorf("not enough values")
	ErrWindowTooBig    = fmt.Errorf("window too big")
)

type Averager struct {
	r *ring.Ring
}

const defaultRingSize = 10

func NewAverager(size int) *Averager {
	if size <= 0 {
		size = defaultRingSize
	}
	return &Averager{
		r: ring.New(defaultRingSize),
	}
}

func (a *Averager) Add(v float64) {
	if math.IsNaN(v) { // skip NaNs
		return
	}
	a.r.Value = big.NewFloat(v)
	a.r = a.r.Next()
}

func (a *Averager) MovingAverage(n int) (float64, error) {
	if n <= 0 {
		return 0, ErrInvalidWindow
	}
	if n > a.r.Len() {
		return 0, ErrWindowTooBig
	}

	p, sum := a.r.Prev(), new(big.Float)
	for i := 0; i < n; i++ {
		if p.Value == nil {
			return 0, ErrNotEnoughValues
		}
		sum = sum.Add(sum, p.Value.(*big.Float))
		p = p.Prev()
	}

	v, _ := sum.Float64()
	return v / float64(n), nil
}
