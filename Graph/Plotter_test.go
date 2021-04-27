package graph

import (
	experiment "MPC/Experiment"
	"testing"
)

func TestPlotting(t *testing.T) {
	xy := []experiment.XY{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8}, {9, 9}, {10, 10}}

	err := plotGraph("test", xy, "Test", "png")
	if err != nil {
		t.Errorf("Failed the test with the error %v", err)
	}
}
