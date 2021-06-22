package service

import (
	"testing"

	"coordinateStoradge/internal/app"
	"github.com/stretchr/testify/require"
)

func TestMorton(t *testing.T) {
	t.Parallel()
	point := app.Point{
		X: uint32(12345),
		Y: uint32(54321),
	}

	t.Run("Convert 2d to Z", func(t *testing.T) {
		z := xy2dMorton(point)
		require.Equal(t, uint64(0xa7200f43), z)
	})

	t.Run("Convert Z to 2d", func(t *testing.T) {
		z := uint64(0xa7200f43)
		p := d2xyMorton(z)
		require.Equal(t, point, p)
	})
}
