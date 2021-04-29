package graph

import (
	"fmt"
	"os"

	plot "gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type XY struct {
	X, Y float64
}

func convertToPlotFormat(xys []XY) plotter.XYs {
	fXY := make(plotter.XYs, len(xys))
	for i, xy := range xys {
		fXY[i].X = xy.X
		fXY[i].Y = xy.Y
	}
	return fXY
}

func PlotGraph(fileName string, xy []XY, title string, format string) error {
	filePath := fileName + "." + format

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Could not create file: %v", filePath)
	}
	defer f.Close()
	p := plot.New()

	scatter, _ := plotter.NewScatter(convertToPlotFormat(xy))
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