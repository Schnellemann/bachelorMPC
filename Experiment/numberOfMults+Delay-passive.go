package experiment

import graph "MPC/Graph"

func IncrementMultAndDelay(plotter graph.Interface) {
	for delay := 0; delay < 110; delay += 10 {
		incrementMultWithDelay(plotter, delay, 200, 4000, 200)
	}
	plotter.Plot()
}
