package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-ml/2-linear-regression/linreg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"io"
	"log"
	"os"

	"gonum.org/v1/plot"
)

func main () {
	itterations := flag.Int("n", 100, "number of itterations")
	outPath := flag.String("0", "out.png", "path to output file")
	flag.Parse()

	inPath := flag.Arg(0)
	if inPath == "" {
		inPath = "data.txt"
	}

	inFile, err := os.Open(inPath)
	if err != nil {
		log.Fatalf("could not open file %s: %v", inPath, err)
	}
	defer inFile.Close()

	xs, ys, err := readData(inFile)
	if err != nil {
		log.Fatalf("could not read $s: %v", inPath, err)
	}

	outFile, err := os.Create(*outPath)
	if err != nil {
		log.Fatalf("could not create output file %s: %v", *outPath, err)
	}

	err = plotData(outFile, xs, ys, *itterations)
	if err != nil {
		log.Fatalf("could not plot data: %v", err)
	}

	err = outFile.Close()
	if err != nil {
		log.Fatalf("could not close file %s: %v", *outPath, err)
	}
}

type xy struct {
	x, y float64
}

func readData(data io.Reader) (xs, ys []float64, err error) {
	s := bufio.NewScanner(data)
	for s.Scan() {
		var x, y float64
		_, err := fmt.Sscanf(s.Text(), "%f,%f", &x, &y)
		if err != nil {
			log.Printf("discarding bad data points %q: %v", s.Text(), err)
			continue
		}
		xs = append(xs, x)
		ys = append(ys, y)
	}

	if err := s.Err(); err != nil {
		return nil, nil, fmt.Errorf("could not scan: %v", err)
	}

	return xs, ys, nil
}

type xyer struct {
	xs, ys []float64
}

func (x xyer) Len() int  { return len(x.xs) }
func (x xyer) XY(i int) (float64, float64) { return x.xs[i], x.ys[i] }

func plotData(out io.Writer, xs, ys []float64, itterations int) error {
	p, err := plot.New()
	if err != nil {
		return fmt.Errorf("could not create plot: %v", err)
	}

	// create scatter plot for all data points
	s, err := plotter.NewScatter(xyer{xs, ys})
	if err != nil {
		return fmt.Errorf("could not create scatter plot: %v", err)
	}
	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)

	x, c := linreg.LinearRegression(xs, ys, itterations, 0.01)

	// create a regression line
	l, err := plotter.NewLine(plotter.XYs{
		{1, 1 * x + c}, {20, 20 * x + c},
	})
	if err != nil {
		return fmt.Errorf("could not create regression line: %v", err)
	}

	p.Add(l)

	wt, err := p.WriterTo(256, 256, "png")
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}

	_, err = wt.WriteTo(out)
	if err != nil {
		return fmt.Errorf("could not write: %v", err)
	}

	return nil
}
