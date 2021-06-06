package motion

import (
	"fmt"
	"sort"
)

type Profile []State

// Duration will return the total time to complete the motion Profile.
func (p Profile) Duration() float64 {
	return p[len(p)-1].t0
}

func (p Profile) State(t float64) State {
	idx := sort.Search(len(p), func(i int) bool {
		return p[i].Time() <= t
	})

	return p[idx].Next(t)
}

func (p Profile) Last() State  { return p[len(p)-1] }
func (p Profile) First() State { return p[0] }

type ProfileConfig struct {
	Epsilon float64
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
	for i, p := range cfg.Params[1:] {
		if p.Start == p.Target {
			continue
		}

		return nil, fmt.Errorf("param[%d]: asymmetric profiles not yet supported (Start != Target)", i)
	}

	if cfg.Params[len(cfg.Params)-1].Max == 0 {
		return nil, fmt.Errorf("param[%d]: final parameter must include Max value", len(cfg.Params)-1)
	}

	return cfg.solveSymmetric()
}
