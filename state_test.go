package motion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_Next(t *testing.T) {
	s := State{s0: []float64{0, 0, 2}}

	n := s.Next(10)
	assert.Equal(t, State{s0: []float64{100, 20, 2}, t0: 10}, n)

	s = n
	n = s.Next(10)
	assert.Equal(t, State{s0: []float64{400, 40, 2}, t0: 20}, n)
}
