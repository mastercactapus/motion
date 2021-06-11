package motion

import (
	"fmt"
)

type Profile []State

// Duration will return the total time to complete the motion Profile.
func (p Profile) Duration() float64 {
	return p[len(p)-1].t0
}

func (p Profile) State(t float64) State {
	for i := len(p) - 1; i >= 0; i-- {
		if t < p[i].t0 {
			continue
		}
		return p[i].Next(t - p[i].t0)
	}

	panic("bad profile")
}

func (p Profile) Last() State  { return p[len(p)-1].clone() }
func (p Profile) First() State { return p[0] }

type ProfileConfig struct {
	Epsilon float64
	TargetT float64
	Params  []Parameter
}

func (cfg ProfileConfig) Order() int { return len(cfg.Params) - 1 }

type Parameter struct {
	Start  float64
	Target float64
	Max    float64
}

// Solve will attempt to find a time-optimized solution to the given targets and constraints.
func (cfg ProfileConfig) Solve() (Profile, error) {
	if cfg.Epsilon == 0 {
		cfg.Epsilon = 0.0001
	}
	for i, p := range cfg.Params[1:] {
		if p.Start == p.Target {
			continue
		}

		return nil, fmt.Errorf("param[%d]: asymmetric profiles not yet supported (Start != Target)", i)
	}

	baseMax := cfg.Params[len(cfg.Params)-1].Max
	if baseMax == 0 && cfg.TargetT == 0 {
		return nil, fmt.Errorf("param[%d]: final parameter must include Max value or TargetT must be specified", len(cfg.Params)-1)
	}

	if cfg.TargetT == 0 {
		return cfg.solveSymmetric()
	}

	last := len(cfg.Params) - 1
	cfg.Params[last].Max = cfg.search(0, func(maxVal float64) (bool, interface{}) {
		cfg.Params[last].Max = maxVal
		prof, err := cfg.solveSymmetric()
		if err != nil {
			return false, nil
		}

		return !cfg.gt(cfg.TargetT, prof.Duration()), prof
	}, func(a, b interface{}) bool {
		if a == nil || b == nil {
			return false
		}

		return cfg.equal(a.(Profile).Duration(), b.(Profile).Duration())
	})

	return cfg.solveSymmetric()
}
