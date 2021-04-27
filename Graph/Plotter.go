package graph

import (
	"fmt"
	"os"
	"time"

	plot "gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func plotGraph(fileName string, Xdata []int, Ydata []time.Duration, title string, format string) error {
	filePath := fileName + "." + format

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Could not create file: %v", filePath)
	}
	defer f.Close()
	p := plot.New()

	scatter, _ := plotter.NewScatter(plotter.XYs{
		{0, 0}, {0.5, 0.5}, {1, 1},
	})
	p.Add(scatter)

	wt, err := p.WriterTo(512, 512, format)
	if err != nil {
		return fmt.Errorf("Could not write to %s: %v", filePath, err)
	}
	wt.WriteTo(f)
	if err := f.Close(); err != nil {
		return fmt.Errorf("Could not close file: %v", filePath)
	}

	return nil

}
