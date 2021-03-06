package service

import (
	"github.com/Haba1234/coordinateStoradge/internal/app"
)

func xy2dMorton(point app.Point) uint64 {
	xx := uint64(point.X)
	xx = (xx | (xx << 16)) & 0x0000FFFF0000FFFF
	xx = (xx | (xx << 8)) & 0x00FF00FF00FF00FF
	xx = (xx | (xx << 4)) & 0x0F0F0F0F0F0F0F0F
	xx = (xx | (xx << 2)) & 0x3333333333333333
	xx = (xx | (xx << 1)) & 0x5555555555555555

	yy := uint64(point.Y)
	yy = (yy | (yy << 16)) & 0x0000FFFF0000FFFF
	yy = (yy | (yy << 8)) & 0x00FF00FF00FF00FF
	yy = (yy | (yy << 4)) & 0x0F0F0F0F0F0F0F0F
	yy = (yy | (yy << 2)) & 0x3333333333333333
	yy = (yy | (yy << 1)) & 0x5555555555555555

	return xx | (yy << 1)
}

func morton1(x uint64) uint64 {
	x &= 0x5555555555555555
	x = (x | (x >> 1)) & 0x3333333333333333
	x = (x | (x >> 2)) & 0x0F0F0F0F0F0F0F0F
	x = (x | (x >> 4)) & 0x00FF00FF00FF00FF
	x = (x | (x >> 8)) & 0x0000FFFF0000FFFF
	x = (x | (x >> 16)) & 0xFFFFFFFFFFFFFFFF
	return x
}

func d2xyMorton(d uint64) app.Point {
	x := morton1(d)
	y := morton1(d >> 1)
	return app.Point{
		X: uint32(x),
		Y: uint32(y),
	}
}
