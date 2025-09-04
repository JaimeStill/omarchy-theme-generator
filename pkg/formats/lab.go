package formats

import (
	"fmt"
)

type LAB struct {
	L float64
	A float64
	B float64
}

func (lab LAB) RGBA() (r, g, b, a uint32) {
	rgba := LABToRGBA(lab)
	r = uint32(rgba.R) * 0x101
	g = uint32(rgba.G) * 0x101
	b = uint32(rgba.B) * 0x101
	a = uint32(rgba.A) * 0x101
	return
}

func NewLAB(l, a, b float64) LAB {
	return LAB{L: l, A: a, B: b}
}

func (lab LAB) IsValid() bool {
	return lab.L >= 0 && lab.L <= 100 &&
		lab.A >= -128 && lab.A <= 127 &&
		lab.B >= -128 && lab.B <= 127
}

func (lab LAB) String() string {
	return fmt.Sprintf("LAB(%.2f, %.2f, %.2f)", lab.L, lab.A, lab.B)
}
