package motion

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfile_Solve(t *testing.T) {

	cfg := ProfileConfig{
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
			{Target: 10},
			{},
			{Max: 2},
		},
	}

	p, err = cfg.Solve()
	assert.NoError(t, err)
	assert.InDelta(t, 10, p.Last().Value(0), .01)
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

	fmt.Println(p.Duration())
	// output: 2h
}
