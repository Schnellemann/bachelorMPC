package graph

import (
	prot "MPC/Protocol"
	"testing"
	"time"
)

func TestExcel(t *testing.T) {
	graph := MkExcel("test", "testVariable")
	graph.NewSeries("firstSeries")
	timer := new(prot.Times)
	timer.Calculate = time.Minute + 2*time.Second
	timer.Network = 5 * time.Second
	timer.Preprocess = 0
	timer.SetupTree = 10 * time.Millisecond
	graph.AddData(2, timer)
	graph.AddData(5, timer)
	graph.NewSeries("secondSeries")
	graph.AddData(1, timer)
	graph.AddData(2, timer)
	graph.AddData(3, timer)
	err := graph.Plot()
	if err != nil {
		t.Errorf("Failed the test with the error %v", err)
	}
}
