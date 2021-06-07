package experiment

import (
	graph "MPC/Graph"
)

//Increment multiplication
//Uses 10 peers
func IncrementMult(plotter graph.Interface) {
	incrementMultWithDelay(plotter, 0, 2000, 50000, 2000, makeRandomBalancedMultExpression, 0)
	plotter.Plot()
}
