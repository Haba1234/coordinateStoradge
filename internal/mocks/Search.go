// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	app "github.com/Haba1234/coordinateStoradge/internal/app"

	mock "github.com/stretchr/testify/mock"
)

// Search is an autogenerated mock type for the Search type
type Search struct {
	mock.Mock
}

// SavePoint provides a mock function with given fields: point
func (_m *Search) SavePoint(point app.Point) {
	_m.Called(point)
}

// SearchNeighbors provides a mock function with given fields: ctx, point
func (_m *Search) SearchNeighbors(ctx context.Context, point app.Point) []app.Point {
	ret := _m.Called(ctx, point)

	var r0 []app.Point
	if rf, ok := ret.Get(0).(func(context.Context, app.Point) []app.Point); ok {
		r0 = rf(ctx, point)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]app.Point)
		}
	}

	return r0
}