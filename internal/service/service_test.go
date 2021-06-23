package service

import (
	"context"
	"testing"

	"github.com/Haba1234/coordinateStoradge/internal/app"
	"github.com/Haba1234/coordinateStoradge/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	t.Parallel()
	t.Run("search points", func(t *testing.T) {
		t.Parallel()
		s := storage.NewStorage()
		search := NewSearch(s)
		for i := 0; i < 10; i++ {
			p := app.Point{X: uint32(i), Y: uint32(i)}
			search.SavePoint(p)
		}
		p := app.Point{X: uint32(4), Y: uint32(7)}
		search.SavePoint(p)
		p = app.Point{X: uint32(3), Y: uint32(6)}
		search.SavePoint(p)
		p = app.Point{X: uint32(7), Y: uint32(4)}
		search.SavePoint(p)
		p = app.Point{X: uint32(7), Y: uint32(2)}
		search.SavePoint(p)

		result := search.SearchNeighbors(context.Background(), p)
		require.Equal(t, []app.Point{{X: 7, Y: 4}, {X: 4, Y: 4}, {X: 5, Y: 5}}, result)
	})

	t.Run("one point", func(t *testing.T) {
		t.Parallel()
		s := storage.NewStorage()
		search := NewSearch(s)
		p := app.Point{X: uint32(2), Y: uint32(1)}
		search.SavePoint(p)
		result := search.SearchNeighbors(context.Background(), p)
		require.Nil(t, result)
	})
}
