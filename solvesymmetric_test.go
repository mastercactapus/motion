package motion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeSymmetricTimes(t *testing.T) {
	check := func(order int, exp []float64) {
		t.Helper()
		times := make([]float64, order)
		for i := range times {
			times[i] = float64(i + 1)
		}

		res := makeSymmetricTimes(nil, order, times)
		assert.Equal(t, exp, res)
	}

	check(1, []float64{1})
	check(2, []float64{2, 1, 2})
	check(3, []float64{3, 2, 3, 1, 3, 2, 3})
	check(4, []float64{4, 3, 4, 2, 4, 3, 4, 1, 4, 3, 4, 2, 4, 3, 4})

}
func TestMakeSymmetricBaseValues(t *testing.T) {
	check := func(order int, exp []float64) {
		t.Helper()
		res := makeSymmetricBaseValues(nil, order, 1)
		assert.Equal(t, exp, res)
	}

	check(1, []float64{1})
	check(2, []float64{1, 0, -1})
	check(3, []float64{1, 0, -1, 0, -1, 0, 1})
	check(4, []float64{1, 0, -1, 0, -1, 0, 1, 0, -1, 0, 1, 0, 1, 0, -1})
}
