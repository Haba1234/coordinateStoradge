package service

import (
	"context"
	"log"
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

func (s *Search) SavePoint(point app.Point) {
	z := xy2dMorton(point.X, point.Y)
	s.storage.AddPoint(z)
}

func (s *Search) SearchNeighbors(ctx context.Context, point app.Point) []app.Point {
	s.storage.RLock()
	defer s.storage.RUnlock()
	p := point
	s.neighbors = [app.MaxPoint]distance{}
	s.ind = 0
	z := xy2dMorton(p.X, p.Y)
	if _, ok := s.storage.ReadPoint(z); !ok || s.storage.Len() <= 1 {
		log.Println("Точка не существует или больше точек нет", s.storage.Len(), ok)
		return nil
	}
	var d uint32 //nolint:gosimple
	d = 1        // Размер области поиска, начинаем с d = 1.
	for {
		zStart, zFinish := extremePoints(d, p)
		for i := zStart; i <= zFinish; i++ {
			select {
			case <-ctx.Done():
				log.Println("Получено прерывание")
				return nil
			default:
			}
			if i == z {
				continue
			}
			_, ok := s.storage.ReadPoint(i)
			if ok {
				x, y := d2xyMorton(i)
				// Расстояние между точками.
				l := math.Pow(float64(int32(x)-int32(p.X)), 2) + math.Pow(float64(int32(y)-int32(p.Y)), 2)
				s.saveNeighbors(l, app.Point{X: x, Y: y})
			}
		}
		// Готовим слайс соседей для отправки.
		if s.ind >= app.MaxPoint || s.ind >= s.storage.Len() {
			points := make([]app.Point, app.MaxPoint)
			for i, val := range s.neighbors {
				points[i] = val.point
			}
			log.Println("Найдены:", points)
			return points
		}
		d++
	}
}

// extremePoints вычисление границ области поиска соседних точек.
func extremePoints(d uint32, point app.Point) (uint64, uint64) {
	dX := point.X - d
	if int32(point.X)-int32(d) < 0 {
		dX = 0
	}
	dY := point.Y - d
	if int32(point.Y)-int32(d) < 0 {
		dY = 0
	}
	zStart := xy2dMorton(dX, dY)

	dX = point.X + d
	if int32(point.X)+int32(d) > app.MaxLimit {
		dX = app.MaxLimit
	}
	dY = point.Y + d
	if int32(point.Y)+int32(d) > app.MaxLimit {
		dY = app.MaxLimit
	}
	zFinish := xy2dMorton(dX, dY)

	return zStart, zFinish
}

// saveNeighbors сохраняет в слайс соседей и при необходимости сортирует по дальности.
func (s *Search) saveNeighbors(l float64, point app.Point) {
	for i, v := range s.neighbors {
		if v.point.X == point.X && v.point.Y == point.Y {
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
