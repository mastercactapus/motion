package motion

import (
	"fmt"
	"log"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfile_Solve(t *testing.T) {
	var cfg ProfileConfig

	cfg = ProfileConfig{
		Epsilon: .01,
		Params: []Parameter{
			{Target: 10},
			{Max: 5},
		},
	}

	p, err := cfg.Solve()
	assert.NoError(t, err)
	assert.InDelta(t, 2.0, p.Duration(), .01)

	cfg = ProfileConfig{
		Epsilon: .01,
		Params: []Parameter{
			{Target: 10}, // 10 steps
			{},           // no velocity config
			{Max: 2},     // max accel 2
		},
	}

	p, err = cfg.Solve()
	assert.NoError(t, err)
	assert.InDelta(t, 10, p.Last().Value(0), .01)

	cfg = ProfileConfig{
		Epsilon: 0.0001,
		Params: []Parameter{
			{Target: 100},
			{},
			{Max: 5},
		},
	}
	p, err = cfg.Solve()
	require.NoError(t, err)
	assert.InDelta(t, 100, p.Last().Value(0), 0.0001)

	// Time based

	cfg = ProfileConfig{
		Epsilon: 0.0001,
		Params: []Parameter{
			{Target: 100},
			{},
		},
		TargetT: 1000,
	}
	p, err = cfg.Solve()
	require.NoError(t, err)
	assert.InDelta(t, 100, p.Last().Value(0), 0.001)
	assert.InDelta(t, 1000, p.Duration(), 0.001)

	cfg = ProfileConfig{
		Epsilon: 0.001,
		Params: []Parameter{
			{Target: 100},
			{}, {},
		},
		TargetT: 1000,
	}
	p, err = cfg.Solve()
	require.NoError(t, err)
	assert.InDelta(t, 100, p.Last().Value(0), 0.001)
	assert.InDelta(t, 1000, p.Duration(), 0.001)

	cfg = ProfileConfig{
		Epsilon: 0.0001,
		Params: []Parameter{
			{Target: 100},
			{}, {}, {},
		},
		TargetT: 1000,
	}
	p, err = cfg.Solve()
	require.NoError(t, err)
	assert.InDelta(t, 100, p.Last().Value(0), 0.001)
	assert.InDelta(t, 1000, p.Duration(), 0.001)

	cfg = ProfileConfig{
		Epsilon: 0.0001,
		Params: []Parameter{
			{Target: 100},
			{}, {}, {}, {},
		},
		TargetT: 1000,
	}
	p, err = cfg.Solve()
	require.NoError(t, err)
	assert.InDelta(t, 100, p.Last().Value(0), 0.001)
	assert.InDelta(t, 1000, p.Duration(), 0.001)

	cfg = ProfileConfig{
		Epsilon: 0.0001,
		Params: []Parameter{
			{Target: 100},
			{}, {}, {}, {}, {},
		},
		TargetT: 1000,
	}
	p, err = cfg.Solve()
	require.NoError(t, err)
	assert.InDelta(t, 100, p.Last().Value(0), 0.001)
	assert.InDelta(t, 1000, p.Duration(), 0.001)

	cfg = ProfileConfig{
		Epsilon: 0.0001,
		Params: []Parameter{
			{Target: 100},
			{}, {}, {}, {}, {}, {},
		},
		TargetT: 1000,
	}
	p, err = cfg.Solve()
	require.NoError(t, err)
	assert.InDelta(t, 100, p.Last().Value(0), 0.001)
	assert.InDelta(t, 1000, p.Duration(), 0.001)

}

func ExampleProfile() {
	// profile config for a target of 10 miles, at a max of
	// 5 mph
	cfg := ProfileConfig{
		Params: []Parameter{
			{Target: 10},
			{Max: 5},
		},
	}

	p, err := cfg.Solve()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(math.Round(p.Duration()))
	// output: 2
}
