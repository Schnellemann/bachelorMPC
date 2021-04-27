package graph

import (
	"testing"
	"time"
)

func TestPlotting(t *testing.T) {
	var x []int
	var y []time.Duration

	err := plotGraph("test", x, y, "Test", "png")
	if err != nil {
		t.Errorf("Failed the test with the error %v", err)
	}
}
