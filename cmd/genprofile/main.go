package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/mastercactapus/motion"
)

func main() {
	steps := flag.Int("steps", 0, "Number of steps to perform.")
	startVel := flag.Float64("start-vel", 0, "Starting velocity.")
	endVel := flag.Float64("end-vel", 0, "Ending velocity.")
	maxVel := flag.Float64("vel", 0, "Max velocity.")
	maxAcc := flag.Float64("acc", 0, "Max acceleration.")
	maxJer := flag.Float64("jer", 0, "Max jerk.")
	maxSna := flag.Float64("sna", 0, "Max snap.")
	maxCra := flag.Float64("cra", 0, "Max crackle.")
	maxPop := flag.Float64("pop", 0, "Max pop.")
	ep := flag.Float64("e", .0001, "Allowable error.")
	t := flag.Float64("t", .001, "Time interval.")
	flag.Parse()

	cfg := motion.ProfileConfig{
		Epsilon: *ep,
		Params: []motion.Parameter{
			{Target: float64(*steps)},
			{Start: *startVel, Target: *endVel, Max: *maxVel},
		},
	}
	orders := []float64{*maxAcc, *maxJer, *maxSna, *maxCra, *maxPop}
	for i, ov := range orders {
		for _, v := range orders[i:] {
			if v > 0 {
				cfg.Params = append(cfg.Params, motion.Parameter{Max: ov})
				break
			}
		}
	}

	p, err := cfg.Solve()
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	w.Write([]string{"Time", "Pos", "Vel", "Acc", "Jer", "Sna", "Cra", "Pop"})

	for i := 0.0; i <= p.Duration(); i += *t {
		s := p.State(i)
		w.Write([]string{
			strconv.FormatFloat(i, 'f', -1, 64),
			strconv.FormatFloat(s.Value(0), 'f', -1, 64),
			strconv.FormatFloat(s.Value(1), 'f', -1, 64),
			strconv.FormatFloat(s.Value(2), 'f', -1, 64),
			strconv.FormatFloat(s.Value(3), 'f', -1, 64),
			strconv.FormatFloat(s.Value(4), 'f', -1, 64),
			strconv.FormatFloat(s.Value(5), 'f', -1, 64),
			strconv.FormatFloat(s.Value(6), 'f', -1, 64),
		})
		w.Flush()
	}
}
