package experiment

import (
	graph "MPC/Graph"
	"time"
)

func IncrementPeerAndDelay(plotter graph.Interface) {
	for delay := 160; delay <= 160; delay *= 2 {
		incrementPeerWithDelay(plotter, time.Duration(delay)*time.Millisecond, 10, 90, 10, makeRandomMultExpression)
	}
	plotter.Plot()
}
