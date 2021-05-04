package graph

import (
	prot "MPC/Protocol"
	"testing"
	"time"
)

func TestPlotting(t *testing.T) {
	graph := MkPlotter("Testing", "TestSeries", "png", "Variable")
	timer := new(prot.Times)
	timer.Calculate = time.Minute + 2*time.Second
	timer.Network = 5 * time.Second
	timer.Preprocess = 0
	timer.SetupTree = 10 * time.Millisecond
	graph.AddData(1, timer)
	graph.AddData(2, timer)
	graph.AddData(10, timer)

	graph.NewSeries("Test2")
	timer.Calculate = time.Minute + 5*time.Second
	timer.Network = 5 * time.Second
	timer.Preprocess = 0
	timer.SetupTree = 10 * time.Millisecond
	graph.AddData(1, timer)
	graph.AddData(2, timer)
	graph.AddData(10, timer)
	err := graph.Plot()
	if err != nil {
		t.Errorf("Failed the test with the error %v", err)
	}
}
