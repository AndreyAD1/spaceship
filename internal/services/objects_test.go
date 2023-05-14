package services_test

import (
	"reflect"
	"testing"

	"github.com/AndreyAD1/spaceship/internal/services"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"
)

func TestBaseObject_GetCornerCoordinates(t *testing.T) {
	type fields struct {
		X float64
		Y float64
	}
	tests := []struct {
		name      string
		fields    fields
		expectedX int
		expectedY int
	}{
		{"common", fields{4.3, 7.2}, 4, 7},
		{"zeros", fields{0, 0}, 0, 0},
		{"negative", fields{-4.3, -7.8}, -4, -8},
		{"integer", fields{1.0, 5.0}, 1, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseObject := &services.BaseObject{
				IsBlocked: false,
				IsDrawn:   false,
				Active:    true,
				X:         tt.fields.X,
				Y:         tt.fields.Y,
				Style:     tcell.StyleDefault.Background(tcell.ColorReset),
				Speed:     1,
				View:      "A",
			}
			x, y := baseObject.GetCornerCoordinates()
			require.Equal(t, tt.expectedX, x)
			require.Equal(t, tt.expectedY, y)
		})
	}
}

func TestBaseObject_GetViewCoordinates(t *testing.T) {
	type fields struct {
		X    float64
		Y    float64
		View string
	}
	tests := []struct {
		name     string
		fields   fields
		expected [][]int
	}{
		{
			"shell 0-0",
			fields{0.0, 0.0, "|"},
			[][]int{{0, 0}},
		},
		{
			"shell 3-5",
			fields{2.8, 5.1, "|"},
			[][]int{{3, 5}},
		},
		{
			"multiline",
			fields{
				2.8,
				5.1,
				`  __  
||  |
`,
			},
			[][]int{{5, 5}, {6, 5}, {3, 6}, {4, 6}, {7, 6}},
		},
		{
			"meteorite",
			fields{
				2.3,
				-1.4,
				services.MeteoriteView,
			},
			[][]int{
				{4, -1}, {5, -1}, {6, -1},
				{3, 0}, {7, 0},
				{2, 1}, {8, 1},
				{2, 2}, {3, 2}, {4, 2}, {5, 2}, {6, 2}, {7, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseObject := &services.BaseObject{
				IsBlocked: false,
				IsDrawn:   false,
				Active:    true,
				X:         tt.fields.X,
				Y:         tt.fields.Y,
				Style:     tcell.StyleDefault.Background(tcell.ColorReset),
				Speed:     1,
				View:      tt.fields.View,
			}
			coords, _ := baseObject.GetViewCoordinates()
			require.Condition(
				t,
				func() bool { return reflect.DeepEqual(coords, tt.expected) },
				"BaseObject.GetViewCoordinates() = %v, want %v",
				coords,
				tt.expected,
			)
		})
	}
}
