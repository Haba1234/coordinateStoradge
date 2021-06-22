package app

type Point struct {
	X, Y uint32
}

const (
	MaxPoint = 3     // Кол-во искомых соседей
	MaxLimit = 10000 // Максимально допустимая координатна точки
)
