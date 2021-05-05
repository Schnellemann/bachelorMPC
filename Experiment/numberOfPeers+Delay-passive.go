package experiment

import (
	graph "MPC/Graph"
	"time"
)

func IncrementPeerAndDelay(plotter graph.Interface) {
	for delay := 10; delay <= 200; delay *= 2 {
		incrementPeerWithDelay(plotter, time.Duration(delay)*time.Millisecond, 10, 50, 10, makeRandomMultExpression)
	}
	plotter.Plot()
}
