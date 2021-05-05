package experiment

import (
	graph "MPC/Graph"
	"time"
)

func IncrementMultBalanced(plotter graph.Interface) {
	for delay := 10; delay <= 200; delay *= 2 {
		incrementMultWithDelay(plotter, time.Duration(delay)*time.Millisecond, 100, 500, 100, makeRandomBalancedMultExpression)
	}
	plotter.Plot()
}
