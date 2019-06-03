package linreg

import "fmt"

func LinearRegression(xs, ys []float64, itterations int, alpha float64) (m, c float64) {
	for i:=0; i < itterations; i++  {
		loss, dm, dc := Gradient(xs, ys, m, c)
		m += -dm * alpha
		c += -dc * alpha

		if (10 * i % itterations) == 0 {
			fmt.Printf("loss (%.2f, %.2f) = %.2f\n", m, c, loss)
		}
	}
	
	return m, c
}

func Gradient(xs, ys []float64, m, c float64) (loss, dm, dc float64) {
	for i := range xs {
		d := ys[i] - (xs[i] * m + c)
		loss += d * d
		dm += -xs[i] * d
		dc += -d
	}

	n := float64(len(xs))

	return loss / n, 2 / n * dm, 2 / n * dc
}
