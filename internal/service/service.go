package service

import (
	"context"
	"math"

	"coordinateStoradge/internal/app"
	"coordinateStoradge/internal/storage"
)

type distance struct {
	point app.Point // Координаты точки.
	l     float64   // Расстояние между точками.
}

// Search структура БЛ.
type Search struct {
	storage   *storage.Storage
	neighbors [app.MaxPoint]distance
	ind       int // Найдено соседей.
}

func NewSearch(storage *storage.Storage) *Search {
	return &Search{
		storage: storage,
	}
}

// SavePoint преобразование координат в Z и запись в архив.
func (s *Search) SavePoint(point app.Point) {
	z := xy2dMorton(point)
	s.storage.AddPoint(z)
}

// SearchNeighbors функция поиска ближайших соседей заданной точки.
func (s *Search) SearchNeighbors(ctx context.Context, point app.Point) []app.Point {
	s.storage.RLock()
	defer s.storage.RUnlock()

	p := point
	s.neighbors = [app.MaxPoint]distance{}
	s.ind = 0
	z := xy2dMorton(p)
	if _, ok := s.storage.ReadPoint(z); !ok || s.storage.Len() <= 1 {
		return nil // Точка не существует или больше точек нет.
	}

	d := uint32(1) // Размер области поиска, начинаем с d = 1.
	for {
		zStart, zFinish := extremePoints(d, p)
		for i := zStart; i <= zFinish; i++ {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			if i == z {
				continue
			}

			if _, ok := s.storage.ReadPoint(i); ok {
				result := d2xyMorton(i)
				// Расстояние между точками.
				l := math.Pow(float64(int32(result.X)-int32(p.X)), 2) + math.Pow(float64(int32(result.Y)-int32(p.Y)), 2)
				s.saveNeighbors(l, result)
			}
		}
		// Готовим слайс соседей для отправки.
		if s.ind >= app.MaxPoint || s.ind >= s.storage.Len() {
			points := make([]app.Point, app.MaxPoint)
			for i, val := range s.neighbors {
				points[i] = val.point
			}
			return points
		}
		d++
	}
}

// extremePoints вычисление границ области поиска соседних точек.
func extremePoints(d uint32, point app.Point) (uint64, uint64) {
	var startP, finishP app.Point
	startP.X = point.X - d
	if int32(point.X)-int32(d) < 0 {
		startP.X = 0
	}
	startP.Y = point.Y - d
	if int32(point.Y)-int32(d) < 0 {
		startP.Y = 0
	}
	zStart := xy2dMorton(startP)

	finishP.X = point.X + d
	if int32(point.X)+int32(d) > app.MaxLimit {
		finishP.X = app.MaxLimit
	}
	finishP.Y = point.Y + d
	if int32(point.Y)+int32(d) > app.MaxLimit {
		finishP.Y = app.MaxLimit
	}
	zFinish := xy2dMorton(finishP)

	return zStart, zFinish
}

// saveNeighbors сохраняет в слайс соседей и при необходимости сортирует по дальности.
func (s *Search) saveNeighbors(l float64, point app.Point) {
	for i, v := range s.neighbors {
		if v.point == point {
			break // Такая точка уже найдена.
		}
		if v.l == 0.0 || v.l > l {
			s.neighbors[i] = distance{
				point: point,
				l:     l,
			}
			s.ind++
			// двигаем остальных соседей дальше, крайнего выкидываем
			for j := i + 1; j < app.MaxPoint; j++ {
				s.neighbors[j], v = v, s.neighbors[j]
			}
			break
		}
	}
}
