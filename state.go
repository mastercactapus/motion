package motion

import (
	"math"
)

type State struct {
	s0 []float64
	t0 float64
}

func (s State) clone() State {
	n := make([]float64, len(s.s0))
	copy(n, s.s0)
	return State{t0: s.t0, s0: n}
}

func factorial(n int) int {
	sum := n
	for i := n - 1; i > 1; i-- {
		sum *= i
	}
	return sum
}
func (s State) Order() int { return len(s.s0) }
func (s State) Value(d Derivitive) float64 {
	if int(d) >= len(s.s0) {
		return 0
	}
	return s.s0[d]
}
func (s State) Time() float64 { return s.t0 }
func (s State) Next(t float64) State {
	s = s.clone()
	vals := make([]float64, len(s.s0))
	copy(vals, s.s0)

	for i := range vals {
		for j, v := range s.s0[i+1:] {
			vals[i] += v * math.Pow(t, float64(j+1)) / float64(factorial(j+1))
		}
	}

	return State{t0: s.t0 + t, s0: vals}
}
func (cfg ProfileConfig) sameState(a, b State) bool {
	if len(a.s0) != len(b.s0) {
		return false
	}

	for i, v := range a.s0 {
		if math.Abs(b.s0[i]-v) <= cfg.Epsilon {
			continue
		}

		return false
	}

	return true
}
