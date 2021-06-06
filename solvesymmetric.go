package motion

import (
	"errors"
	"math"
	"time"
)

type plan struct {
	T     time.Duration
	Val   float64
	Order int
}

func makeSymmetricBaseValues(src []float64, order int, value float64) []float64 {
	if order < 1 {
		panic("invalid order")
	}
	if order == 1 {
		return append(src[:0], value)
	}

	src = makeSymmetricBaseValues(src, order-1, value)
	src = append(src, 0)
	for _, r := range src[:len(src)-1] {
		src = append(src, -r)
	}

	return src
}

func makeSymmetricTimes(src []float64, order int, times []float64) []float64 {
	if order < 1 {
		panic("invalid order")
	}
	if order == 1 {
		return append(src[:0], times[0])
	}

	src = makeSymmetricTimes(src, order-1, times[1:])
	src = append(src, times[0])
	src = append(src, src[:len(src)-1]...)

	return src
}

func (cfg ProfileConfig) isValid(p Profile) bool {
	for _, s := range p {
		if cfg.gt(s.s0[0], cfg.Params[0].Target) {
			return false
		}
		for i, v := range s.s0 {
			if cfg.Params[i].Max == 0 {
				continue
			}

			if cfg.gt(v, cfg.Params[i].Max) {
				return false
			}
		}
	}

	return true
}
func (cfg ProfileConfig) equal(a, b float64) bool {
	return math.Abs(a-b) < cfg.Epsilon
}
func (cfg ProfileConfig) gt(a, b float64) bool {
	return a-cfg.Epsilon > b
}
func (cfg ProfileConfig) isExact(p Profile) bool {
	final := p[len(p)-1]

	for i, param := range cfg.Params {
		if cfg.equal(final.s0[i], param.Target) {
			continue
		}

		return false
	}

	return true
}

func (cfg ProfileConfig) stateCmp(a, b interface{}) bool {
	if a == nil || b == nil {
		return false
	}

	return cfg.sameState(a.(State), b.(State))
}

func (cfg ProfileConfig) solveSymmetric() (Profile, error) {
	base := make([]float64, 0, 1<<cfg.Order())
	times := make([]float64, cfg.Order())
	timeSlots := make([]float64, 0, 1<<cfg.Order())
	p := make(Profile, 1, 1<<cfg.Order())
	for _, param := range cfg.Params {
		p[0].s0 = append(p[0].s0, param.Start)
	}

	for order := 1; order <= cfg.Order(); order++ {
		base := makeSymmetricBaseValues(base, order, cfg.Params[cfg.Order()].Max)

		calc := func() (bool, interface{}) {
			timeSlots = makeSymmetricTimes(timeSlots, order, times)

			p = p[:1]
			for i, t := range timeSlots {
				if t == 0 {
					continue
				}
				p[len(p)-1].s0[cfg.Order()] = base[i]

				p = append(p, p.Last().Next(t))
			}

			p[len(p)-1].s0[cfg.Order()] = cfg.Params[cfg.Order()].Target

			return cfg.isValid(p), p.Last()
		}

		doSearch := func(index int) bool {
			if cfg.Params[index+1].Max == 0 {
				times[cfg.Order()-index-1] = 0
				calc()
				return cfg.isExact(p)
			}

			times[index] = cfg.search(times[index], func(v float64) (bool, interface{}) {
				times[index] = v
				return calc()
			}, cfg.stateCmp)

			calc()
			return cfg.isValid(p) && cfg.isExact(p)
		}
		var searchAll func(startIndex int) (bool, interface{})
		searchAll = func(startIndex int) (bool, interface{}) {
			if startIndex == cfg.Order() {
				return false, nil
			}

			// reset
			for i := range times {
				times[i] = 0
			}

			for i := range times {
				if !doSearch(cfg.Order() - i - 1) {
					continue
				}

				return true, p.Last()
			}

			if times[startIndex] > 0 {
				times[startIndex] = cfg.search(times[startIndex], func(v float64) (bool, interface{}) {
					times[startIndex] = v
					return searchAll(startIndex + 1)
				}, cfg.stateCmp)
			}

			calc()
			return cfg.isValid(p) && cfg.isExact(p), p.Last()
		}

		if ok, _ := searchAll(0); !ok {
			continue
		}

		return p, nil
	}

	return nil, errors.New("could not find a solution for the given constraints")
}
