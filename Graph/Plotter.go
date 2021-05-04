package graph

import (
	prot "MPC/Protocol"
	"fmt"
	"os"

	plot "gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg/draw"
)

type Series struct {
	Data []XY
	Name string
}

type Plotter struct {
	data     []Series
	format   string
	seriesNr int
	title    string
	xAxis    string
}

type XY struct {
	X float64
	Y prot.Times
}

func MkPlotter(title string, firstSeriesName string, format string, xAxisName string) *Plotter {
	plotter := new(Plotter)
	plotter.format = format
	plotter.title = title
	firstSerie := Series{Name: firstSeriesName}
	plotter.data = append(plotter.data, firstSerie)
	plotter.xAxis = xAxisName
	return plotter
}

func (gp *Plotter) convertToPlotFormat() []plotter.XYs {
	var result []plotter.XYs
	for _, series := range gp.data {
		fXY := make(plotter.XYs, len(series.Data))
		for j, data := range series.Data {
			fXY[j].X = data.X
			fXY[j].Y = float64(SumTimes(data.Y))
		}
		result = append(result, fXY)
	}

	return result
}

func (gp *Plotter) Plot() error {
	filePath := gp.title + "." + gp.format

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Could not create file: %v", filePath)
	}
	defer f.Close()
	p := plot.New()
	p.Add(plotter.NewGrid())
	p.Title.Text = gp.title
	p.Title.TextStyle.Font.Variant = "Mono"
	p.Y.Label.Text = "Time (ms)"
	p.Y.Label.Padding = font.Length(20)
	p.X.Label.Text = gp.xAxis
	p.X.Label.Padding = font.Length(20)
	series := gp.convertToPlotFormat()
	for i, serie := range series {
		scatter, _ := plotter.NewScatter(serie)
		scatter.GlyphStyle.Color = plotutil.SoftColors[i]
		scatter.Shape = draw.CircleGlyph{}
		p.Add(scatter)
		p.Legend.Add(gp.data[i].Name, scatter)
	}

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
	gp.data[gp.seriesNr].Data = append(gp.data[gp.seriesNr].Data, XY{X: float64(variable), Y: *data})
}

func (gp *Plotter) NewSeries(name string) {
	gp.seriesNr = gp.seriesNr + 1
	newSerie := Series{Name: name}
	gp.data = append(gp.data, newSerie)
}

func SumTimes(timer prot.Times) int64 {
	protTime := timer.Calculate + timer.Preprocess + timer.SetupTree
	return protTime.Milliseconds()
}
