package motion

import (
	"math"
)

func (cfg ProfileConfig) search(max float64, check func(float64) (ok bool, val interface{}), same func(a, b interface{}) bool) float64 {
	var n, min float64
	var val, newVal interface{}
	var ok bool

	n = 10
	for max == 0 {
		ok, newVal = check(n)
		if !ok {
			max = n
			break
		}

		if n == math.Inf(1) {
			return 0
		}
		val = newVal
		n *= n
	}

	for {
		n = (min + max) / 2
		ok, newVal = check(n)
		if ok {

			min = n
		} else {
			max = n
			if math.Abs(max-min) < (cfg.Epsilon / 10) {
				return min
			}
		}

		if same(val, newVal) {
			return min
		}

		val = newVal
	}
}
