package formats

type XYZ struct {
	X float64
	Y float64
	Z float64
}

var (
	D65Illuminant = XYZ{X: 95.047, Y: 100.000, Z: 108.883}
	D50Illuminant = XYZ{X: 96.422, Y: 100.00, Z: 82.521}
)

func NewXYZ(x, y, z float64) XYZ {
	return XYZ{X: x, Y: y, Z: z}
}
