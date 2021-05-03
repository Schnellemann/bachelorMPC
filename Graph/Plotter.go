package graph

import (
	prot "MPC/Protocol"
	"fmt"
	"os"

	plot "gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type Plotter struct {
	data   []XY
	format string
}

type XY struct {
	X float64
	Y prot.Times
}

func MkPlotter(format string) *Plotter {
	plotter := new(Plotter)
	plotter.format = format
	return plotter
}

func (gp *Plotter) convertToPlotFormat() plotter.XYs {
	fXY := make(plotter.XYs, len(gp.data))
	for i, data := range gp.data {
		fXY[i].X = data.X
		fXY[i].Y = float64(SumTimes(data.Y))
	}
	return fXY
}

func (gp *Plotter) Plot(fileName string, title string) error {
	filePath := fileName + "." + gp.format

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Could not create file: %v", filePath)
	}
	defer f.Close()
	p := plot.New()

	scatter, _ := plotter.NewScatter(gp.convertToPlotFormat())
	p.Add(scatter)

	wt, err := p.WriterTo(512, 512, gp.format)
	if err != nil {
		return fmt.Errorf("Could not write to %s: %v", filePath, err)
	}
	wt.WriteTo(f)
	if err := f.Close(); err != nil {
		return fmt.Errorf("Could not close file: %v", filePath)
	}

	return nil

}

func (gp *Plotter) AddData(variable int, data *prot.Times) {
	gp.data = append(gp.data, XY{X: float64(variable), Y: *data})
}

func SumTimes(timer prot.Times) int64 {
	protTime := timer.Calculate + timer.Preprocess + timer.SetupTree
	return protTime.Milliseconds()
}
